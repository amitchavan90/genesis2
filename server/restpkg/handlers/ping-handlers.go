package handlers

import (
	"net/http"
)

// Ping test
func Ping(w http.ResponseWriter, r *http.Request) {
	// cookie := r.Header.Get("cookie")
	// fmt.Println(cookie)

	w.Write([]byte("ping"))
}

// King test
func King(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("kong"))
}
