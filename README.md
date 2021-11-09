# backend

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [go v1.14](https://golang.org/doc/go1.14)

### Development
- Install dependencies
```
./bin/setup_env.sh
```
- Run server

```
make run
```

Backend server will be served at [localhost:8080](http://localhost:8080) by default

### Deploy to Heroku with Terraform

#### Prerequisites

- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) latest version
- [Terraform](https://www.terraform.io/downloads.html)

To deploy the application to Heroku with Terraform, we need to create the Heroku API Key first:

```bash
$ heroku login
$ heroku authorizations:create --description <api key description>
```

And then, move to the `deploy/heroku` folder and run the following steps:

_Step 1:_ Copy the variable file and update the variables

```sh
$ cp terraform.tfvars.sample terraform.tfvars≠
```

*You can get the `tfvars` files from 1Password*

_Step 2:_ Initialize Terraform

```sh
$ terraform init
```

_Step 3:_ Generate an execution plan

```sh
$ terraform plan -var-file="terraform.tfvars"
```

_Step 5:_ Execute the generated plan

```sh
$ terraform apply -var-file="terraform.tfvars"
```

_Step 6:_ Build the application and push to heroku

### Wiki

##### Google OAuth2 Flow
[https://developers.google.com/identity/sign-in/web/server-side-flow](https://developers.google.com/identity/sign-in/web/server-side-flow)
