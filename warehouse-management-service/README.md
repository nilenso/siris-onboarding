# Warehouse Management System

Refer [project document](https://docs.google.com/document/d/19Gtx_CG7ebD9HDvTI_KJS9Ge-hfUYPUMx85ZurVgLXs/edit?usp=sharing)

This application is a web service, with reasonable domain complexity. \
The project is intended to demonstrate/achieve:

1. Relational database modelling
2. Project management
3. CI/CD using GitHub Actions and a shell script

## Installation

1. Clone the repo `https://github.com/nilenso/siris-onboarding.git`
2. `cd` into `warehouse-management-system`
3. Download dependencies `go mod download`
4. Setup database: _TODO_
5. Run tests: _TODO_
6. Start the server `go run main.go`,
    1. Starts the server on `http://localhost:80/`
    2. Status check route is available at `http://localhost:80/ping`

> This web server uses the [Gin framework](https://github.com/gin-gonic/gin)
