# Spanner
### Intelligent code generator for REST APIs and serverless applications
 [![asciicast](https://asciinema.org/a/rBDloIglqTWvcOT7qDGg1HUAp.svg)](https://asciinema.org/a/rBDloIglqTWvcOT7qDGg1HUAp)
 
 Spanner is a CLI tool to auto-generate REST APIs from just a simple definition of the object model. Spanner can setup the handlers, the repository and database functions and can also prewire middlewares for SSO based authentication and CSRF prevention. The templates can be modified to suit your needs. 

 Spanner connects to MongoDB out-of-the-box but it also generates a repo interface which can be implemented to support other databases. It can convert an example model given in JSON format to Golang struct with type-safety.

 Spanner generates an [echo](https://echo.labstack.com/) server and thus supports a wide range of pre-made middlewares.

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
 OIDC can be enabled by creating a `config.yaml` file before running the `spanner` command.
```yaml
oauth:
    enable: True
    clientid: "<clientID>"
    clientsecret: "<clientSecret>"
    redirecturl: "<redirectUrl>"
    configurl: "<Eg: https://YOUR_DOMAIN/.well-known/openid-configuration>"
port: "7000"
dockerusername: "<docker username>"
```
 then run
 ```bash
 spanner example.json
 ```

Right now, if OIDC is enabled then authentication is enabled for all endpoints. The client must be authenticated at all times. This can be fixed by setting route level middleware in `main.go`

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
Additional business logic can be added in `handler/handler.go`.

### Deploy
`Spanner` autogenerates **KNative** deployment file which can be used to deploy the REST API to a Kubernetes cluster. The `Makefile` also has a few recipes like `docker-build`, `docker-push`, `docker-buildx` to build and push docker images to DockerHub. `kdeploy` and `kdeploy-rpi` can be used to build, push the image and deploy to knative in a single command.
The buildx recipe and be used to deploy to Raspberry Pi too.  