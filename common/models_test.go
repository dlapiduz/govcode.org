package common

import (
	"testing"
)

func TestHelpWantedLabels(t *testing.T) {

	var labelTests = []struct {
		labels   string
		expected bool
	}{
		{"help wanted", true},
		{"help_wanted", true},
		{"needs-help", true},
		{"needs_help,help-wanted", true},
		{"42,NEEDS_HELP,HELP WANTED,BUG", true},

		{"42", false},
		{"FALSE", false},
		{"TRUE", false},
		{"tomorrow and tomorrow and tomorrow creeps in this petty place from day to day to the last syllable of recorded time,Life... is a tale told by an idiot, full of sound and fury, signifying nothing", false},
	}
	for _, test := range labelTests {
		issue := Issue{Labels: test.labels}
		if issue.HelpWanted() != test.expected {
			t.Errorf("Expected Labels %v to be %v but were not", issue.Labels, test.expected)
		}
	}
}
