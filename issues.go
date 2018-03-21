package main

import (
	"context"
	"sort"
	"time"

	"github.com/google/go-github/github"
)

type ByYouth []*github.Issue

func (is ByYouth) Len() int           { return len(is) }
func (is ByYouth) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
func (is ByYouth) Less(i, j int) bool { return is[i].GetNumber() > is[j].GetNumber() }

func FetchIssues(client *github.Client, ctx context.Context, cfg Config) ([]*github.Issue, error) {
	issues := []*github.Issue{}

	for _, uname := range cfg.Repository.Users {
		isss, _, err := client.Issues.ListByRepo(ctx, uname, cfg.Repository.Name, nil)
		if err != nil {
			return nil, err
		}
		issues = append(issues, isss...)
	}

	// filter issues
	zeroTime := time.Time{}
	toPublish := []*github.Issue{}
	for _, issue := range issues {
		if issue.GetClosedAt() != zeroTime {
			// Issue was closed so we don't publish it.
			continue
		}
		found := false
		for _, uname := range cfg.Repository.Users {
			if issue.User.GetLogin() == uname {
				// Only if the user is an allowed publisher.
				found = true
				break
			}
		}
		if found {
			toPublish = append(toPublish, issue)
		}
	}

	// Sort issues: newest first
	sort.Sort(ByYouth(issues))

	return issues, nil
}
