package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"web-example/src/db"
	"web-example/src/handlers"
	"web-example/src/types"

	"github.com/google/uuid"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := handlers.InitRouter()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable() {
	db.DB.Exec("DELETE FROM wallets")
}

func createWalletInDB() types.Wallet {
	user_id := uuid.New()
	wallet := types.Wallet{
		Id:        user_id,
		Amount:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.DB.Create(&wallet)
	return wallet
}
