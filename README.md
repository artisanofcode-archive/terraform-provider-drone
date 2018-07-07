# Drone Terraform Provider

A [Terraform](https://www.terraform.io) provider for configuring the 
[Drone](https://drone.io) continuous delivery platform.

## Installing

You can download the plugin from the [Releases](https://github.com/artisanofcode/terraform-provider-drone/releases/latest) page,
for help installing please refer to the [Official Documentation](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).


## Example

```terraform
provider "drone" {
  server = "https:://ci.example.com/"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoib2N0b2NhdCIsInR5cGUiOiJ1c2VyIn0.Fg0eYxO9x2CfGIvIHDZKhQbCGbRAsSB_iRDJlDEW6vc"
}

resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visability = "public"
  hooks      = ["push", "pull_request", "tag", "deployment"]
}

resource "drone_secret" "master_password" {
  repository = "${resource.hello_world.repository}"
  name       = "master_password"
  value      = "correct horse battery staple"
  events     = ["push", "pull_request", "tag", "deployment"]
}
```

## Provider

#### Example Usage

```terraform
provider "drone" {
  server = "https://ci.example.com/"
  token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0Ijoib2N0b2NhdCIsInR5cGUiOiJ1c2VyIn0.Fg0eYxO9x2CfGIvIHDZKhQbCGbRAsSB_iRDJlDEW6vc"
}
````

#### Argument Reference

* `server` - (Optional) The Drone servers url, It must be provided, but can also
  be sourced from the `DRONE_SERVER` environment variable.
* `token` - (Optional) The Drone servers api token, It must be provided, but can
  also be sourced from the `DRONE_TOKEN` environment variable.

## Resources

### `drone_registry`

Manage a repository registry.

#### Example Usage

```terraform
resource "drone_registry" "docker_io" {
  repository = "octocat/hello-world"
  address    = "docker.io"
  username   = "octocat"
  password   = "correct horse battery staple"
}
```

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `address` - (Required) Registry address.
* `username` - (Required) Registry username.
* `password` - (Required) Registry password.

### `drone_repo`

Activate and configure a repository.

#### Example Usage

```terraform
resource "drone_repo" "hello_world" {
  repository = "octocat/hello-world"
  visability = "public"
  hooks      = ["push", "pull_request", "tag", "deployment"]
}
```

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `trusted` - (Optional) Repository is trusted (default: `false`).
* `gated` - (Optional) Repository is gated (default: `false`).
* `timeout` - (Optional) Repository timeout (default: `0`).
* `visibility` - (Optional) Repository visibility (default: `private`).
* `hooks` - (Optional) List of hooks this repository should setup is limited to, 
  values must be `push`, `pull_request`, `tag`, and/or `deployment`.

### `drone_secret`

Manage a repository secret.

#### Example Usage

```terraform
resource "drone_secret" "master_password" {
  repository = "octocat/hello-world"
  name       = "master_password"
  value      = "correct horse battery staple"
  events     = ["push", "pull_request", "tag", "deployment"]
}
````

#### Argument Reference

* `repository` - (Required) Repository name (e.g. `octocat/hello-world`).
* `name` - (Required) Secret name.
* `value` - (Required) Secret value.
* `images` - (Optional) List of images this secret is limited to.
* `events` - (Optional) List of events this repository should setup is limited to, 
  values must be `push`, `pull_request`, `tag`, and/or `deployment` (default: `["push", "tag", "deployment"]`).

### `drone_user`

Manage a user.

#### Example Usage

```terraform
resource "drone_user" "octocat" {
  login = "octocat"
}
````

#### Argument Reference

* `login` - (Required) Login name.

## Source

To install from source:

```shell
git clone git://github.com/artisanofcode/terraform-provider-drone.git
cd terraform-provider-drone
go get
go build
```

## Licence

This project is licensed under the [MIT licence](http://dan.mit-license.org/).

## Meta

This project uses [Semantic Versioning](http://semver.org/).
