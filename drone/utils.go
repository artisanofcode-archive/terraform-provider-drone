package drone

import (
	"fmt"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	keyConfig     = "configuration"
	keyProtected  = "protected"
	keyRepository = "repository"
	keyTimeout    = "timeout"
	keyTrusted    = "trusted"
	keyVisibility = "visibility"
	keyName       = "name"
	keyValue      = "value"
	keyAllowPR    = "allow_pull_requests"
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

func parseID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("Error: Invalid identity (e.g. octocat/hello-world/fancy_pants)")
	}

	return parts[0], parts[1], parts[2], nil
}
