package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthHandler struct {
	Users  UserStore
	Secret []byte
	TTL    time.Duration
	Now    func() time.Time
}

func NewAuthHandler(users UserStore, secret []byte, ttl time.Duration) *AuthHandler {
	return &AuthHandler{
		Users:  users,
		Secret: secret,
		TTL:    ttl,
		Now:    time.Now}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(requestIDKey).(string)

	fmt.Println(1)
	answerUnauthorized := func() {
		writeError(w, http.StatusUnauthorized, ErrorMessageJSON{
			ErrorCode: "NOT_AUTHORIZED",
			Message:   "Username or password is wrong",
			RequestID: requestID,
		})
	}
	fmt.Println(2)
	answerInternalError := func(msg string) {
		writeError(w, http.StatusInternalServerError, ErrorMessageJSON{
			ErrorCode: "INTERNAL_SERVER_ERROR",
			Message:   msg,
			RequestID: requestID,
		})
	}

	fmt.Println(3)
	submitted := UserLogin{}
	if err := json.NewDecoder(r.Body).Decode(&submitted); err != nil {
		answerUnauthorized()
		return
	}

	fmt.Println(submitted)
	user, err := h.Users.GetByUsername(r.Context(), submitted.Username)
	if err != nil {
		answerUnauthorized()
		return
	}
	fmt.Println(5)
	if err := VerifyPassword(user.Password, submitted.Password); err != nil {
		answerUnauthorized()
		return
	}

	fmt.Println(56)
	token, err := IssueToken(user, h.Secret, h.TTL, h.Now())
	if err != nil {
		answerInternalError("Token signing failed")
		return
	}

	fmt.Println(56)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.TTL.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writeJSON(w, http.StatusOK, requestID, map[string]string{"status": "ok"})
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(requestIDKey).(string)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	writeJSON(w, http.StatusOK, requestID, map[string]string{"status": "ok"})
}
func (h *AuthHandler) WhoAmI(r *http.Request) (*AppResponse, *AppError) {
	claims := r.Context().Value(userContextKey).(UserContext)
	return newAppResponse(http.StatusOK, claims), nil
}

func (h *AuthHandler) IsAuthenticated(r *http.Request) (*AppResponse, *AppError) {
	return newAppResponse(http.StatusNoContent, nil), nil
}
