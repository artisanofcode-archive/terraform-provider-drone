package main

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func testRegistryConfigBasic(user, repo, address, username, password string) string {
	return fmt.Sprintf(`
    resource "drone_repo" "repo" {
      repository = "%s/%s"
    }
    
    resource "drone_registry" "registry" {
      repository = "${drone_repo.repo.repository}"
			address    = "%s"
			username   = "%s"
			password   = "%s"
    }
    `,
		user,
		repo,
		address,
		username,
		password,
	)
}

func TestRegistry(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRegistryConfigBasic(
					testDroneUser,
					"repository-1",
					"example.com",
					"user",
					"pass",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_registry.registry",
						"repository",
						fmt.Sprintf("%s/repository-1", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_registry.registry",
						"address",
						"example.com",
					),
					resource.TestCheckResourceAttr(
						"drone_registry.registry",
						"username",
						"user",
					),
					resource.TestCheckResourceAttr(
						"drone_registry.registry",
						"password",
						"pass",
					),
				),
			},
		},
	})
}

func testRegistryDestroy(state *terraform.State) error {
	client := testProvider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_registry" {
			continue
		}

		owner, repo, err := parseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		err = client.RegistryDelete(owner, repo, resource.Primary.Attributes["address"])

		if err == nil {
			return fmt.Errorf(
				"Registry still exists: %s/%s:%s",
				owner,
				repo,
				resource.Primary.Attributes["address"],
			)
		}
	}

	return nil
}
