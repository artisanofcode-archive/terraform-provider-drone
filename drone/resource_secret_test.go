package drone

import (
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/suite"

	"testing"
)

type SecretTestSuite struct {
	suite.Suite
}

func (s *SecretTestSuite) TestSecretCreate() {
	data := resourceSecret().Data(&terraform.InstanceState{})
	data.Set(keyName, "test")
	data.Set(keyValue, "password")
	data.Set(keyAllowPR, true)

	secret := createSecret(data)

	s.Equal("test", secret.Name)
	s.Equal("password", secret.Data)
	s.Equal(true, secret.PullRequest)
}

func (s *SecretTestSuite) TestReadSecret() {
	data := resourceSecret().Data(&terraform.InstanceState{})
	secret := &drone.Secret{
		Name:        "test",
		Data:        "password",
		PullRequest: false,
	}

	if err := readSecret(data, "test", "repo", secret); err != nil {
		s.Fail("readSecret returned an error", err)
	}

	s.Equal("secret/test/repo/test", data.Id())
	s.Equal("test", data.Get(keyName).(string))
	s.Equal("password", data.Get(keyValue).(string))
	s.Equal(false, data.Get(keyAllowPR).(bool))
}

func TestSecretSuite(t *testing.T) {
	suite.Run(t, &SecretTestSuite{})
}
