package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"web-example/src/db"
	"web-example/src/types"

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
		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonStr))
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
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonStr))
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
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
