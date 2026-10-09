package usecase

import (
	"errors"
	"shoego/database"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"

	"gorm.io/gorm"
)

func CreditWallet(userID uint, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var wallet domain.Wallet

		err := tx.Where("user_id = ?", userID).First(&wallet).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = domain.Wallet{
				UserID:  userID,
				Balance: 0,
			}

			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		wallet.Balance += amount

		if err := tx.Model(&wallet).
			Update("balance", wallet.Balance).Error; err != nil {
			return err
		}

		transaction := domain.WalletTransaction{
			WalletID:    wallet.ID,
			Amount:      amount,
			Type:        "credit",
			Description: description,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// func CreditWallet(userID uint, amount float64, description string) error {

// 	wallet, err := repository.GetWalletByUserID(userID)

// 	if err != nil {
// 		wallet = &domain.Wallet{
// 			UserID:  userID,
// 			Balance: 0,
// 		}
// 		err = repository.CreateWallet(wallet)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	wallet.Balance += amount
// 	err = repository.UpdateWalletBalance(wallet.ID, wallet.Balance)
// 	if err != nil {
// 		return err
// 	}
// 	transaction := &domain.WalletTransaction{
// 		WalletID:    wallet.ID,
// 		Amount:      amount,
// 		Type:        "credit",
// 		Description: description,
// 	}
// 	return repository.CreateWalletTransaction(transaction)
// }

func GetWallet(userID uint) (*models.WalletResponse, error) {

	wallet, err := repository.GetWalletByUserID(userID)

	if err != nil {

		return &models.WalletResponse{
			Balance: 0,
		}, nil
	}

	return &models.WalletResponse{
		Balance: wallet.Balance,
	}, nil
}

func GetWalletHistory(userID uint) ([]models.WalletTransactionResponse, error) {
	wallet, err := repository.GetWalletByUserID(userID)

	if err != nil {
		return []models.WalletTransactionResponse{}, nil
	}

	transactions, err := repository.GetWalletTransactions(wallet.ID)
	if err != nil {
		return nil, err
	}
	var result []models.WalletTransactionResponse
	for _, tx := range transactions {
		result = append(result, models.WalletTransactionResponse{
			Amount:      tx.Amount,
			Type:        tx.Type,
			Description: tx.Description,
			CreatedAt:   tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}
