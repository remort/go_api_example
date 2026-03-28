package tests

import (
	"bytes"
	"encoding/json"
	"go-api-example/src/db"
	"go-api-example/src/types"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
)

var invalidUUID = "e2266d5a-8804-4d0c-91de-"

func TestMain(m *testing.M) {
	db.InitDB(&db.DB)
	defer db.DB.Close()
	db.DB.AutoMigrate(&types.Wallet{})
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/wallet/e2266d5a-8804-4d0c-91de-3b7dd18c12c6", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestInvalidID(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/api/v1/wallet/"+invalidUUID, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestWalletExists(t *testing.T) {
	clearTable()
	wallet := createWalletInDB()
	log.Printf("Wallet: %v", wallet)

	req, _ := http.NewRequest("GET", "/api/v1/wallet/"+wallet.Id.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	resp := &types.Wallet{}
	json.Unmarshal(response.Body.Bytes(), resp)

	if resp.Id != wallet.Id {
		t.Errorf("Expected %v. Got %v", wallet, resp)
	}
}

func TestWalletDepositWithdraw(t *testing.T) {
	opTypes := []types.WalletOperation{"deposit", "withdraw"}
	for _, opType := range opTypes {
		clearTable()
		wallet := createWalletInDB()
		log.Printf("Wallet: %v", wallet)

		var jsonStr = []byte(`{
			"valletId": "` + wallet.Id.String() + `",
			"operationType": "` + string(opType) + `",
			"amount": 1
		}`)
		req, _ := http.NewRequest("POST", "/api/v1/wallet/change-balance", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		response := executeRequest(req)

		checkResponseCode(t, http.StatusCreated, response.Code)

		resp := &types.Wallet{}
		json.Unmarshal(response.Body.Bytes(), resp)

		if resp.Id != wallet.Id {
			t.Errorf("Expected ID %v. Got %v", wallet.Id, resp.Id)
		}

		var expectedAmount int
		switch opType {
		case types.Deposit:
			expectedAmount = wallet.Amount + 1
		case types.Withdraw:
			expectedAmount = wallet.Amount - 1
		}

		if resp.Amount != expectedAmount {
			t.Errorf("Expected amount %v. Got %v", expectedAmount, resp.Amount)
		}
	}
}

func TestWalletDepositWithdrawInvalidID(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{
		"valletId": "` + invalidUUID + `",
		"operationType": "` + "deposit" + `",
		"amount": 1
	}`)
	req, _ := http.NewRequest("POST", "/api/v1/wallet/change-balance", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestWalletDepositWithdrawInvalidOpType(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{
		"valletId": "` + uuid.New().String() + `",
		"operationType": "` + "decrease" + `",
		"amount": 1
	}`)
	req, _ := http.NewRequest("POST", "/api/v1/wallet/change-balance", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateWalletSuccess(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("POST", "/api/v1/wallet", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var wallet types.Wallet
	json.Unmarshal(response.Body.Bytes(), &wallet)

	if wallet.Id == uuid.Nil {
		t.Errorf("Expected valid UUID, got %v", wallet.Id)
	}

	if wallet.Amount != 0 {
		t.Errorf("Expected amount 0, got %d", wallet.Amount)
	}

	var dbWallet types.Wallet
	db.DB.First(&dbWallet, "id = ?", wallet.Id)
	if dbWallet.Id != wallet.Id {
		t.Errorf("Wallet not found in database")
	}
}

func TestListWalletsWithMultipleWallets(t *testing.T) {
	clearTable()

	expectedWallets := make([]types.Wallet, 3)
	for i := range 3 {
		wallet := createWalletInDB()
		expectedWallets[i] = wallet
	}
	req, _ := http.NewRequest("GET", "/api/v1/wallet", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var wallets []types.Wallet
	json.Unmarshal(response.Body.Bytes(), &wallets)

	if len(wallets) != 3 {
		t.Errorf("Expected 3 wallets, got %d", len(wallets))
	}
	log.Printf("Wallets %v+ %d", wallets, len(wallets))

	for _, expected := range expectedWallets {
		found := false
		for _, actual := range wallets {
			log.Printf("actual.Id %s == expected.Id %s", actual.Id, expected.Id)
			if actual.Id == expected.Id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Wallet %v not found in response", expected.Id)
		}
	}
}

func TestListWalletsEmpty(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/wallet", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var wallets []types.Wallet
	json.Unmarshal(response.Body.Bytes(), &wallets)

	if len(wallets) != 0 {
		t.Errorf("Expected empty list, got %d wallets", len(wallets))
	}
}

func TestDeleteWalletSuccess(t *testing.T) {
	clearTable()
	wallet := createWalletInDB()

	req, _ := http.NewRequest("DELETE", "/api/v1/wallet/"+wallet.Id.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNoContent, response.Code)

	var dbWallet types.Wallet
	result := db.DB.First(&dbWallet, "id = ?", wallet.Id)
	if result.Error == nil {
		t.Errorf("Wallet %v still exists in database after deletion", wallet.Id)
	}
}

func TestDeleteWalletNotFound(t *testing.T) {
	clearTable()

	nonExistentID := uuid.New()

	req, _ := http.NewRequest("DELETE", "/api/v1/wallet/"+nonExistentID.String(), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var errResponse string
	json.Unmarshal(response.Body.Bytes(), &errResponse)

	if errResponse != "Wallet not found" {
		t.Errorf("Expected error message in response")
	}
}

func TestDeleteWalletInvalidUUID(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("DELETE", "/api/v1/wallet/invalid-uuid-format", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var errResponse string
	json.Unmarshal(response.Body.Bytes(), &errResponse)
	log.Printf("Response %v", errResponse)

	if errResponse != "Invalid UUID" {
		t.Errorf("Expected error message for invalid UUID")
	}
}
