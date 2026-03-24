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
	ValletId      uuid.UUID       `json:"valletId"`
	OperationType WalletOperation `json:"operationType"`
	Amount        int             `json:"amount"`
}

type Wallet struct {
	Id        uuid.UUID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
