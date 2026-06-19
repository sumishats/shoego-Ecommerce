package models

type WalletResponse struct{
	Balance float64 `json:"balance"`
}
type WalletTransactionResponse struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}