# ![RealWorld Example App](logo.png)

> ### Go (golang) codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged backend application built with **GO** (golang) including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the golang community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

There are three layers: handlers, domain and model.

Handlers layer controls routing and retrieving user inputs: serialization of requests and query parameters. Gorilla mux is used for routing.

Domain layer controls what can be called as our application's core business-logic. Functions are named from user's perspective. Types represent responses that are defined in specs.  There are extensive tests for this layer.

Model layer represents how the data is stored persistently. Postgres db is used, with GORM for migrations and some queries. Types represent tables in the database.

There are utils packages for auth, errors and app initialization.

# Getting started

To start the service
`docker-compose up -d --build`

To check that it works
`curl localhost:4000/ping`

To run the test suite
`./test_requests.sh`
