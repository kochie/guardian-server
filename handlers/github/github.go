package github

import (
    "fmt"
    "net/http"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Disconnect(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hii there, I love %s!", r.URL.Path[1:])
}

func Addkey(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hii there, I love %s!", r.URL.Path[1:])
}
