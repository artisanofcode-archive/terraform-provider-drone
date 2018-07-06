package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func testRepoConfigBasic(user, repo string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
      hooks      = ["push", "pull_request", "tag", "deployment"]
    }
    `, user, repo)
}

func TestRepo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRepoConfigBasic(testDroneUser, "repository-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"repository",
						fmt.Sprintf("%s/repository-1", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"hooks.#",
						"4",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"hooks.1329302135",
						"deployment",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"hooks.1396138718",
						"pull_request",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"hooks.398155140",
						"tag",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"hooks.696883710",
						"push",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"timeout",
						"0",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"timeout",
						"0",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"gated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"drone_repo.repo",
						"trusted",
						"false",
					),
				),
			},
		},
	})
}

func testRepoDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_repo" {
			continue
		}

		owner, repo, err := parseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		repositories, err := client.RepoList()

		for _, repository := range repositories {
			if (repository.Owner == owner) && (repository.Name == repo) {
				client.RepoDel(owner, repo)
				return fmt.Errorf("Repo still exists: %s/%s", owner, repo)
			}
		}
	}

	return nil
}
