package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

var validRepoHooks = []string{
	drone.EventPull,
	drone.EventPush,
	drone.EventTag,
	drone.EventDeploy,
}

func resourceRepo() *schema.Resource {
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
			"trusted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
			},
			"hooks": {
				Type:     schema.TypeSet,
				Optional: true,
				// ValidateFunc: validation.ValidateListUniqueStrings,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validRepoHooks, true),
				},
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceRepoCreate,
		Read:   resourceRepoRead,
		Update: resourceRepoUpdate,
		Delete: resourceRepoDelete,
		Exists: resourceRepoExists,
	}
}

func resourceRepoCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	_, err = client.RepoPost(owner, repo)

	if err != nil {
		return err
	}

	repository, err := client.RepoPatch(owner, repo, createRepo(data))

	if err != nil {
		return err
	}

	return readRepo(data, repository, err)
}

func resourceRepoRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Id())

	if err != nil {
		return err
	}

	repository, err := client.Repo(owner, repo)

	return readRepo(data, repository, err)
}

func resourceRepoUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	repository, err := client.RepoPatch(owner, repo, createRepo(data))

	return readRepo(data, repository, err)
}

func resourceRepoDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Id())

	if err != nil {
		return err
	}

	return client.RepoDel(owner, repo)
}

func resourceRepoExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Id())

	if err != nil {
		return false, err
	}

	repository, err := client.Repo(owner, repo)

	exists := (repository.Owner == owner) && (repository.Name == repo) && (err == nil)

	return exists, err
}

func createRepo(data *schema.ResourceData) (repository *drone.RepoPatch) {
	hooks := data.Get("hooks").(*schema.Set)

	trusted := data.Get("trusted").(bool)
	gated := data.Get("gated").(bool)
	timeout := int64(data.Get("timeout").(int))
	visibility := data.Get("visibility").(string)
	pull := hooks.Contains(drone.EventPull)
	push := hooks.Contains(drone.EventPush)
	deploy := hooks.Contains(drone.EventDeploy)
	tag := hooks.Contains(drone.EventTag)

	repository = &drone.RepoPatch{
		IsTrusted:   &trusted,
		IsGated:     &gated,
		Timeout:     &timeout,
		Visibility:  &visibility,
		AllowPull:   &pull,
		AllowPush:   &push,
		AllowDeploy: &deploy,
		AllowTag:    &tag,
	}

	return
}

func readRepo(data *schema.ResourceData, repository *drone.Repo, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s", repository.Owner, repository.Name))

	hooks := make([]string, 0)

	if repository.AllowPull == true {
		hooks = append(hooks, drone.EventPull)
	}

	if repository.AllowPush == true {
		hooks = append(hooks, drone.EventPush)
	}

	if repository.AllowDeploy == true {
		hooks = append(hooks, drone.EventDeploy)
	}

	if repository.AllowTag == true {
		hooks = append(hooks, drone.EventTag)
	}

	data.Set("repository", fmt.Sprintf("%s/%s", repository.Owner, repository.Name))
	data.Set("trusted", repository.IsTrusted)
	data.Set("gated", repository.IsGated)
	data.Set("timeout", repository.Timeout)
	data.Set("visibility", repository.Visibility)
	data.Set("hooks", hooks)

	return nil
}
