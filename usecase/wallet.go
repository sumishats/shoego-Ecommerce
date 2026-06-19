package usecase

import (
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
)

func CreditWallet(userID uint, amount float64, description string) error {

	wallet, err := repository.GetWalletByUserID(userID)

	if err != nil {
		wallet = &domain.Wallet{
			UserID:  userID,
			Balance: 0,
		}
		err = repository.CreateWallet(wallet)
		if err != nil {
			return err
		}
	}
	wallet.Balance += amount
	err = repository.UpdateWalletBalance(wallet.ID, wallet.Balance)
	if err != nil {
		return err
	}
	transaction := &domain.WalletTransaction{
		WalletID:    wallet.ID,
		Amount:      amount,
		Type:        "credit",
		Description: description,
	}
	return repository.CreateWalletTransaction(transaction)
}

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
