package drone

import (
	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/suite"

	"testing"
)

type RepoTestSuite struct {
	suite.Suite
}

func (s *RepoTestSuite) TestRepoRead() {
	data := resourceRepo().Data(&terraform.InstanceState{})
	repo := &drone.Repo{
		Slug:       "test/test",
		Config:     "foo.yml",
		Protected:  true,
		Visibility: "private",
	}
	err := readRepo(data, repo)
	if err != nil {
		s.Fail("readRepo returned an error", err)
	}

	s.Equal("repo/test/test", data.Id())
	s.Equal("foo.yml", data.Get(keyConfig))
	s.Equal(true, data.Get(keyProtected))
	s.Equal(false, data.Get(keyTrusted))
	s.Equal("private", data.Get(keyVisibility))
}

func (s *RepoTestSuite) TestUpdateRepo() {
	data := resourceRepo().Data(&terraform.InstanceState{})
	data.Set(keyConfig, "foo.yml")
	data.Set(keyProtected, true)
	data.Set(keyVisibility, "private")

	repo := &drone.Repo{}
	patch := updateRepo(data, repo)

	s.Equal("foo.yml", repo.Config)
	s.Equal("foo.yml", *patch.Config)
	s.Equal(true, repo.Protected)
	s.Equal(true, *patch.Protected)
	s.Equal(false, repo.Trusted)
	s.Equal(false, *patch.Trusted)
	s.Equal("private", repo.Visibility)
	s.Equal("private", *patch.Visibility)
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, &RepoTestSuite{})
}
