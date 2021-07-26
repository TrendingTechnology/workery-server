package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/utils"
)

// To run this API, try running in your console:
// $ http post 127.0.0.1:5000/api/v1/register email="fherbert@dune.com" password="the-spice-must-flow" name="Frank Herbert"
func (h *Controller) registerEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize our array which will store all the results from the remote server.
	var requestData models.RegisterRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// // For debugging purposes, print our output so you can see the code working.
	// fmt.Println(requestData.Name)
	// fmt.Println(requestData.Email)
	// fmt.Println(requestData.Password)

	// Lookup the email and if it is not unique we need to generate a `400 Bad Request` response.
	if userFound, _ := h.UserRepo.CheckIfExistsByEmail(ctx, requestData.Email); userFound {
		http.Error(w, "Email alread exists", http.StatusBadRequest)
		return
	}

	// Secure our password.
	passwordHash, err := utils.HashPassword(requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := &models.User{
		Uuid:         uuid.NewString(),
		TenantId:     requestData.TenantId,
		Email:        requestData.Email,
		FirstName:    requestData.FirstName,
		LastName:     requestData.LastName,
		PasswordHash: passwordHash,
		State:        1,
		Role:         4,
		Timezone:     "utc",
		CreatedTime:  time.Now(),
		ModifiedTime: time.Now(),
	}

	// Save our new user account.
	if err := h.UserRepo.Insert(ctx, m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate our response.
	responseData := models.RegisterResponse{
		Message: "You have successfully registered an account.",
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// To run this API, try running in your console:
// $ http post 127.0.0.1:5000/api/v1/login email="fherbert@dune.com" password="the-spice-must-flow"
func (h *Controller) loginEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	var requestData models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// // For debugging purposes, print our output so you can see the code working.
	// fmt.Println(requestData.Email)
	// fmt.Println(requestData.Password)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	user, err := h.UserRepo.GetByEmail(ctx, requestData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Incorrect email or password", http.StatusBadRequest)
		return
	}

	// Verify the inputted password and hashed password match.
	passwordMatch := utils.CheckPasswordHash(requestData.Password, user.PasswordHash)
	if passwordMatch == false {
		http.Error(w, "Incorrect email or password", http.StatusBadRequest)
		return
	}

	// Start our session.
	sessionExpiryTime := time.Hour * 24 * 7 // 1 week
	sessionUuid := uuid.NewString()
	err = h.SessionManager.SaveUser(ctx, sessionUuid, user, sessionExpiryTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate our JWT token.
	accessToken, refreshToken, err := utils.GenerateJWTTokenPair(h.SecretSigningKeyBin, sessionUuid, sessionExpiryTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Finally return success.
	responseData := models.LoginResponse{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Role:         user.Role,
		TenantId:     user.TenantId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For debugging purposes only.
	log.Println("loginEndpoint | Response:", responseData)
}

// To run this API, try running in your console:
// $ http post 127.0.0.1:5000/api/v1/refresh-token value="xxx"
func (h *Controller) postRefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var requestData models.RefreshTokenRequest

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For debugging purposes, print our output so you can see the code working.
	log.Println(requestData.Value)

	ctx := r.Context()

	// Verify our refresh token.
	sessionUuid, err := utils.ProcessJWTToken(h.SecretSigningKeyBin, requestData.Value)
	if err != nil {
		http.Error(w, "Unauthorized - refresh token expired or invalid", http.StatusUnauthorized)
		return
	}

	// Lookup our user profile in the session or return 500 error.
	user, err := h.SessionManager.GetUser(ctx, sessionUuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate our JWT token.
	newSessionUuid := uuid.NewString()
	newSessionExpiryTime := time.Hour * 24 * 7 // 1 week
	accessToken, refreshToken, err := utils.GenerateJWTTokenPair(h.SecretSigningKeyBin, newSessionUuid, newSessionExpiryTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save our new session.
	err = h.SessionManager.SaveUser(ctx, newSessionUuid, user, newSessionExpiryTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Finally return success.
	responseData := models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SPECIAL THANKS:
// [1][2]: Learned from:
// a. https://blog.golang.org/json
// b. https://stackoverflow.com/questions/21197239/decoding-json-using-json-unmarshal-vs-json-newdecoder-decode
