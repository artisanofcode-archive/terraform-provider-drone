package drone

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL for the drone server",
				DefaultFunc: schema.EnvDefaultFunc("DRONE_SERVER", nil),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Token for the drone server",
				DefaultFunc: schema.EnvDefaultFunc("DRONE_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"drone_registry": resourceRegistry(),
			"drone_repo":     resourceRepo(),
			"drone_secret":   resourceSecret(),
			"drone_user":     resourceUser(),
		},
		ConfigureFunc: providerConfigureFunc,
	}
}

func providerConfigureFunc(data *schema.ResourceData) (interface{}, error) {
	config := new(oauth2.Config)

	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{AccessToken: data.Get("token").(string)},
	)

	client := drone.NewClient(data.Get("server").(string), auther)

	if _, err := client.Self(); err != nil {
		return nil, fmt.Errorf("drone client failed: %s", err)
	}

	return client, nil
}
