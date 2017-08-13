package user

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hii there, I love %s!", r.URL.Path[1:])
}
