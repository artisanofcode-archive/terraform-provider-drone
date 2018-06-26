package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

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

func toStringList(set *schema.Set) (result []string) {
	result = make([]string, set.Len())

	for i, v := range set.List() {
		result[i] = v.(string)
	}

	return
}
