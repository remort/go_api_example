package handlers

import (
	"net/http"
)

func InitRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST /api/v1/wallet/change-balance", HandleWalletDepositWithdraw)
	r.HandleFunc("GET /api/v1/wallet/{wallet_id}", HandleGetWallet)
	r.HandleFunc("POST /api/v1/wallet", HandleCreateWallet)
	r.HandleFunc("GET /api/v1/wallet", HandleListWallets)
	r.HandleFunc("DELETE /api/v1/wallet/{wallet_id}", HandleDeleteWallet)

	return r
}
