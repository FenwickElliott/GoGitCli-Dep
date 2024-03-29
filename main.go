package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var client *github.Client
var ctx context.Context
var appDir = path.Join(os.Getenv("HOME"), "Library", "Application Support", "GoGitCli")

func main() {
	args := os.Args[1:]
	switch strings.ToLower(args[0]) {
	case "authorize":
		authorize()
	case "deauthorize":
		deauthorize()
	case "new":
		createRepo(os.Args[1:][1])
	default:
		fmt.Println("Not an option")
	}
}

func setToken() {
	token, err := ioutil.ReadFile(path.Join(appDir, "apiToken"))
	if os.IsNotExist(err) {
		log.Fatal("No API token found. Run: 'GoGitCli authorize'")
	}
	check(err)
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(token)})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func createRepo(name string) {
	setToken()
	repo := &github.Repository{
		Name: github.String(name),
	}
	_, _, err := client.Repositories.Create(ctx, "", repo)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully created ", name)
	}
}

func authorize() {
	fmt.Println("Authorizing, please check your browser")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		os.Mkdir(appDir, 0744)
	}
	done := make(chan bool)
	go serve(done)
	exec.Command("open", "https://github.com/login/oauth/authorize?client_id=974e6b9d6153b79b9fb9&redirect_uri=http://localhost:3456/catch&scope=repo&state=rabbits&allow_signup=true").Start()
	finished := <-done
	if finished {
		fmt.Println("Authorization Complete")
	}
}

func deauthorize() {
	err := os.Remove(path.Join(appDir, "apiToken"))
	check(err)
	fmt.Println("Successfully remove authorization")
}

func serve(done chan bool) {
	http.HandleFunc("/catch", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Thank you, GoGitCli can now access your GitHub account.\nYou may close this window.\n")
		response := r.URL.Query()
		if response["state"][0] == "rabbits" {
			exchangeCode(response["code"][0], done)
		} else {
			log.Fatal("Missmatched states!")
		}
	})
	http.ListenAndServe(":3456", nil)
}

func exchangeCode(code string, done chan bool) {
	clientSecret, err := ioutil.ReadFile(".client_secret")
	if err != nil {
		log.Fatal("Client_secret not present, you can't exchange code")
	}
	body := strings.NewReader("client_id=974e6b9d6153b79b9fb9&client_secret=" + string(clientSecret) + "&code=" + code)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", body)
	req.Header.Set("Accept", "application/json")
	check(err)
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	tokenBytes, err := ioutil.ReadAll(resp.Body)
	check(err)
	tokenResp := make(map[string]string)
	err = json.Unmarshal(tokenBytes, &tokenResp)
	check(err)

	err = ioutil.WriteFile(path.Join(appDir, "apiToken"), []byte(tokenResp["access_token"]), 0644)
	check(err)

	done <- true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
