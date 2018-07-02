package main

import (
	"testing"
)

func TestParseRepo(t *testing.T) {
	for _, test := range []struct {
		name, str, user, repo string
		is_error              bool
	}{
		{"Test valid repository", "octocat/hello-world", "octocat", "hello-world", false},
		{"Test another valid repository", "drone/drone", "drone", "drone", false},
		{"Test invalid repository without slash", "foobar", "", "", true},
		{"Test invalid repository with too many slashes", "foo/bar/baz", "", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			user, repo, err := parseRepo(test.str)

			if (test.is_error == true) && (err == nil) {
				t.Errorf("expected error")
			}

			if (test.is_error == false) && (err != nil) {
				t.Errorf("unexpected error")
			}

			if test.user != user {
				t.Errorf("unexpected user")
			}

			if test.repo != repo {
				t.Errorf("unexpected repo")
			}
		})
	}
}
