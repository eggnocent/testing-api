package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"tugaspagik/models"
	"tugaspagik/repositories"
	"tugaspagik/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var sessions = make(map[string]time.Time)

func main() {
	dsn := "egiwira:12345@tcp(127.0.0.1:3306)/testuser?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	db.AutoMigrate(&models.User{})
	fmt.Println("Database berhasil terkoneksi")

	// Seed admin user
	admin := models.User{Username: "admin", Password: "123"}
	if err := db.FirstOrCreate(&admin, models.User{Username: "admin"}).Error; err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	// Terminal-based login
	var username, password string
	fmt.Println("Silahkan login terlebih dahulu,")
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	user, err := authService.Authenticate(username, password)
	if err != nil {
		fmt.Println("Login gagal:", err)
		return
	}

	fmt.Println("Login berhasil, selamat datang", user.Username)

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Lihat tabel user")
		fmt.Println("2. Buat user baru")
		fmt.Println("3. Update user")
		fmt.Println("4. Hapus user")
		fmt.Println("5. Logout")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			users, err := userService.GetAllUsers()
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Daftar User:")
				for _, u := range users {
					fmt.Printf("ID: %d, Username: %s\n", u.ID, u.Username)
				}
			}
		case 2:
			var newUser models.User
			fmt.Print("Masukkan username baru: ")
			fmt.Scanln(&newUser.Username)
			fmt.Print("Masukkan password baru: ")
			fmt.Scanln(&newUser.Password)
			err := userService.CreateUser(newUser)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("User berhasil dibuat")
			}
		case 3:
			var id string
			var updateUser models.User
			fmt.Print("Masukkan ID user yang ingin diupdate: ")
			fmt.Scanln(&id)
			fmt.Print("Masukkan username baru: ")
			fmt.Scanln(&updateUser.Username)
			fmt.Print("Masukkan password baru: ")
			fmt.Scanln(&updateUser.Password)
			idUint, _ := strconv.ParseUint(id, 10, 32)
			updateUser.ID = uint(idUint)
			err := userService.UpdateUser(updateUser)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("User berhasil diupdate")
			}
		case 4:
			var id string
			fmt.Print("Masukkan ID user yang ingin dihapus: ")
			fmt.Scanln(&id)
			idUint, _ := strconv.ParseUint(id, 10, 32)
			err := userService.DeleteUser(uint(idUint))
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("User berhasil dihapus")
			}
		case 5:
			fmt.Println("Logout berhasil")
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}
