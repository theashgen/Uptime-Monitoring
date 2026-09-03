package auth

import (
	"net/http"
	"os"
)

func SetCookie(w http.ResponseWriter, val string) {

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	}

	http.SetCookie(w, cookie)
}
