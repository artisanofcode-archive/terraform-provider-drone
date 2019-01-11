package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"regexp"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			keyRepository: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			keyName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			keyValue: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			keyAllowPR: {
				Type:     schema.TypeBool,
				Computed: true,
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

	namespace, repo, err := parseRepo(data.Get("repository").(string))
	if err != nil {
		return err
	}
	secret, err := client.SecretCreate(namespace, repo, createSecret(data))
	if err != nil {
		return err
	}

	return readSecret(data, namespace, repo, secret)
}

func resourceSecretRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}
	name := data.Get(keyName).(string)

	secret, err := client.Secret(namespace, repo, name)
	if err != nil {
		return err
	}

	return readSecret(data, namespace, repo, secret)
}

func resourceSecretUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}

	secret, err := client.SecretUpdate(namespace, repo, createSecret(data))
	if err != nil {
		return err
	}
	return readSecret(data, namespace, repo, secret)
}

func resourceSecretDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}
	log.Print("[ERROR] foobar")
	name := data.Get(keyName).(string)

	return client.SecretDelete(namespace, repo, name)
}

func resourceSecretExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return false, err
	}
	name := data.Get(keyName).(string)

	secret, err := client.Secret(namespace, repo, name)
	exists := (secret.Name == name) && (err == nil)

	return exists, err
}

func createSecret(data *schema.ResourceData) *drone.Secret {
	return &drone.Secret{
		Name:        data.Get(keyName).(string),
		Data:        data.Get(keyValue).(string),
		PullRequest: data.Get(keyAllowPR).(bool),
	}
}

func readSecret(data *schema.ResourceData, namespace, repo string, secret *drone.Secret) error {
	slug := fmt.Sprintf("%s/%s", namespace, repo)
	data.SetId(fmt.Sprintf("secret/%s/%s", slug, secret.Name))
	err := setResourceData(nil, data, keyRepository, slug)
	err = setResourceData(err, data, keyName, secret.Name)
	err = setResourceData(err, data, keyValue, secret.Data)
	return setResourceData(err, data, keyAllowPR, secret.PullRequest)
}
