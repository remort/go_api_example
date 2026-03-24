package handlers

import (
	"net/http"
)

func InitRouter() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /api/v1/wallet", HandleWalletDepositWithdraw)
	r.HandleFunc("GET /api/v1/wallet/{wallet_id}", HandleGetWallet)
	return r
}
