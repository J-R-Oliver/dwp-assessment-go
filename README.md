# DWP Assessment Go

[![Build](https://github.com/J-R-Oliver/dwp-assessment-go/actions/workflows/build.yml/badge.svg)](https://github.com/J-R-Oliver/dwp-assessment-go/actions/workflows/build.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/J-R-Oliver/dwp-assessment-go)](https://github.com/gomods/athens)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/J-R-Oliver/dwp-assessment-go)](https://github.com/J-R-Oliver/dwp-assessment-go/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/J-R-Oliver/dwp-assessment-go.svg)](https://pkg.go.dev/github.com/J-R-Oliver/dwp-assessment-go)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)

<table>
<tr>
<td>
An API which calls the API at https://dwp-techtest.herokuapp.com, and returns people who are listed as either living 
in London, or whose current coordinates are within 50 miles of London.
</td>
</tr>
</table>

## Contents

- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [OpenAPI Specification](#openapi-specification)
- [Configuration](#configuration)
- [Testing](#testing)
- [Conventional Commits](#conventional-commits)
- [GitHub Actions](#github-actions)

## Getting Started

### Prerequisites

To install, run and modify this project you will need to have:

- [Go](https://go.dev)
- [Git](https://git-scm.com)
- [Docker](https://www.docker.com)

### Installation

To start, please `fork` and `clone` the repository to your local machine. You are able to run the service directly on
the command line, using _Docker_, or with your IDE of choice.

### Running

#### Command Line

To run the service from the command line first execute the following command:

```shell
go run cmd/web/*
```

#### Docker

The included `Dockerfile` allows the service to be containerised and run using _Docker_. You can build an image and then
start a container from it, exposing the internal port of `8080`.

```shell
docker build --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') -t j-r-oliver/dwp-assessment-go .
docker run --name dwp-assessment-go -p 8080:8080 j-r-oliver/dwp-assessment-go
```

#### Docker Compose

A _Docker Compose_ file has been provided to facilitate running the service in conjunction with
a [WireMock](http://wiremock.org) of the upstream DWP API. This _WireMock_ provides a local instance of the DWP API
allowing for development regardless of external connectivity.

To start the service and the _WireMock_ you will need to execute _Docker Compose_ `up dwp-assessment-go`. This will
start both services simultaneously and set the `PEOPLE_ENDPOINT` environment variable to the local _WireMock_ instance.

```shell
docker-compose up dwp-assessment-go
```

By passing the name of the `service` defined in the `docker-compose.yml` file you may start the services individually.
By default the `PEOPLE_ENDPOINT` environment variable will be set to the local _WireMock_ instance when starting from _
Docker Compose_. This can be overridden to the hosted _Heroku_ API by passing in the `PEOPLE_ENDPOINT` environment
variable as demonstrated below.

```shell
docker-compose up wiremock
```

```shell
PEOPLE_ENDPOINT=https://dwp-techtest.herokuapp.com docker-compose up dwp-assessment-go
```

## API Endpoints

There are two RESTful API endpoints available:

> `/api/people`

Returns all people available from the DWP API.

> `/api/people/{city}`

Returns all people who are listed as either living in the `city` or whose current coordinates are within 50 miles of the
coordinates of the `city`. Currently, only `london` has been configured as an available `city`. All other requests will
result in a `404 - City Not Found` response. For example `/api/people/london` will return all people living in London or
whose coordinates are within 50 miles of London.

An optional query has also been configured for the `{city}` endpoint to amend the default distance. For example, the
path and query below would return all people living in London or whose coordinates are within 25 miles of London

> `/api/people/london?distance=25`

## OpenAPI Specification

An [OpenAPI Specification](https://spec.openapis.org/oas/v3.1.0) has been provided and can be found
in [./openapi-specification](./openapi-specification/openapi-specification.yml). The specification hasn't been used for
code generation due the desire to explore Go's capabilities.

## Configuration

Service configuration is managed using environment variables and the [configuration.yaml](./configuration.yaml). This
allows the service to be aware of necessary configuration and to provide sensible defaults. City coordinates can be
configured
by adding to the `cities` key.

```yaml
cities:
  London:
    lat: 51.514248
    lon: -0.093145
  manchester:
    lat: 53.480759
    lon: -2.242630
```

### Environment Variables

The following environment variables are available for configuration:

| Environment Variable | Default                            | Description                                                               |
|----------------------|------------------------------------|---------------------------------------------------------------------------|
| PORT                 | 8080                               | Port number for the service                                               |
| LOGGING_LEVEL        | info                               | Sets the logging level to be outputted to the logs (error, info or debug) |
| PEOPLE_ENDPOINT      | https://dwp-techtest.herokuapp.com | People / Users API endpoint                                               |
| $PEOPLE_DISTANCE     | 50                                 | Default distance in miles from city's coordinates                         |

## Testing

All tests have been written using the [testing](https://pkg.go.dev/testing) package from the
[Standard library](https://pkg.go.dev/std).

### Unit Tests

To run the unit tests execute:

```shell
go test -v ./...
```

Code coverage is also measured by using the `testing` package. To run tests with coverage execute:

```shell
go test -coverprofile=coverage.out  ./...
```

### Component Tests

Component tests has been written that exercise the service against the `WireMock`. The easiest way to execute the tests
is to use `docker-compose`. The following command will start the `WireMock`, then the service under tests, and then
execute the tests.

```shell
docker-compose up --renew-anon-volumes --exit-code-from api-tests
```

To run the tests locally, without `Docker`, execute the following command:

```shell
go test -v -tags=component ./api-tests
```

The URL of the service inside the tests can be configured using the `BASE_URL` environment variable. A default
of `http://localhost:8080`has been configured.

## Conventional Commits

This project uses the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for commit
messages. The specification provides a simple rule set for creating commit messages, documenting features, fixes, and
breaking changes in commit messages.

A [pre-commit](https://pre-commit.com) [configuration file](.pre-commit-config.yaml) has been provided to automate
commit
linting. Ensure that *pre-commit* has been [installed](https://www.conventionalcommits.org/en/v1.0.0/) and execute...

```shell
pre-commit install
````

...to add a commit [Git hook](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) to your local machine.

An automated pipeline job has been [configured](.github/workflows/build.yml) to lint commit messages on a push.

## GitHub Actions

A CI/CD pipeline has been created using [GitHub Actions](https://github.com/features/actions) to automated tasks such as
linting and testing.

### Build Workflow

The [build](./.github/workflows/build.yml) workflow handles integration tasks. This workflow consists of three
jobs, `Git`and `Go`, that run in parallel, and then `Docker`, if `Git` and `Go` are successful. This workflow is
triggered on a push to a branch.

#### Git

This job automates tasks relating to repository linting and enforcing best practices.

#### Go

This job automates `Go` specific tasks.

#### Docker

This job automates `Docker` specific tasks and the component tests.
