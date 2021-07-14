# Backend for Food Ordering App

![App Image](https://food-ordering.app/assets/img/sm_image.png) 

## Table of contents
* [Overview](#overview)
* [Technologies](#technologies)
* [Features](#features)


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
* Full e2e test coverage.
* Service-based architecture via Dependency Injection.
* CRUD operations.
* Database interactions via Repositories.
* Cookie-based session management using JWT.
* Roles management.
* Routes guarded with [AuthMiddleware](https://github.com/vladlent-portfolio/food-ordering-backend/blob/main/controllers/user/middlewares.go#L22).
* File upload with MIME-type and size check using [Upload](https://github.com/vladlent-portfolio/food-ordering-backend/blob/main/services/upload.go#L12) service.
* Model constraints.
* Validation for user-provided data.

