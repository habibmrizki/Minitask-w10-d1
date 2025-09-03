package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var users = make(map[string]User)

var (
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	numberRegex      = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[!@#$%^&*()_+=\-<>?|~{}]`)
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error"`
	UserID  int    `json:"user_id,omitempty"`
}

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("format email tidak valid")
	}
	return nil
}

// ValidateUserCredentials memeriksa validitas password
func ValidateUserCredentials(user User) error {
	if !lowercaseRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu huruf kecil")
	}

	if !uppercaseRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu huruf besar")
	}

	if !numberRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu angka")
	}

	if !specialCharRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu karakter khusus")
	}

	return nil
}

func RegisterHandler(ctx *gin.Context) {
	var newUser User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Gagal mengikat data. Pastikan format JSON benar dan semua field terisi.",
		})
		return
	}

	// Memanggil fungsi validasi email
	if err := ValidateEmail(newUser.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Email tidak valid: " + err.Error(),
		})
		return
	}

	// Memanggil fungsi validasi password
	if err := ValidateUserCredentials(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Kata sandi tidak memenuhi kriteria keamanan: " + err.Error(),
		})
		return
	}

	if _, found := users[newUser.Email]; found {
		ctx.JSON(http.StatusConflict, Response{
			Success: false,
			Message: "Email sudah terdaftar",
		})
		return
	}

	newUser.ID = len(users) + 1

	users[newUser.Email] = newUser

	fmt.Println("Pengguna baru terdaftar:", newUser)

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Pengguna berhasil terdaftar",
		UserID:  newUser.ID,
	})
}

func LoginHandler(ctx *gin.Context) {
	var loginUser User

	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Gagal mengikat data. Pastikan format JSON benar dan semua field terisi.",
		})
		return
	}

	storedUser, found := users[loginUser.Email]
	if !found || storedUser.Password != loginUser.Password {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Email atau password tidak valid",
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login berhasil",
		UserID:  storedUser.ID,
	})
}

func main() {
	router := gin.Default()

	api := router.Group("/auth")
	{
		api.POST("/register", RegisterHandler)
		api.POST("/login", LoginHandler)

	}

	router.Run()
}
