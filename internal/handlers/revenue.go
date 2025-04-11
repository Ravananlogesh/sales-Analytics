package handlers

import (
	"customersales/internal/models"
	"customersales/internal/utils"
	database "customersales/migrations"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTotalRevenu(ctx *gin.Context) {
	log := new(utils.Logger)
	log.Log(utils.INFO, "GetTotalRevenu (+)")
	defer log.Log(utils.INFO, "GetTotalRevenu (-)")

	sDateStr := ctx.DefaultQuery("s_date", ctx.Param("s_date"))
	eDateStr := ctx.DefaultQuery("e_date", ctx.Param("e_date"))

	if sDateStr == "" || eDateStr == "" {
		log.Log(utils.ERROR, "GR001", "start date and end date are required")
		utils.JSONErrorResponse(ctx, 400, errors.New("start date and end date are required"))
		return
	}

	sDate, err := time.Parse("2006-01-02", sDateStr)
	if err != nil {
		log.Log(utils.ERROR, "GR002", err.Error())
		utils.JSONErrorResponse(ctx, 400, errors.New("invalid start date format, use YYYY-MM-DD"))
		return
	}

	eDate, err := time.Parse("2006-01-02", eDateStr)
	if err != nil {
		log.Log(utils.ERROR, "GR003", err.Error())
		utils.JSONErrorResponse(ctx, 400, errors.New("invalid end date format, use YYYY-MM-DD"))
		return
	}

	if sDate.After(eDate) {
		log.Log(utils.ERROR, "GR004", "start date is greater than end date")
		utils.JSONErrorResponse(ctx, 400, errors.New("start date cannot be after end date"))
		return
	}

	var total struct {
		Amount float64 `json:"amount"`
	}

	err = database.GDB.Model(&models.OrderItem{}).
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("orders.date_of_sale BETWEEN ? AND ?", sDate, eDate).
		Select("SUM(order_items.unit_price * order_items.quantity_sold * (1 - order_items.discount)) as amount").
		Scan(&total).Error

	if err != nil {
		log.Log(utils.ERROR, "GR005", err.Error())
		utils.JSONErrorResponse(ctx, 500, errors.New("failed to calculate total revenue"))
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "S",
		"message": "total revenue",
		"data": gin.H{
			"total_amount": total.Amount,
			"start_date":   sDate.Format("2006-01-02"),
			"end_date":     eDate.Format("2006-01-02"),
		},
	})
}
