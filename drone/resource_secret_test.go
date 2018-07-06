package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func testSecretConfigBasic(user, repo, name, value string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }
    
    resource "drone_secret" "secret" {
      repository = "${drone_repo.repo.repository}"
      name       = "%s"
      value      = "%s"
      events     = ["push", "pull_request", "tag", "deployment"]
    }
    `,
		user,
		repo,
		name,
		value,
	)
}

func TestSecret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSecretConfigBasic(
					testDroneUser,
					"repository-1",
					"password",
					"1234567890",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"repository",
						fmt.Sprintf("%s/repository-1", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"name",
						"password",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"value",
						"1234567890",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"images.#",
						"0",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"events.#",
						"4",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"events.1329302135",
						"deployment",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"events.1396138718",
						"pull_request",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"events.398155140",
						"tag",
					),
					resource.TestCheckResourceAttr(
						"drone_secret.secret",
						"events.696883710",
						"push",
					),
				),
			},
		},
	})
}

func testSecretDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_secret" {
			continue
		}

		owner, repo, err := parseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		err = client.SecretDelete(owner, repo, resource.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf(
				"Secret still exists: %s/%s:%s",
				owner,
				repo,
				resource.Primary.Attributes["name"],
			)
		}
	}

	return nil
}
