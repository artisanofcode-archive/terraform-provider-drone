package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"regexp"
)

var (
	defaultSecretEvents = []string{
		drone.EventPush,
		drone.EventTag,
		drone.EventDeploy,
	}
	validSecretEvents = []string{
		drone.EventPull,
		drone.EventPush,
		drone.EventTag,
		drone.EventDeploy,
	}
)

func resourceSecret() *schema.Resource {
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"images": {
				Type:     schema.TypeSet,
				Optional: true,
				// ValidateFunc: validation.ValidateListUniqueStrings,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"events": {
				Type:     schema.TypeSet,
				Optional: true,
				// ValidateFunc: validation.ValidateListUniqueStrings,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validSecretEvents, true),
				},
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Exists: resourceSecretExists,
	}
}

func resourceSecretCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	secret, err := client.SecretCreate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, name, err := parseId(data.Id(), "secret_password")

	if err != nil {
		return err
	}

	secret, err := client.Secret(owner, repo, name)

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, err := parseRepo(data.Get("repository").(string))

	if err != nil {
		return err
	}

	secret, err := client.SecretUpdate(owner, repo, createSecret(data))

	data.Set("value", data.Get("value").(string))

	return readSecret(data, owner, repo, secret, err)
}

func resourceSecretDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	owner, repo, name, err := parseId(data.Id(), "secret_password")

	if err != nil {
		return err
	}

	return client.SecretDelete(owner, repo, name)
}

func resourceSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	owner, repo, name, err := parseId(data.Id(), "secret_password")

	if err != nil {
		return false, err
	}

	secret, err := client.Secret(owner, repo, name)

	exists := (secret.Name == name) && (err == nil)

	return exists, err
}

func createSecret(data *schema.ResourceData) (secret *drone.Secret) {
	events := []string{}
	eventSet := data.Get("events").(*schema.Set)
	for _, v := range eventSet.List() {
		events = append(events, v.(string))
	}

	images := []string{}
	imageSet := data.Get("images").(*schema.Set)
	for _, v := range imageSet.List() {
		images = append(images, v.(string))
	}

	secret = &drone.Secret{
		Name:   data.Get("name").(string),
		Value:  data.Get("value").(string),
		Images: images,
		Events: events,
	}

	if len(secret.Events) == 0 {
		secret.Events = defaultSecretEvents
	}

	return
}

func readSecret(data *schema.ResourceData, owner, repo string, secret *drone.Secret, err error) error {
	if err != nil {
		return err
	}

	data.SetId(fmt.Sprintf("%s/%s/%s", owner, repo, secret.Name))

	data.Set("repository", fmt.Sprintf("%s/%s", owner, repo))
	data.Set("name", secret.Name)
	data.Set("images", secret.Images)
	data.Set("events", secret.Events)

	return nil
}
