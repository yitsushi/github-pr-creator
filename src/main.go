package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)


var (
	sourceOwner   = flag.String("source-owner", "", "Name of the owner (user or org) of the repo to create the pull request.")
	sourceRepo    = flag.String("source-repo", "", "Name of repo to create the pull request.")
	sourceBranch  = flag.String("source-branch", "", "Name of the branch to create the pull request.")
	prRepoOwner   = flag.String("merge-repo-owner", "", "Name of the owner (user or org) of the repo to create the PR against. If not specified, the value of the `-source-owner` flag will be used.")
	prRepo        = flag.String("merge-repo", "", "Name of repo to create the PR against. If not specified, the value of the `-source-repo` flag will be used.")
	prBranch      = flag.String("merge-branch", "master", "Name of branch to create the PR against (the one you want to merge your branch in via the PR).")
	prSubject     = flag.String("pr-title", "", "Title of the pull request. If not specified, no pull request will be created.")
	prDescription = flag.String("pr-text", "", "Text to put in the description of the pull request.")

	client *github.Client
	ctx    = context.Background()
)

func createPR() (err error) {
	if *prSubject == "" {
		return errors.New("missing `-pr-title` flag; skipping PR creation")
	}

	if *prRepoOwner != "" && *prRepoOwner != *sourceOwner {
		*sourceBranch = fmt.Sprintf("%s:%s", *sourceOwner, *sourceBranch)
	} else {
		prRepoOwner = sourceOwner
	}

	if *prRepo == "" {
		prRepo = sourceRepo
	}

	newPR := &github.NewPullRequest{
		Title:               prSubject,
		Head:                sourceBranch,
		Base:                prBranch,
		Body:                prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, *prRepoOwner, *prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func init() {
	flag.Parse()

	if *sourceOwner == "" || *sourceRepo == "" || *sourceBranch == "" || *prSubject == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-source-branch`, `-pr-title`")
	}

	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_OAUTH_TOKEN")})
	client = github.NewClient(oauth2.NewClient(ctx, token))
}

func main() {
	if err := createPR(); err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
}
