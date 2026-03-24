package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"web-example/src/db"
	"web-example/src/types"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func sendJSON(w http.ResponseWriter, status int, data any) error {
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
	log.Println("In handleWalletDepositWithdraw")

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
	log.Println("In handleGetWallet")
	wallet_id := r.PathValue("wallet_id")

	if err := uuid.Validate(wallet_id); err != nil {
		sendJSON(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var wallet types.Wallet
	result := db.DB.Where("id = ?", wallet_id).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendJSON(w, http.StatusNotFound, "Wallet not found")
		return
	}

	log.Printf("Wallet in DB: %v", wallet)

	sendJSON(w, http.StatusOK, wallet)
}
