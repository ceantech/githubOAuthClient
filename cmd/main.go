package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	sourceOwner   = flag.String("source-owner", "", "Name of the owner (user or org) of the repo to create the commit in.")
	sourceRepo    = flag.String("source-repo", "", "Name of repo to create the commit in.")
	commitMessage = flag.String("commit-message", "", "Content of the commit message.")
	commitBranch  = flag.String("commit-branch", "", "Name of branch to create the commit in. If it does not already exists, it will be created using the `base-branch` parameter")
	baseBranch    = flag.String("base-branch", "main", "Name of branch to create the `commit-branch` from.")
	prRepoOwner   = flag.String("merge-repo-owner", "", "Name of the owner (user or org) of the repo to create the PR against. If not specified, the value of the `-source-owner` flag will be used.")
	prRepo        = flag.String("merge-repo", "", "Name of repo to create the PR against. If not specified, the value of the `-source-repo` flag will be used.")
	prBranch      = flag.String("merge-branch", "main", "Name of branch to create the PR against (the one you want to merge your branch in via the PR).")
	prSubject     = flag.String("pr-title", "", "Title of the pull request. If not specified, no pull request will be created.")
	prDescription = flag.String("pr-text", "", "Text to put in the description of the pull request.")
	sourceFiles   = flag.String("files", "", `Comma-separated list of files to commit and their location.
The local file is separated by its target location by a semi-colon.
If the file should be in the same location with the same name, you can just put the file name and omit the repetition.
Example: README.md,main.go:github/examples/commitpr/main.go`)
	authorName  = flag.String("author-name", "", "Name of the author of the commit.")
	authorEmail = flag.String("author-email", "", "Email of the author of the commit.")
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(genIndexHtml(*authorName, *prRepo)))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := authInit()
	//fmt.Printf("Redirect URL from Github: %s\n", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// github_callback handler, invoked after successful user authorization
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	client, ctx := getAuthClient(w, r)
	if client == nil {
		log.Fatal("Unable to initialize Auth Client, exiting..")
	}

	// Get User login details
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("\nLogged in as GitHub user: %s\n", *user.Login)
	//fmt.Println(*user)

	*sourceOwner = *user.Login
	*sourceRepo = "gomods" //DEFAULT, need to take as i/p ideally
	*prRepoOwner = *user.Login
	*prRepo = "gomods" //DEFAULT
	*authorName = *user.Name
	*authorEmail = "dummy@gmail.com"

	// list repositories for authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Printf("client.Repositories.List() failed with '%s'\n", err)
	} else {
		fmt.Printf("\nRepos associated with user '%s':\n", *user.Name)
		for _, repo := range repos {
			fmt.Println(*repo.Name)
		}
	}

	// get Source Repo details
	getRepository(client, ctx)

	// Generate Current time reference in new branch name
	currTimeStr := fmt.Sprintf("%d", time.Now().Nanosecond())
	newRefName := "UserStory-" + currTimeStr

	*commitMessage = fmt.Sprintf("%s: initial changes for the feature", newRefName)

	// Create new branch
	newRef, err := createNewBranch(client, ctx, newRefName)
	if err != nil {
		fmt.Printf("Error while creating branch")
		return
	}

	createPR(client, ctx, newRef, currTimeStr, newRefName)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
	// Register Handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)

	// Invoked by github on successful user autherization
	http.HandleFunc("/github/callback", callbackHandler)

	// Start the web server
	fmt.Println("Starting WebApp on http://127.0.0.1:9080")
	log.Fatal(http.ListenAndServe(":9080", nil))
}
