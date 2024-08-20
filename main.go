package main

import (
	"fmt"
	"log"
	"time"
	"tugaspagik/controllers"
	"tugaspagik/middlewares"
	"tugaspagik/models"
	"tugaspagik/repositories"
	"tugaspagik/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var sessions = make(map[string]time.Time)

func main() {
	dsn := "egiwira:12345@tcp(127.0.0.1:3306)/testsiank?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	db.AutoMigrate(&models.User{})
	fmt.Println("Database berhasil terkoneksi")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-16388.c334.asia-southeast2-1.gce.redns.redis-cloud.com:16388", // Alamat Redis server
		Password: "RlpxlJnI8lBnnVVFAB7LOSp5UQJ6vNrJ",                                   // Kosongkan jika Redis tidak memerlukan password
		DB:       0,                                                                    // Database yang digunakan (0 adalah default database)
	})

	// Cek koneksi Redis
	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis!")

	// Seed admin user
	admin := models.User{Username: "admin", Password: "123"}
	if err := db.FirstOrCreate(&admin, models.User{Username: "admin"}).Error; err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, rdb) // Menambahkan Redis client ke dalam UserService
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	r := gin.Default()

	// Routes untuk login dan logout
	r.POST("/login", authController.Login)
	r.POST("/logout", authController.Logout)

	// Middleware untuk melindungi routes yang memerlukan otentikasi
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware()) // Middleware tanpa parameter
	protected.GET("/users", userController.GetAllUsers)
	protected.GET("/users/:id", userController.GetUserByID)
	protected.POST("/users", userController.CreateUser)
	protected.PUT("/users/:id", userController.UpdateUser)
	//protected.DELETE("/users/:id", userController.DeleteUser) // Di-comment karena tidak digunakan

	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
