package handlers

import (
	"customersales/internal/models"
	"customersales/internal/utils"
	database "customersales/migrations"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRevenueByProduct(ctx *gin.Context) {
	log := new(utils.Logger)
	log.SetSid(ctx.Request)
	log.Log(utils.INFO, "GetRevenueByProduct (+)")
	defer log.Log(utils.INFO, "GetRevenueByProduct (-)")

	sDate, err := time.Parse("2006-01-02", ctx.Query("s_date"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid start date"})
		return
	}
	eDate, err := time.Parse("2006-01-02", ctx.Query("e_date"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid end date"})
		return
	}

	var results []models.RevenueResult

	err = database.GDB.Table("order_items").
		Select("products.name as product_name, SUM(order_items.unit_price * order_items.quantity_sold * (1 - order_items.discount)) as revenue").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("orders.date_of_sale BETWEEN ? AND ?", sDate, eDate).
		Group("products.name").
		Order("revenue DESC").
		Scan(&results).Error

	if err != nil {
		log.Log(utils.ERROR, "GR005", err.Error())
		ctx.JSON(500, gin.H{"status": "E", "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "S",
		"message": "Revenue by product",
		"data": gin.H{
			"start":    sDate.Format("2006-01-02"),
			"end":      eDate.Format("2006-01-02"),
			"revenues": results,
		},
	})
}
