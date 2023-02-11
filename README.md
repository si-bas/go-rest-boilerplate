# GO REST BOILERPLATE

Service to serve REST API.

## Prerequisite 

### Setup Configuration
There are 2 ways to setup configuration:
1. Local: will read complete onfig from .config.json locally
    - Copy `.config` file to `.config.json` on the root folder
    - Fill config value directly there
2. (Work in progress) Consul: will read complete config from consul
    - Copy `.example.env` file to `.env` on the root folder
    - Fill consul configuration
    - ```
      set -o allexport
      source .env
      set +o allexport```

### How do I run locally? ###

* Setup Configuration in above section
* Run API: `go run main.go serve`

### Database Migrations ###

* Folder: databases/mysql
* Library: https://github.com/pressly/goose
* Getting started:
    * `brew install goose`
    * create new mysql database, e.g. appdb
    * cd `databases/mysql`
    * `export GOOSE_DRIVER=mysql GOOSE_DBSTRING="root:password@/appdb?parseTime=true"`
    * `goose status`
    * create new migration file `goose create file_name_in_snake_case sql`
    * `goose up`
* Configurations: `.config.json`

### Contribution guidelines ###

* Writing tests
* Code review
* Other guidelines

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact