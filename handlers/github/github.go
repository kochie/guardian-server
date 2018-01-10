package github

import (
    "fmt"
    "net/http"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
//    "golang.org/x/oauth2/github"
    "context"
    "log"
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

func CreateConnection() {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: ""},
    )
    tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)

    // list all repositories for the authenticated user
    repos, _, err := client.Repositories.List(ctx, "", nil)
    handleError(err)

    fmt.Println(repos)
}

func handleError(err error) {
    log.Fatal(err)
}
