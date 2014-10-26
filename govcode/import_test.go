package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestHelpWantedLabelCount(t *testing.T) {
	var issues []github.Issue
	issues = append(issues, makeIssue([]string{}))
	assertHelpWantedCount(t, issues, 0)
	issues = append(issues, makeIssue([]string{"help wanted"}))
	assertHelpWantedCount(t, issues, 1)
	issues = append(issues, makeIssue([]string{"help_wanted"}))
	assertHelpWantedCount(t, issues, 2)
	issues = append(issues, makeIssue([]string{"needs-help"}))
	assertHelpWantedCount(t, issues, 3)
	//ensure multiple labels that match only result in a single increment
	issues = append(issues, makeIssue([]string{"help wanted", "needs-help", "bug"}))
	assertHelpWantedCount(t, issues, 4)
	issues = append(issues, makeIssue([]string{"42"}))
	assertHelpWantedCount(t, issues, 4)
	issues = append(issues, makeIssue([]string{"tomorrow and tomorrow and tomorrow creeps in this petty place from day to day to the last syllable of recorded time;", "Life... is a tale told by an idiot, full of sound and fury, signifying nothing"}))
	assertHelpWantedCount(t, issues, 4)
	issues = append(issues, makeIssue([]string{"HELP WANTED", "NEEDS-HELP", "BUG"}))
	assertHelpWantedCount(t, issues, 5)
}

func makeIssue(labels []string) github.Issue {
	issue := github.Issue{}
	for idx := range labels {
		issue.Labels = append(issue.Labels, github.Label{Name: &labels[idx]})
	}
	return issue
}

func assertHelpWantedCount(t *testing.T, issues []github.Issue, count int64) {
	if actual := countHelpWantedIssues(issues); actual != count {
		t.Errorf("Should have had %d issues but had %d\n", count, actual)
	}
}
