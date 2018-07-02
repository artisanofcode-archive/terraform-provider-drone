package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	testDroneUser string = os.Getenv("DRONE_USER")
	testProviders map[string]terraform.ResourceProvider
	testProvider  *schema.Provider
)

func init() {
	testProvider = Provider()
	testProviders = map[string]terraform.ResourceProvider{
		"drone": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DRONE_SERVER"); v == "" {
		t.Fatal("DRONE_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_TOKEN"); v == "" {
		t.Fatal("DRONE_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("DRONE_USER"); v == "" {
		t.Fatal("DRONE_USER must be set for acceptance tests")
	}
}
