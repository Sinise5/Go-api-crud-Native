package controllers

import (
	"encoding/json"
	"log"
	"myapp/config"
	"myapp/models"
	"myapp/utils"
	"net/http"
	"regexp"
	"strings"
)

// Response struct untuk konsistensi format response
type Response struct {
	Message string `json:"message"`
}

// RegisterUser - Mendaftarkan pengguna baru
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi username
	user.Username = strings.TrimSpace(user.Username)
	if len(user.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters", http.StatusBadRequest)
		return
	}

	// Validasi password
	if !isValidPassword(user.Password) {
		http.Error(w, "Password must be at least 8 characters, including uppercase, lowercase, number, and special character", http.StatusBadRequest)
		return
	}

	// Cek apakah username sudah ada
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
	if err != nil {
		log.Printf("Error checking username existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Simpan user ke database
	_, err = config.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hashedPassword)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "User registered successfully"})
}

// Validasi username (alphanumeric, min 3 karakter)
func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,}$`)
	return re.MatchString(username)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input tidak boleh kosong
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Validasi format username
	if !isValidUsername(user.Username) {
		http.Error(w, "Invalid username format. Min 3 alphanumeric characters.", http.StatusBadRequest)
		return
	}

	// Validasi format password
	if !isValidPassword(user.Password) {
		http.Error(w, "Invalid password format. Min 8 chars with uppercase, lowercase, and a number.", http.StatusBadRequest)
		return
	}

	// Cek apakah user ada di database
	row := config.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", user.Username)
	var storedUser models.User
	err = row.Scan(&storedUser.ID, &storedUser.Password)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Cek apakah password cocok
	if !utils.CheckPasswordHash(user.Password, storedUser.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Generate token jika berhasil login
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Response jika berhasil login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token, "message": "Login successful"})
}

// isValidPassword - Validasi password dengan regex
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}
