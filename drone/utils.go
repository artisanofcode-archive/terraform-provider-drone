package drone

import (
	"fmt"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

func findRepo(data *schema.ResourceData, repos []*drone.Repo) (*drone.Repo, error) {
	slug := data.Get("repository").(string)
	for _, repo := range repos {
		if slug == repo.Slug {
			return repo, nil
		}
	}
	return nil, fmt.Errorf("repository %q not found", slug)
}

func parseRepo(str string) (user, repo string, err error) {
	parts := strings.Split(str, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("Error: Invalid repository (e.g. octocat/hello-world).")
		return
	}

	user = parts[0]
	repo = parts[1]
	return
}

func parseId(str, example string) (user, repo, id string, err error) {
	parts := strings.Split(str, "/")

	if len(parts) < 3 {
		err = fmt.Errorf(
			"Error: Invalid identity (e.g. octocat/hello-world/%s).",
			example,
		)
		return
	}

	user = parts[0]
	repo = parts[1]

	id = strings.Join(parts[2:], "/")

	return
}
