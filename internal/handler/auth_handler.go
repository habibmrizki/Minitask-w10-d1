package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/day1/internal/models"
)

// RegisterHandler menangani pendaftaran pengguna
func RegisterHandler(ctx *gin.Context) {
	var newUser models.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Gagal mengikat data. Pastikan format JSON benar dan semua field terisi.",
		})
		return
	}

	if err := models.ValidateEmail(newUser.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Email tidak valid: " + err.Error(),
		})
		return
	}

	if err := models.ValidateUserCredentials(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Kata sandi tidak memenuhi kriteria keamanan: " + err.Error(),
		})
		return
	}

	users := models.GetUsersMap()
	if _, found := users[newUser.Email]; found {
		ctx.JSON(http.StatusConflict, models.Response{
			Success: false,
			Message: "Email sudah terdaftar",
		})
		return
	}

	newUser.ID = len(users) + 1
	users[newUser.Email] = newUser

	fmt.Println("Pengguna baru terdaftar:", newUser)

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Pengguna berhasil terdaftar",
		UserID:  newUser.ID,
	})
}

// LoginHandler menangani login pengguna
func LoginHandler(ctx *gin.Context) {
	var loginUser models.User

	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Gagal mengikat data. Pastikan format JSON benar dan semua field terisi.",
		})
		return
	}

	users := models.GetUsersMap()
	storedUser, found := users[loginUser.Email]
	if !found || storedUser.Password != loginUser.Password {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   "Email atau password tidak valid",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Login berhasil",
		UserID:  storedUser.ID,
	})
}
