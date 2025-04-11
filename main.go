package main

import (
	"customersales/config"
	"customersales/internal/handlers"
	"customersales/internal/utils"
	database "customersales/migrations"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.OpenFile("./log/log"+time.Now().Format("02012006.15.04.05.000000000")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error open file ", err.Error())
	}
	defer file.Close()
	log.SetOutput(file)
	err = database.ConnectDatabase()
	if err != nil {
		log.Fatal("Error Occur in DB Connection : ", err)
	}
	// go Schedular()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Page is unused so please contact to admin or use correct endpoint"})
	})

	router.GET("/refersh", handlers.DataRefershAPI)

	revenuGroup := router.Group("/revenue")
	revenuGroup.GET("/total/:s_date/:e_date", handlers.GetTotalRevenu)
	revenuGroup.GET("/by-product", handlers.GetRevenueByProduct)

	port := config.GetConfig().Service.Port
	router.Run(fmt.Sprintf(":%d", port))
}
func init() {
	config.LoadGlobalConfig("toml/config.toml")
}
func Schedular() {
	log := new(utils.Logger)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Log(utils.INFO, "Data refersh ")
			err := handlers.DataRefersh(log, "csv/sample.csv")
			if err != nil {
				log.Log(utils.ERROR, "Error Occur in schedular :", err.Error())
			}
		}
	}
}
