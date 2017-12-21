package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var client *github.Client
var ctx context.Context
var appDir = path.Join(os.Getenv("HOME"), "Library", "Application Support", "GoGitCli")

func main() {

	setToken()

	// initialize()
	// getKey()
	// checkArgs()
	// createRepo(os.Args[1:2][0])
}

func setToken() {
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		os.Mkdir(appDir, 0744)
	}
	token, err := ioutil.ReadFile(path.Join(appDir, "apiToken"))
	if os.IsNotExist(err) {
		// TODO get token thorugh OAuth
		log.Fatal("Get token")
	}
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(token)})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func checkArgs() {
	if len(os.Args[1:]) < 1 {
		log.Fatal("No arguments given, we need at least one for the name")
	}
	// return nil
}

func createRepo(name string) {
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

func check(err error) {
	if err != nil {
		// fmt.Println(err)
		// panic(err)
		log.Fatal(err)
	}
}
