# Spanner
### Intelligent code generator for REST APIs and serverless applications
 [![asciicast](https://asciinema.org/a/rBDloIglqTWvcOT7qDGg1HUAp.svg)](https://asciinema.org/a/rBDloIglqTWvcOT7qDGg1HUAp)
 
 Spanner is a CLI tool to auto-generate REST APIs from just a simple definition of the object model. Spanner can setup the handlers, the repository and database functions and can also prewire middlewares for SSO based authentication and CSRF prevention. The templates can be modified to suit your needs. 

 Spanner connects to MongoDB out-of-the-box but it also generates a repo interface which can be implemented to support other databases. It can convert an example model given in JSON format to Golang struct with type-safety.

 ## Features
 - [x] Generates REST APIs (CRUD) 
 - [x] Auto-generates swagger spec
 - [x] Auto-generates SSO middleware
 - [x] Follows RESTful API standards
 - [x] Has default middlewares for logging, CSRF prevention etc.
 - [x] Auto-generates MongoDB functions  
 - [x] Auto-generates makefile to launch the server and database
 #### WIP
 - [ ] SQL support
 - [ ] Expose SQL database queries as REST API
 - [ ] Auto-generate serverless config files for KNative deployments
 - [ ] Auto-generate docker-compose files for local development and testing
 - [ ] Support for Pagination and PATCH requests
 - [ ] Prometheus and Grafana support to gather metrics


 ## Example
 The user has to define an example model over which the CRUD API needs to be built and Spanner takes care of the rest.
 Save this as `example.json`
 ```json
 {
   "user": {
        "firstName": "Bruce",
        "lastName": "Wayne",
        "username": "darkknight",
        "age": 35,
        "designation": "Vigilante",
        "isAwesome": true,
        "address": {
            "street": "1007 Mountain Drive",
            "city": "Gotham",
            "location": {
                "lat": "23.67N",
                "long": "34.567E"
            },
            "isResidence": true
        }
    }
}
 ```
 and run
 ```bash
 spanner example.json
 ```
 ### Folder structure
```
 crud-app/
├─ go.mod
├─ example.json
├─ config.yaml
├─ main.go
├─ handler/
│  ├─ handler.go
│  ├─ middlewares.go
├─ modules/
│  ├─ user/
│  │  ├─ model/
│  │  │  ├─ model.gen.go
│  │  ├─ repo/
│  │  │  ├─ mongorepo.gen.go
│  │  │  ├─ repo.gen.go
├─ Makefile
```

