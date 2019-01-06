package drone

import (
	"regexp"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	keyConfig     = "configuration"
	keyProtected  = "protected"
	keyRepository = "repository"
	keyTimeout    = "timeout"
	keyTrusted    = "trusted"
	keyVisibility = "visibility"
)

func resourceRepo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			keyConfig: {
				Type:     schema.TypeString,
				Default:  ".drone.yml",
				Optional: true,
			},
			keyProtected: {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			keyRepository: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[^/ ]+/[^/ ]+$"),
					"Invalid repository (e.g. octocat/hello-world)",
				),
			},
			keyTimeout: {
				Type:     schema.TypeInt,
				Default:  60,
				Optional: true,
			},
			keyTrusted: {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			keyVisibility: {
				Type:     schema.TypeString,
				Default:  "private",
				Optional: true,
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

	// sync and get the repo list
	repos, err := client.RepoListSync()
	if err != nil {
		return err
	}

	// search the repo, don't make a second api call to client.Repo
	repo, err := findRepo(data, repos)
	if err != nil {
		return err
	}

	if !repo.Active {
		repo, err = client.RepoEnable(repo.Namespace, repo.Name)
		if err != nil {
			return err
		}
	}
	repo, err = client.RepoUpdate(repo.Namespace, repo.Name, updateRepo(data, repo))
	if err != nil {
		return err
	}

	return readRepo(data, repo)
}

func resourceRepoRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespeace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}

	repository, err := client.Repo(namespeace, repo)
	if err != nil {
		return err
	}

	return readRepo(data, repository)
}

func resourceRepoUpdate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)
	owner, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}

	repository, err := client.RepoUpdate(owner, repo, updateRepo(data, &drone.Repo{}))
	if err != nil {
		return err
	}

	return readRepo(data, repository)
}

func resourceRepoDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return err
	}

	return client.RepoDisable(namespace, repo)
}

func resourceRepoExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	namespace, repo, err := parseRepo(data.Get(keyRepository).(string))
	if err != nil {
		return false, err
	}

	repository, err := client.Repo(namespace, repo)
	// FIXME: check if repo is active?
	exists := (repository.Namespace == namespace) && (repository.Name == repo) && (err == nil)

	return exists, err
}

func updateRepo(data *schema.ResourceData, repository *drone.Repo) *drone.RepoPatch {

	configuration, ok := data.GetOk(keyConfig)
	if ok {
		repository.Config = configuration.(string)
	}

	protected, ok := data.GetOk(keyProtected)
	if ok {
		repository.Protected = protected.(bool)
	}

	timeout, ok := data.GetOk(keyTimeout)
	if ok {
		repository.Timeout = timeout.(int64)
	}

	trusted, ok := data.GetOk(keyTrusted)
	if ok {
		repository.Trusted = trusted.(bool)
	}

	visibility, ok := data.GetOk(keyVisibility)
	if ok {
		repository.Visibility = visibility.(string)
	}

	patch := &drone.RepoPatch{
		Config:     &repository.Config,
		Protected:  &repository.Protected,
		Timeout:    &repository.Timeout,
		Visibility: &repository.Visibility,
	}
	return patch
}

func readRepo(data *schema.ResourceData, repository *drone.Repo) error {
	err := setResourceData(nil, data, keyConfig, repository.Config)
	err = setResourceData(err, data, keyProtected, repository.Protected)
	err = setResourceData(err, data, keyTimeout, repository.Timeout)
	err = setResourceData(err, data, keyTrusted, repository.Trusted)
	return setResourceData(err, data, keyVisibility, repository.Visibility)
}

func setResourceData(err error, data *schema.ResourceData, key string, value interface{}) error {
	if err != nil {
		return err
	}
	return data.Set(key, value)
}
