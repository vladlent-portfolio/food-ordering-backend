# Backend for Food Ordering App

![App Image](https://food-ordering.app/assets/img/sm_image.png) 

## Table of contents
* [Overview](#overview)
* [Technologies](#technologies)
* [Features](#features)
* [Running locally](#running-locally)
    * [Development server](#development-server)
    * [Test](#test)
    * [Build](#build)

## Overview
API for [Food Ordering App](https://github.com/vladlent-portfolio/food-ordering-frontend) written in Go.  
**Schema and interaction are available via [Swagger](https://api.food-ordering.app/swagger/index.html).**  
*Since Swagger doesn't support cookie-based authorizations you should **Sign In** [here](https://food-ordering.app) 
(it already has ready-to-use users, available with a single click) to be able to interact with guarded routes.*


## Technologies
* [Gin](https://gin-gonic.com)
* [GORM](https://gorm.io)
* [PostgreSQL](https://postgresql.org)
* [Wire](https://github.com/google/wire)
* [Testify](https://github.com/stretchr/testify)

## Features
* Full integration testing coverage.
* GitHub Actions for CI/CD.
* Service-based architecture via Dependency Injection.
* CRUD operations.
* Database interactions via Repositories.
* Cookie-based session management using JWT.
* Roles management.
* Routes guarded with [AuthMiddleware](https://github.com/vladlent-portfolio/food-ordering-backend/blob/main/controllers/user/middlewares.go#L22).
* File upload with MIME-type and size check using [Upload](https://github.com/vladlent-portfolio/food-ordering-backend/blob/main/services/upload.go#L12) service.
* Model constraints.
* Validation for user-provided data.

## Running locally
The project uses Bash scripts for process automation.  
All of them are located in [scripts folder](https://github.com/vladlent-portfolio/food-ordering-backend/tree/main/scripts).

### Development server
Use `run.sh` script to run a local server in development mode. Make sure to update [.env][.env link] 
file with correct database credentials.  
By default, the server runs on port 8080. The port can be changed in [.env][.env link] file by changing `HOST_PORT` variable.

### Test
Use `test.sh` script to recursively run tests in all project's folders.

### Build
Use `build.sh` script to compile the app. By default, it creates a 64-bit executable for Linux named `food_ordering_app` in the current working directory.
You can customize build settings by updating `os`, `arch` and `outputname` variables in the script itself.  
Valid combinations of `os` and `arch` can be found in the [official golang docs](https://golang.org/doc/install/source#environment).

### Running in prod mode
In a directory where you are going to run the binary, create a file named `.production.env`. It should have the same structure as 
[.env][.env link] file, so you can just copy it. Update all variables in `.production.env` to your production credentials.

To run the app in prod mode you will need to set `GIN_MODE=release` environment variable in your terminal.

```bash
$ GIN_MODE=release ./food_ordering_api
```


[.env link]: https://github.com/vladlent-portfolio/food-ordering-backend/blob/main/.env
