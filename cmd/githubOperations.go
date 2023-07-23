package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

// Repo data
type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

var fileForAutoCommit = "commitTimeStamps.txt"

func editSampleFile(file string, data string) bool {
	fmt.Printf("\nEditing file for commit: %s\n", file)
	fHdl, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error in opening file to edit: ", err)
		return false
	}

	defer fHdl.Close()

	if _, err = fHdl.WriteString(data + "\n"); err != nil {
		fmt.Println("Error in writing to : ", err)
		return false
	}

	return true
}

func getTree(client *github.Client, ctx context.Context, ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(*sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	//fmt.Printf("srcOwner[%s], srcRepo [%s], objSHA[%s]", *sourceOwner, *sourceRepo, *ref.Object.SHA)
	//fmt.Println(entries)
	tree, _, err = client.Git.CreateTree(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA, entries)
	if err != nil {
		fmt.Println("Error in github client CreateTree: ", err)
	}
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty '-files parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	if err != nil {
		fmt.Println("Error in os.Readfile for file: ", err)
		fmt.Printf("File: %s\n", localFile)
	}
	return targetName, b, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(client *github.Client, ctx context.Context, ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, *sourceOwner, *sourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: authorName, Email: authorEmail}
	commit := &github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []github.Commit{*parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, *sourceOwner, *sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, *sourceOwner, *sourceRepo, ref, false)
	return err
}

func getRepository(client *github.Client, ctx context.Context) {
	repo, _, err := client.Repositories.Get(ctx, *sourceOwner, *sourceRepo)

	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}
	pack := &Package{
		FullName:    *repo.FullName,
		Description: *repo.Description,
		ForksCount:  *repo.ForksCount,
		StarsCount:  *repo.StargazersCount,
	}

	fmt.Printf("\n'%s' Repo Details:\n%+v\n", *sourceRepo, pack)
}

func createNewBranch(client *github.Client, ctx context.Context, newRefName string) (newRef *github.Reference, err error) {
	*prSubject = "Change Request for " + newRefName
	*prDescription = "Incremental additions for user story " + newRefName

	var baseRef *github.Reference
	baseRefName := "refs/heads/" + *baseBranch
	baseRef, _, err = client.Git.GetRef(ctx, *sourceOwner, *prRepo, baseRefName)
	if err != nil {
		fmt.Printf("Error in GetRef for %s: %v\n", baseRefName, err)
		return nil, err
	}
	newRef = &github.Reference{Ref: github.String("refs/heads/" + newRefName), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	if _, _, err := client.Git.CreateRef(ctx, *sourceOwner, *prRepo, newRef); err != nil {
		fmt.Printf("Error in CreateRef for %s: %v\n", newRefName, err)
		return nil, err
	}
	fmt.Printf("\nSuccessfully Created new Branch: %s\n", newRefName)
	return newRef, nil
}

func createPR(client *github.Client, ctx context.Context,
	newRef *github.Reference, appendData string,
	newRefName string) {

	//Edit a file in repo for commit
	editSampleFile(fileForAutoCommit, appendData)

	// if no source files given from command line, just include auto-edited
	// file to include in commit, else append to the list of files
	if *sourceFiles == "" {
		*sourceFiles = fileForAutoCommit
	} else {
		*sourceFiles += "," + fileForAutoCommit
	}

	// Create Commit
	tree, err := getTree(client, ctx, newRef)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(client, ctx, newRef, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	// PR Creation in newly created Branch
	*commitBranch = fmt.Sprintf("%s:%s", *sourceOwner, newRefName)

	newPR := &github.NewPullRequest{
		Title:               prSubject,
		Head:                commitBranch,
		Base:                prBranch,
		Body:                prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, *prRepoOwner, *prRepo, newPR)
	if err != nil {
		fmt.Printf("Error in CreatePR against %s, error: %v\n", newRefName, err)
	}

	fmt.Printf("\nPR created: %s\n", pr.GetHTMLURL())
}
