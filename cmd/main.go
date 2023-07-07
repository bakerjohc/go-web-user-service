package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	"github.com/bakerjohc/go-web-user-service/internal/users"
	"github.com/bakerjohc/go-web-user-service/common"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {

	db := common.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// test 1 to 1
	tx1 := db.Begin()
	userA := users.UserModel{
		Username: "test",
		Email:    "test@test.com",
		Bio:      "test",
		Image:    nil,
	}
	tx1.Save(&userA)
	tx1.Commit()
	fmt.Println(userA)

	r.Run() // listen and serve on 0.0.0.0:8080
}