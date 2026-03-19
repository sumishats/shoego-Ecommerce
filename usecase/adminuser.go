package usecase

import (
	"errors"
	"math"
	"shoego/models"
	"shoego/repository"
	"strconv"
	"time"
)

func GetAdminUsers(pageStr string, limitStr string, search string) (*models.AdminUserListResponse, error) {
	page := 1
	limit := 10

	//convert string to int 
	if pageStr != "" { 
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	//how mant row to skip
	offset := (page - 1) * limit 

	users, err := repository.GetUsers(search, limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := repository.CountUsers(search)
	if err != nil {
		return nil, err
	}

	//calculate total page 
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit))) 

	//user response 
	var userResponses []models.AdminUserResponse

	for _, user := range users {
		userResponses = append(userResponses, models.AdminUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Blocked:   user.Blocked,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	return &models.AdminUserListResponse{
		Users:      userResponses,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func BlockUser(id uint) error {
	user, err := repository.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Blocked {
		return errors.New("user already blocked")
	}

	return repository.BlockUser(id)
}

func UnblockUser(id uint) error {
	user, err := repository.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.Blocked {
		return errors.New("user already unblocked")
	}

	return repository.UnblockUser(id)
}