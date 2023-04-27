# Project Title

One or two sentences describing the purpose of your project.

## Table of Contents

- [Project Title](#project-title)
  - [Table of Contents](#table-of-contents)
  - [About the Project](#about-the-project)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [License](#license)
  - [Contact](#contact)

## About the Project

The purpose of the project is to create a robust APIs and web client application for user task management that allows users to manage their tasks with comprehensive authorization.

For full API document [please checkout this section](https://github.com/marktrs/simple-todo/tree/main/apidoc)

Web application

## Getting Started

To start using this application on local hosted server you need to clone this repository along with [the submodule repository](https://github.com/marktrs/simple-todo-client)

```
git clone --recurse-submodules https://github.com/marktrs/simple-todo.git
```

This project is using Git submodule to allows managing dependencies in a larger Git repository, referencing other repositories as subdirectories while maintaining separate version control and revision history for each.

### Prerequisites

To build using docker:

- Docker [Docker](https://www.docker.com/)

To build from source without docker:

- Installed [Golang 1.20](https://go.dev/dl/) for API server testing, building from source.
- Installed [Node.JS 19](https://nodejs.org/dist/v19.8.0/) for sveltekit client application building on local and run unit testing
- Installed Postgres with a user configuration from the [environment variable file](https://github.com/marktrs/simple-todo/blob/main/.env.example)

## Usage

Using docker

1. Start running services using docker compose

```sh
$ docker compose up -d --build
```

You should see the following output:

```shell
[+] Running 4/4
 ⠿ Network simple-todo_default   Created                                                                                                                                      0.0s
 ⠿ Container postgres            Started                                                                                                                                      0.4s
 ⠿ Container simple-todo-api     Started                                                                                                                                      0.6s
 ⠿ Container simple-todo-client  Started
```

2. Inspect server logs to ensure the server is running

```sh
$ docker logs -f simple-todo-api
```

You should see the following output:

```shell
{"level":"info","time":"2023-04-27T09:09:37Z","message":"seeding user data..."}
{"level":"info","time":"2023-04-27T09:09:45Z","message":"server is listening on port 3000..."}
```

3. Checkout /metrics on [http://localhost:3000/metrics](http://localhost:3000/metrics) for server metrics and /health for server health check

```sh
curl --location 'http://localhost:3000/api/health'
```

4. Then navigate to web application on [http://localhost:8080/](http://localhost:8080/)
5. login with default user or create a new one on registration page [http://localhost:8080/register](http://localhost:8080/register)

```sh
username: test01
password: 1111

username: test02
password: 2222

username: test03
password: 3333

username: test04
password: 4444
```

Stop running services

```sh
$ docker compose down
```

### Run Test

Run unit testing from source

API server

```sh
make test-unit
```

or get test coverage profile

```sh
make test-coverage
make test-report
```

Sveltekit - client application

```sh
cd client/
npm install -g pnpm
pnpm install
pnpm test
```

## Architecture Design

The API server is designed using a layered architecture pattern. The layered architecture pattern separates an application into logical layers that interact with each other to perform a specific task and easily to write a unit test of a specific function without dependencies implementation concerns.

To ensure the application's scalability and maintainability, I implemented with Dependency injection, Separation of concerns and Single Responsibility Principle (SRP) design principles.

I use the OpenAPI specification to document the API server's endpoints and responses to provides a clear and standardized way of describing the API, Enables the automatic generation of documentation and client libraries, which can save significant development time and effort and provides a machine-readable format that can be used for automated testing.

Tools usage in this project:

- Database : Postgres
  - ORM : GORM
- HTTP Router : Go/Fiber
- Client web application : SvelteKit

### Project Layout

```
.
├── Dockerfile
├── Makefile
├── README.md
├── apidoc
│   ├── README.md
│   ├── open-api.yml
│   └── swagger.json
├── config
│   └── config.go
├── coverage.out
├── database
│   └── ...
├── docker-compose.yml
├── go.mod
├── go.sum
├── handler
│   └── ...
├── integration_test.go
├── logger
│   └── ...
├── middleware
│   └── ...
├── model
│   └── ...
├── repository
│   └── ...
├── router
│   └── ...
├── server
│   └── ...
├── start.sh
├── temp
├── testutil
│   └── mocks
│       └── repository
│           └── ...
└── main.go
```

`apidoc`: Contains OpenAPI specification and Swagger JSON files

`client`: Contains files for the SvelteKit web client application, including dependencies and tests

`config`: Contains configuration files for the project

`database`: Contains files related to the database, including connector, migration, and seeder files

`handler`: Contains files that define the HTTP request handlers for each route of the API

`logger`: Contains files related to logging and log management

`middleware`: Contains middleware functions for the API, including authentication, error handling, and logging

`model`: Contains data transfer object and schema that map to the application's database tables

`repository`: Contains implementations of data access layer for the CRUD operations for each entity in the application

`router`: Contains the router configuration for the API

`server`: Contains the server initialization and starting code

`testutil`: Contains mock implementations for repositories used in testing

`temp`: Contains temporary log files

`root directory`: Contains miscellaneous files, including Makefile, Dockerfile, integration tests, and shell scripts

## Future works

- `Testing` More coverage test % beside the core functionalities
- `Caching` Caching strategy
- `gRPC` to provide more efficient binary protocol, which can result in faster and more efficient communication between client and server
- `document` An endpoint to serve API document as HTML

## Notes

- Makefile **MUST NOT** change well-defined command semantics, see Makefile for details.
