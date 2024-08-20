package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"tugaspagik/models"
	"tugaspagik/repositories"
	"tugaspagik/services"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Alamat Redis server
		Password: "",               // Kosongkan jika Redis tidak memerlukan password
		DB:       0,                // Database yang digunakan (0 adalah default database)
	})

	// Cek koneksi Redis
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis!")

	dsn := "egiwira:12345@tcp(127.0.0.1:3306)/testsiank?charset=utf8mb4&parseTime=True&loc=Local"
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
	userService := services.NewUserService(userRepo, rdb) // Tambahkan Redis client di sini

	// Terminal-based login
	reader := bufio.NewReader(os.Stdin)
	var username, password string
	fmt.Println("Silahkan login terlebih dahulu,")
	fmt.Print("Username: ")
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username) // Trim untuk menghapus karakter newline
	fmt.Print("Password: ")
	password, _ = reader.ReadString('\n')
	password = strings.TrimSpace(password)

	user, err := authService.Authenticate(username, password)
	if err != nil {
		fmt.Println("Login gagal:", err)
		return
	}

	// Set waktu login sebagai waktu aktivitas terakhir di Redis
	sessionID := user.Username
	rdb.Set(sessionID, time.Now().Format(time.RFC3339), time.Minute)

	fmt.Println("Login berhasil, selamat datang", user.Username)

	for {
		// Ambil waktu terakhir dari Redis
		lastActivity, err := rdb.Get(sessionID).Result()
		if err != nil || time.Since(parseTime(lastActivity)) > time.Minute {
			fmt.Println("Sesi Anda telah berakhir karena tidak ada aktivitas selama 1 menit.")
			rdb.Del(sessionID) // Hapus sesi di Redis
			break
		}

		fmt.Println("\nMenu:")
		fmt.Println("1. Lihat tabel user")
		fmt.Println("2. Buat user baru")
		fmt.Println("3. Update user")
		fmt.Println("4. Logout")

		var choice int
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&choice)

		// Periksa kembali sebelum memperbarui aktivitas
		lastActivity, err = rdb.Get(sessionID).Result()
		if err != nil || time.Since(parseTime(lastActivity)) > time.Minute {
			fmt.Println("Sesi Anda telah berakhir karena tidak ada aktivitas selama 1 menit.")
			rdb.Del(sessionID) // Hapus sesi di Redis
			break
		}

		// Perbarui waktu aktivitas terakhir di Redis
		rdb.Set(sessionID, time.Now().Format(time.RFC3339), time.Minute)

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
			fmt.Print("Masukkan nama lengkap: ")
			newUser.FullName, _ = reader.ReadString('\n')
			newUser.FullName = strings.TrimSpace(newUser.FullName)
			fmt.Print("Masukkan username baru: ")
			newUser.Username, _ = reader.ReadString('\n')
			newUser.Username = strings.TrimSpace(newUser.Username)
			fmt.Print("Masukkan password baru: ")
			newUser.Password, _ = reader.ReadString('\n')
			newUser.Password = strings.TrimSpace(newUser.Password)
			fmt.Print("Masukkan email: ")
			newUser.Email, _ = reader.ReadString('\n')
			newUser.Email = strings.TrimSpace(newUser.Email)

			err := userService.CreateUser(newUser)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("User berhasil dibuat")
			}

		case 3:
			var id string
			fmt.Print("Masukkan ID user yang ingin diupdate: ")
			id, _ = reader.ReadString('\n')
			id = strings.TrimSpace(id)

			user, err := userService.GetUserByID(id)
			if err != nil {
				fmt.Println("Error: id tidak ditemukan")
				break
			}

			fmt.Print("Masukkan nama lengkap baru: ")
			user.FullName, _ = reader.ReadString('\n')
			user.FullName = strings.TrimSpace(user.FullName)
			fmt.Print("Masukkan username baru: ")
			user.Username, _ = reader.ReadString('\n')
			user.Username = strings.TrimSpace(user.Username)
			fmt.Print("Masukkan password baru: ")
			user.Password, _ = reader.ReadString('\n')
			user.Password = strings.TrimSpace(user.Password)
			fmt.Print("Masukkan email baru: ")
			user.Email, _ = reader.ReadString('\n')
			user.Email = strings.TrimSpace(user.Email)

			err = userService.UpdateUser(user)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("User berhasil diupdate")
			}

		case 4:
			fmt.Println("Logout berhasil")
			rdb.Del(sessionID) // Hapus sesi dari Redis saat logout
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}

// parseTime adalah helper function untuk mengonversi string waktu dari Redis ke time.Time
func parseTime(value string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, value)
	return parsedTime
}
