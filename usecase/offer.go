package usecase

import (
	"errors"
	"shoego/domain"
	"shoego/models"
	"shoego/repository"
	"time"
)

func CreateProductOffer(req models.CreateProductOfferRequest) error {

	_, err := repository.GetProductByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return errors.New("invalid start date")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return errors.New("invalid end date")
	}

	if endDate.Before(startDate) {
		return errors.New("end date must be after start date")
	}

	if req.DiscountPercentage <= 0 || req.DiscountPercentage > 90 {
		return errors.New("discount percentage must be between 1 and 90")
	}

	offer := &domain.ProductOffer{
		ProductID:          req.ProductID,
		OfferName:          req.OfferName,
		DiscountPercentage: req.DiscountPercentage,
		StartDate:          startDate,
		EndDate:            endDate,
		IsActive:           true,
	}

	return repository.CreateProductOffer(offer)
}

func GetAllProductOffers() ([]models.ProductOfferResponse, error) {

	offers, err := repository.GetAllProductOffers()
	if err != nil {
		return nil, err
	}

	var result []models.ProductOfferResponse

	for _, offer := range offers {

		result = append(result, models.ProductOfferResponse{
			ID:                 offer.ID,
			ProductID:          offer.ProductID,
			ProductName:        offer.Product.Name,
			OfferName:          offer.OfferName,
			DiscountPercentage: offer.DiscountPercentage,
			IsActive:           offer.IsActive,
		})
	}

	return result, nil
}
func DeleteProductOffer(offerID uint) error {

	return repository.DeleteProductOffer(offerID)
}

//category offer
func CreateCategoryOffer(req models.CreateCategoryOfferRequest) error {

	_, err := repository.GetCategoryByID(req.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return errors.New("invalid start date")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return errors.New("invalid end date")
	}

	if endDate.Before(startDate) {
		return errors.New("end date must be after start date")
	}

	if req.DiscountPercentage <= 0 || req.DiscountPercentage > 90 {
		return errors.New("discount percentage must be between 1 and 90")
	}

	offer := &domain.CategoryOffer{
		CategoryID:         req.CategoryID,
		OfferName:          req.OfferName,
		DiscountPercentage: req.DiscountPercentage,
		StartDate:          startDate,
		EndDate:            endDate,
		IsActive:           true,
	}

	return repository.CreateCategoryOffer(offer)
}

func GetAllCategoryOffers() ([]models.CategoryOfferResponse, error) {

	offers, err := repository.GetAllCategoryOffers()
	if err != nil {
		return nil, err
	}

	var result []models.CategoryOfferResponse

	for _, offer := range offers {

		result = append(result, models.CategoryOfferResponse{
			ID:                 offer.ID,
			CategoryID:         offer.CategoryID,
			CategoryName:       offer.Category.Name,
			OfferName:          offer.OfferName,
			DiscountPercentage: offer.DiscountPercentage,
			StartDate:          offer.StartDate.Format("2006-01-02"),
			EndDate:            offer.EndDate.Format("2006-01-02"),
			IsActive:           offer.IsActive,
		})
	}

	return result, nil
}

func DeleteCategoryOffer(id uint) error {
	return repository.DeleteCategoryOffer(id)
}

func GetReferralDetails(userID uint) (*models.ReferralResponse, error) {

	user, err := repository.GetUserByID(userID)

	if err != nil {
		return nil, err
	}

	return &models.ReferralResponse{
		ReferralCode: user.ReferralCode,
		ReferralLink: "https://shoego.com/signup?ref=" + user.ReferralCode,
	}, nil
}