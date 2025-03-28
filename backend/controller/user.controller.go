package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IgorBrizack/scale-from-0-to-1-million/api/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserController struct {
	cacheDB  *redis.Client
	masterDB *gorm.DB
	slaveDB  *gorm.DB
}

func NewController(cacheDB *redis.Client, masterDB, slaveDB *gorm.DB) *UserController {
	return &UserController{
		cacheDB:  cacheDB,
		masterDB: masterDB,
		slaveDB:  slaveDB,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uc.masterDB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Users created successfully",
		"user":    user,
	})
}

func (uc *UserController) GetUsers(c *gin.Context) {
	cacheKey := "users_cache"
	ctx := context.Background()

	cachedUsers, err := uc.cacheDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []model.User
		if err := json.Unmarshal([]byte(cachedUsers), &users); err != nil {
			log.Println("Error unmarshaling cached users:", err)
			c.JSON(500, gin.H{"error": "Error unmarshaling cached users"})
			return
		}

		c.JSON(200, gin.H{
			"source": "cache",
			"users":  users,
		})
		return
	}

	var users []model.User
	if err := uc.slaveDB.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error fetching users from database"})
		return
	}

	cachedData, err := json.Marshal(users)
	if err != nil {
		log.Println("Error marshaling users:", err)
		c.JSON(500, gin.H{"error": "Error serializing users"})
		return
	}

	err = uc.cacheDB.Set(ctx, cacheKey, cachedData, 10*time.Second).Err()
	if err != nil {
		log.Println("Error caching users:", err)
	}

	c.JSON(200, gin.H{
		"source": "database",
		"users":  users,
	})
}
