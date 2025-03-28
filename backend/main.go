package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/scale-from-0-to-1-million/api/controller"
	"github.com/IgorBrizack/scale-from-0-to-1-million/api/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env.")
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	cacheDB := database.NewRedisClient()
	database := database.NewDatabase()
	masterDB := database.MasterDB()
	slaveDB := database.SlaveDB()

	userController := controller.NewController(cacheDB, masterDB, slaveDB)

	router := gin.Default()

	router.GET("/users", userController.GetUsers)
	router.POST("/users", userController.CreateUser)

	fmt.Printf("Running on %s\n", port)
	router.Run(":" + port)
}
