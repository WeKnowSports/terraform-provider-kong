# terraform-provider-kong (rapid7 fork)
**NOTE**: This repository was forked from https://github.com/WeKnowSports/terraform-provider-kong.

This fork's master branch updates main.go to point at this fork's kong package for use with `go get`/`go install`.

The master branch includes several patches contributed (or in review) from the rapid7 fork back to the upstream project.

## Installation
Install a binary from [releases](https://github.com/rapid7/terraform-provider-kong/releases) into your `terraform.d/plugins/${GOOS}_${GOARCH}` directory.

**NOTE**: Optionally this can be installed under `~/.terraform.d/plugins` in your home directory if you have multiple projects using the plugin.

```
GOOS=darwin
GOARCH=amd64
VERSION=0.1.0

wget -O terraform-provider-kong_v${VERSION} https://github.com/rapid7/terraform-provider-kong/releases/download/${VERSION}/${GOOS}_${GOARCH}-terraform-provider-kong_v${VERSION}
chmod +x terraform-provider-kong_v${VERSION}

# mv terraform-provider-kong_v${VERSION} ~/.terraform.d/plugins/${GOOS}_${GOARCH}/
mv terraform-provider-kong_v${VERSION} terraform.d/plugins/${GOOS}_${GOARCH}/
```

## Usage
Once [installed](#installation) you can use the following to confirm you can use the plugin:

```
terraform init
terraform plan
```

## Development
See the [provider plugin codebase](https://www.terraform.io/docs/plugins/provider.html) information at terraform.io.

1. Ensure `GOPATH`, `GOBIN`, and `PATH` are set.

  ```bash
  export GOPATH="${GOPATH:-$(echo -n ~/go)}"
  export GOBIN="${GOPATH}/bin"
  export PATH="${GOPATH}/bin:${PATH}"
  ```

2. Clone this project into `${GOPATH}/src/github.com/rapid7/terraform-provider-kong`.
3. Install `terraform` into your `GOPATH` (required since `terraform init` will search for providers inside of `GOPATH` only):

  ```
  go get github.com/hashicorp/terraform
  go install github.com/hashicorp/terraform
  ```

4. Make your intended changes.
5. Use `go install` inside of `${GOPATH}/src/github.com/rapid7/terraform-provider-kong`.
6. Use `terraform init` in your terraform project.
7. Confirm your changes work. Profit! :tada:

## Distribution
1. Figure out what the last released version was by looking at https://github.com/rapid7/terraform-provider-kong/releases.
2. Create binaries for different platforms by bumping versions (use http://semver.org):

  ```
  # Ensure gox is installed and in your PATH.
  go get github.com/mitchellh/gox

  VERSION="0.1.0"
  gox -osarch="darwin/amd64 linux/amd64" -output="dist/{{.OS}}_{{.Arch}}/{{.Dir}}_v${VERSION}"  github.com/rapid7/terraform-provider-kong
  ```

3. Create an annotated `git tag` for your version (e.g. `v0.1.0`).
4. Push the tag and create a github release.
    * Attach your binaries.
    * Add the checksums into the release notes (using `shasum -a 256`).

---

# Terraform provider for Kong
Uses [Terraform](http://www.terraform.io) to configure APIs in [Kong](http://www.getkong.org). It fully supports creating APIs and consumers, but plugins and credentials are not complete (most plugins will work though).

```
go build -o tests/terraform-provider-kong
```

## Compile and terraform plan / apply

### Start kong


```Shell
docker-compose up -d
```

## Run plan
```Shell
./start plan
```

## Run apply
```Shell
./start apply
```

## Example usage

Please refer to terraform/tests
