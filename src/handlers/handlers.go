package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-api-example/src/db"
	"go-api-example/src/types"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func sendJSON(w http.ResponseWriter, status int, data any) error {
	log.Printf("Response with status code: %d", status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func decodeBody(w http.ResponseWriter, r *http.Request) (types.DepositWithdrawRequest, error) {
	log.Println("Parsing request body")
	var requestBody types.DepositWithdrawRequest
	var errMsg string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		errMsg = "Error parsing request body"
		sendJSON(w, http.StatusBadRequest, errMsg)
		return requestBody, fmt.Errorf("%s", errMsg)
	}
	defer r.Body.Close()

	log.Printf("Parsed request body %v , %v", requestBody.ValletId, requestBody.OperationType)

	switch requestBody.OperationType {
	case types.Deposit:
	case types.Withdraw:
		requestBody.Amount = -requestBody.Amount
	default:
		errMsg = "Invalid operation type"
		sendJSON(w, http.StatusBadRequest, errMsg)
		return requestBody, fmt.Errorf("%s", errMsg)
	}
	return requestBody, nil
}

func HandleWalletDepositWithdraw(w http.ResponseWriter, r *http.Request) {
	log.Println("In POST Deposit/Waithdraw handler")

	requestBody, err := decodeBody(w, r)
	if err != nil {
		log.Println("validation error, return")
		return
	}

	var wallet types.Wallet
	log.Printf("Requested wallet ID: %v", requestBody.ValletId)
	result := db.DB.Where("id = ?", requestBody.ValletId).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendJSON(w, http.StatusNotFound, "Wallet not found")
		return
	}

	log.Printf("Wallet in DB: %v", wallet)
	newAmount := wallet.Amount + requestBody.Amount
	db.DB.Model(&wallet).Update("Amount", newAmount)

	sendJSON(w, http.StatusCreated, wallet)
}

func HandleGetWallet(w http.ResponseWriter, r *http.Request) {
	log.Println("In GET wallet handler")
	walletId := r.PathValue("wallet_id")

	if err := uuid.Validate(walletId); err != nil {
		sendJSON(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var wallet types.Wallet
	result := db.DB.Where("id = ?", walletId).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendJSON(w, http.StatusNotFound, "Wallet not found")
		return
	}

	log.Printf("Wallet in DB: %v", wallet)

	sendJSON(w, http.StatusOK, wallet)
}

func HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	log.Println("In POST Create wallet handler")

	wallet := types.Wallet{
		Id:        uuid.New(),
		Amount:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := db.DB.Create(&wallet)
	if result.Error != nil {
		sendJSON(w, http.StatusBadRequest, "Create wallet error: "+result.Error.Error())
	}

	log.Printf("New wallet ID: %v", wallet.Id)

	sendJSON(w, http.StatusCreated, wallet)
}

func HandleListWallets(w http.ResponseWriter, r *http.Request) {
	log.Println("In GET List wallets handler")

	var wallets []types.Wallet
	result := db.DB.Model(&types.Wallet{}).Limit(10).Find(&wallets)
	if result.Error != nil {
		sendJSON(w, http.StatusBadRequest, "Lookup wallet list error: "+result.Error.Error())
	}

	log.Printf("Found wallets in DB: %v", len(wallets))

	sendJSON(w, http.StatusOK, wallets)
}

func HandleDeleteWallet(w http.ResponseWriter, r *http.Request) {
	log.Println("In DELETE Wallet handler")

	walletId := r.PathValue("wallet_id")
	if err := uuid.Validate(walletId); err != nil {
		log.Print(111)
		sendJSON(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var wallet types.Wallet
	result := db.DB.First(&wallet, "id = ?", walletId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			sendJSON(w, http.StatusNotFound, "Wallet not found")
		} else {
			sendJSON(w, http.StatusBadRequest, "Lookup wallet error: "+result.Error.Error())
		}
		return
	}

	result = db.DB.Delete(&wallet)
	if result.Error != nil {
		sendJSON(w, http.StatusBadRequest, "Delete wallet error: "+result.Error.Error())
		return
	}
	log.Printf("Deleted wallet ID: %v", walletId)

	w.WriteHeader(http.StatusNoContent)
}
