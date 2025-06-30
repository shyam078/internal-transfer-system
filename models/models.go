package models

type Account struct {
    AccountID int64  `json:"account_id"`
    Balance   string `json:"balance"`
}

type Transaction struct {
    SourceAccountID      int64  `json:"source_account_id"`
    DestinationAccountID int64  `json:"destination_account_id"`
    Amount               string `json:"amount"`
}
