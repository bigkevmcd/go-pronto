package main

import (
	"io/ioutil"
	"path"
	"testing"

	gc "gopkg.in/v1/check"
)

func Test(t *testing.T) { gc.TestingT(t) }

type ConfigSuite struct {
	path string
}

var _ = gc.Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpTest(c *gc.C) {
	s.path = path.Join(c.MkDir(), "config.yml")
}

// Calling ConfigFromYaml with an invalid path generates an error.
func (s *ConfigSuite) TestConfigFromYamlReturnsErrorOnMissingFile(c *gc.C) {
	_, err := ConfigFromYaml(s.path)
	c.Assert(err, gc.ErrorMatches, "opening config file: open.*no such file or directory")
}

// The Config should extract Swift credentials
func (s *ConfigSuite) TestConfigHasCredentials(c *gc.C) {
	data := []byte(`
credentials:
  auth-url: https://keystone.example.com/v2.0/
  tenant-name: example_tenant
  region: example
  username: testing
  password: letmein
`)
	ioutil.WriteFile(s.path, data, 0600)
	config, _ := ConfigFromYaml(s.path)
	c.Assert(config.Credentials.AuthUrl, gc.Equals, "https://keystone.example.com/v2.0/")
	c.Assert(config.Credentials.TenantName, gc.Equals, "example_tenant")
	c.Assert(config.Credentials.Region, gc.Equals, "example")
	c.Assert(config.Credentials.Username, gc.Equals, "testing")
	c.Assert(config.Credentials.Password, gc.Equals, "letmein")
}

// The Config should extract port and container details.
func (s *ConfigSuite) TestConfigHasOtherDetails(c *gc.C) {
	data := []byte(`
container: 1233f97438f0456788273d757ab49101
port: :9080
`)
	ioutil.WriteFile(s.path, data, 0600)
	config, err := ConfigFromYaml(s.path)
	c.Assert(err, gc.IsNil)
	c.Assert(config.Port, gc.Equals, ":9080")
	c.Assert(config.Container, gc.Equals, "1233f97438f0456788273d757ab49101")
}
