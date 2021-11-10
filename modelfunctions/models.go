package modelfunctions

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id,omitempty"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type Listaccount struct {
	Account
	ListAccountsParams
}

type Listjson struct {
	Data []Account `json:"data"`
}

type DepositMoneyParams struct {
	Amount    int64     `json:"amount"`
	Owner     string    `json:"owner"`
	Reference string    `json:"reference"`
	DepositID string    `json:"deposit_id"`
	CreatedAt time.Time `json:"created_at"`
}

type WithdrawMoneyParams struct {
	Amount     int64     `json:"amount"`
	Owner      string    `json:"owner"`
	Reference  string    `json:"reference"`
	WithdrawID string    `json:"withdraw_id"`
	CreatedAt  time.Time `json:"created_at"`
}
type WithdrawInfo struct {
	Amount    int64     `json:"amount"`
	Owner     string    `json:"owner"`
	Reference string    `json:"reference"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	ID        int64     `json:"id,omitempty"`
}

type DepositInfo struct {
	Amount    int64     `json:"amount"`
	Owner     string    `json:"owner"`
	Reference string    `json:"reference"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	ID        int64     `json:"id,omitempty"`
}

type TransferInfo struct {
	Amount    int64     `json:"amount"`
	Sender    string    `json:"sender"`
	Reference string    `json:"reference"`
	Balance   int64     `json:"balance"`
	Receiver  string    `json:"receiver"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
}
type TransferParams struct {
	Amount    int64  `json:"amount"`
	Sender    string `json:"sender"`
	Reference string `json:"reference"`
	Receiver  string `json:"receiver"`
}
