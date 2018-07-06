package drone

import (
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
)

var defaultUserEvents = []string{
	drone.EventPush,
	drone.EventTag,
	drone.EventDeploy,
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
	}
}

func resourceUserCreate(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.UserPost(createUser(data))

	return readUser(data, user, err)
}

func resourceUserRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	user, err := client.User(data.Id())

	return readUser(data, user, err)
}

func resourceUserDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(drone.Client)

	return client.UserDel(data.Id())
}

func resourceUserExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(drone.Client)

	login := data.Id()

	user, err := client.User(login)

	exists := (user.Login == login) && (err == nil)

	return exists, err
}

func createUser(data *schema.ResourceData) (user *drone.User) {
	user = &drone.User{
		Login: data.Get("login").(string),
	}

	return
}

func readUser(data *schema.ResourceData, user *drone.User, err error) error {
	if err != nil {
		return err
	}

	data.SetId(user.Login)

	data.Set("login", user.Login)

	return nil
}
