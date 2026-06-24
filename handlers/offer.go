package handlers

import (
	"net/http"
	"shoego/models"
	"shoego/response"
	"shoego/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProductOffer(c *gin.Context) {

	var req models.CreateProductOfferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid request body", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := usecase.CreateProductOffer(req)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to create offer", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "product offer created successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetAllProductOffers(c *gin.Context) {

	offers, err := usecase.GetAllProductOffers()
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to fetch offers", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "offers fetched successfully", offers, nil))
}
func DeleteProductOffer(c *gin.Context) {

	offerID64, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid offer id", nil, err.Error()))
		return
	}

	err = usecase.DeleteProductOffer(uint(offerID64))
	if err != nil {

		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to delete offer", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "offer deleted successfully", nil, nil))
}

// category offer
func CreateCategoryOffer(c *gin.Context) {

	var req models.CreateCategoryOfferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid request", nil, err.Error()))
		return
	}

	err := usecase.CreateCategoryOffer(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to create category offer", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "category offer created successfully", nil, nil))
}

func GetAllCategoryOffers(c *gin.Context) {

	offers, err := usecase.GetAllCategoryOffers()
	if err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "failed to fetch category offers", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "category offers fetched successfully", offers, nil)
	c.JSON(http.StatusOK, successRes)
}

func DeleteCategoryOffer(c *gin.Context) {

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "invalid offer id", nil, err.Error()))
		return
	}

	err = usecase.DeleteCategoryOffer(uint(id64))

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to delete category offer", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "category offer deleted successfully", nil, nil))
}

//referral

func GetReferralDetails(c *gin.Context) {

	userID := c.GetUint("user_id")

	data, err := usecase.GetReferralDetails(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ClientResponse(http.StatusBadRequest, "failed to get referral details", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "referral details fetched successfully", data, nil))
}
