package repository

import (
	"shoego/database"
	"shoego/domain"
)

func GetWalletByUserID(userID uint) (*domain.Wallet, error) {
	var wallet domain.Wallet

	err := database.DB.Where("user_id=?", userID).First(&wallet).Error

	if err != nil {
		return nil, err
	}
	return &wallet, nil

}
func CreateWallet(wallet *domain.Wallet) error {
	return database.DB.Create(wallet).Error
}

func UpdateWalletBalance(walletID uint, balance float64) error {

	return database.DB.Model(&domain.Wallet{}).Where("id = ?", walletID).Update("balance", balance).Error
}

func CreateWalletTransaction(transaction *domain.WalletTransaction) error {
	return database.DB.Create(transaction).Error
}

// wallet history
func GetWalletTransactions(walletID uint) ([]domain.WalletTransaction, error) {
	var transactions []domain.WalletTransaction

	err := database.DB.Where("wallet_id=?", walletID).Order("created_at desc").Find(&transactions).Error

	return transactions, err
}
