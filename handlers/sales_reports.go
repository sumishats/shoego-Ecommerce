package handlers

import (
	"net/http"
	"shoego/response"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func GetSalesReport(c *gin.Context) {

	filter := c.DefaultQuery("filter", "daily")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := usecase.GetSalesReport(filter, startDate, endDate)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "failed to generate sales report", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "sales report fetched successfully", data, nil)

	c.JSON(http.StatusOK, successRes)
}

func DownloadSalesReportPDF(c *gin.Context) {

	period := c.DefaultQuery("period", "daily")
	from := c.Query("from")
	to := c.Query("to")

	pdfBytes, err := usecase.DownloadSalesReportPDF(period, from, to)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError,"failed to generate sales report pdf",nil,err.Error(),)
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=sales-report.pdf")

	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
