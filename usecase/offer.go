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

func GetAllProductOffers(page, limit int) (models.PaginatedResponse, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offers, totalCount, err := repository.GetAllProductOffers(page, limit)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	var result []models.ProductOfferResponse

	for _, offer := range offers {
		result = append(result, models.ProductOfferResponse{
			ID:                 offer.ID,
			ProductID:          offer.ProductID,
			ProductName:        offer.Product.Name,
			OfferName:          offer.OfferName,
			DiscountPercentage: offer.DiscountPercentage,
			StartDate:          offer.StartDate.Format("2006-01-02"),
			EndDate:            offer.EndDate.Format("2006-01-02"),
			IsActive:           offer.IsActive,
		})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return models.PaginatedResponse{
		Offers:     result,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}
func DeleteProductOffer(offerID uint) error {

	return repository.DeleteProductOffer(offerID)
}

// category offer
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

func GetAllCategoryOffers(page, limit int) (models.PaginatedCategoryOfferResponse, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offers, totalCount, err := repository.GetAllCategoryOffers(page, limit)
	if err != nil {
		return models.PaginatedCategoryOfferResponse{}, err
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

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return models.PaginatedCategoryOfferResponse{
		Offers:     result,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
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
func UpdateProductOffer(offerID uint, req models.UpdateProductOfferRequest) error {

	offer, err := repository.GetProductOfferByID(offerID)
	if err != nil {
		return errors.New("offer not found")
	}

	_, err = repository.GetProductByID(req.ProductID)
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

	offer.ProductID = req.ProductID
	offer.OfferName = req.OfferName
	offer.DiscountPercentage = req.DiscountPercentage
	offer.StartDate = startDate
	offer.EndDate = endDate

	return repository.UpdateProductOffer(offer)
}
func UpdateCategoryOffer(offerID uint, req models.UpdateCategoryOfferRequest) error {

	offer, err := repository.GetCategoryOfferByID(offerID)
	if err != nil {
		return errors.New("category offer not found")
	}

	_, err = repository.GetCategoryByID(req.CategoryID)
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

	offer.CategoryID = req.CategoryID
	offer.OfferName = req.OfferName
	offer.DiscountPercentage = req.DiscountPercentage
	offer.StartDate = startDate
	offer.EndDate = endDate

	return repository.UpdateCategoryOffer(offer)
}
