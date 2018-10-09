package command

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var rgx = regexp.MustCompile(`(\w+)/(.*)$`)

func OpenPullRequest(githubToken string, repo string, branch string, title string, body string) error {
	matches := rgx.FindAllStringSubmatch(repo, 1)
	if len(matches) == 0 || len(matches[0]) != 3 {
		return errors.New("could not parse repo owner and repo name")
	}

	owner := matches[0][1]
	repoName := matches[0][2]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	_, _, err := client.PullRequests.Create(ctx, owner, repoName, &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(branch),
		Base:                github.String("master"),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	})

	fmt.Printf("'%s' create-pr> : ", repo)

	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return err
	}

	fmt.Println("done")

	return nil
}
