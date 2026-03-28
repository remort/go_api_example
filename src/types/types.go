package types

import (
	"time"

	"github.com/google/uuid"
)

type WalletOperation string

const (
	Deposit  WalletOperation = "deposit"
	Withdraw WalletOperation = "withdraw"
)

type DepositWithdrawRequest struct {
	WalletId      uuid.UUID       `json:"walletId"`
	OperationType WalletOperation `json:"operationType"`
	Amount        int             `json:"amount"`
}

type Wallet struct {
	Id        uuid.UUID `json:"id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
