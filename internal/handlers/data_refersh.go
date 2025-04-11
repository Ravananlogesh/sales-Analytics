package handlers

import (
	"customersales/internal/models"
	"customersales/internal/utils"
	database "customersales/migrations"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DataRefershAPI(ctx *gin.Context) {
	log := new(utils.Logger)
	log.SetSid(ctx.Request)
	log.Log(utils.INFO, "DataRefershAPI (+)")

	err := DataRefersh(log, "csv/sample.csv")
	if err != nil {
		log.Log(utils.ERROR, "DR001", err.Error())
		ctx.JSON(500, gin.H{"status": "E", "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "S", "message": "refresh successful"})
	log.Log(utils.INFO, "DataRefershAPI (-)")
}

func DataRefersh(log *utils.Logger, path string) error {
	log.Log(utils.INFO, "DataRefersh (+)")

	file, err := os.Open(path)
	if err != nil {
		log.Log(utils.ERROR, "DR001", err.Error())
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		log.Log(utils.ERROR, "DR002", err.Error())
		return err
	}

	for _, row := range rows[1:] {
		orderID := row[0]
		productID := row[1]
		customerID := row[2]

		quantity, err := strconv.Atoi(row[7])
		if err != nil {
			log.Log(utils.ERROR, "DR003", err.Error())
			return err
		}

		unitPrice, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR004", err.Error())
			return err
		}

		discount, err := strconv.ParseFloat(row[9], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR005", err.Error())
			return err
		}

		shippingCost, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			log.Log(utils.ERROR, "DR006", err.Error())
			return err
		}

		dateOfSale, err := time.Parse("2006-01-02", row[6])
		if err != nil {
			log.Log(utils.ERROR, "DR007", err.Error())
			return err
		}

		customer := models.Customer{CustomerID: customerID}
		err = database.GDB.Where("customer_id = ?", customerID).FirstOrCreate(&customer, models.Customer{
			Name:    row[12],
			Email:   row[13],
			Address: row[14],
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR008", err.Error())
			return err
		}

		product := models.Product{ProductID: productID}
		err = database.GDB.Where("product_id = ?", productID).FirstOrCreate(&product, models.Product{
			Name:     row[3],
			Category: row[4],
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR009", err.Error())
			return err
		}

		order := models.Order{OrderID: orderID}
		err = database.GDB.Where("order_id = ?", orderID).FirstOrCreate(&order, models.Order{
			CustomerID:   customerID,
			Region:       row[5],
			DateOfSale:   dateOfSale,
			PaymentType:  row[11],
			ShippingCost: shippingCost,
		}).Error
		if err != nil {
			log.Log(utils.ERROR, "DR010", err.Error())
			return err
		}

		orderItem := models.OrderItem{
			OrderID:      orderID,
			ProductID:    productID,
			QuantitySold: quantity,
			UnitPrice:    unitPrice,
			Discount:     discount,
		}

		err = database.GDB.Create(&orderItem).Error
		if err != nil {
			log.Log(utils.ERROR, "DR011", err.Error())
			return err
		}
	}

	log.Log(utils.INFO, "CSV data processed and refreshed successfully")
	log.Log(utils.INFO, "DataRefersh (-)")
	return nil
}
