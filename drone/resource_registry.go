package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,
		Exists: resourceRegistryExists,
	}
}

func resourceRegistryCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	registry, err := client.RegistryCreate(owner, repo, createRegistry(data))

	data.Set("password", data.Get("password").(string))

	return readRegistry(data, owner, repo, registry, err)
}

func resourceRegistryRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, address, err := parseId(data.Id(), "drone.io")

	if err != nil {
		return err
	}

	registry, err := client.Registry(owner, repo, address)

	return readRegistry(data, owner, repo, registry, err)
}

func resourceRegistryUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	registry, err := client.RegistryUpdate(owner, repo, createRegistry(data))

	data.Set("password", data.Get("password").(string))

	return readRegistry(data, owner, repo, registry, err)
}

func resourceRegistryDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, address, err := parseId(data.Id(), "drone.io")

	if err != nil {
		return err
	}

	return client.RegistryDelete(owner, repo, address)
}

func resourceRegistryExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, address, err := parseId(data.Id(), "drone.io")

	if err != nil {
		return false, err
	}

	registry, err := client.Registry(owner, repo, address)

	exists := (registry.Address == address) && (err == nil)

	return exists, err
}

func createRegistry(data *schema.ResourceData) (registry *drone.Registry) {
	registry = &drone.Registry{
		Address:  data.Get("address").(string),
		Username: data.Get("username").(string),
		Password: data.Get("password").(string),
	}

	return
}

func readRegistry(data *schema.ResourceData, owner, repo string, registry *drone.Registry, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s/%s", owner, repo, registry.Address))

	data.Set("repository", fmt.Sprintf("%s/%s", owner, repo))
	data.Set("address", registry.Address)
	data.Set("username", registry.Username)

	return nil
}
