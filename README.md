# Github Client using OAuth2.0

A simple github web app for performing git operations, built over go-github package utilizing github REST APIs, OAuth2.0 for authroization.

To authorize user with GitHub using OAuth2.0, register the application and generate client_id, client_secret and keep them safely, specially the secret. No where we are going to store the client id or secret in the code.

Refer https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow

Use https://github.com/settings/applications/new


## Build
Build docker container with supplied Dockerfile in parent dir

    cd GithubClient
    docker build -t ghclient .


## Running the docker container

### Pre-requisites
Make sure to export env variables in current terminal enviroment already

    export CLIENT_ID_ENV=\<GITHUB-REGISTERED-APP-OAUTH2-CLIENT-ID\>

    export CLIENT_SEC_ENV=\<GITHUB-REGISTERED-APP-OAUTH2-CLIENT-SECRET\> 


### Running docker container with pre-built image
    docker run -dp <hostname:port:<container-port>> -e CLIENT_ID_ENV -e CLIENT_SEC_ENV ghclient

### Example
    docker run -dp 127.0.0.1:9080:9080 -e CLIENT_ID_ENV -e CLIENT_SEC_ENV ghclient

*Note: Intentionally values of env variables not supplied in docker run command line, so its not visible in 'docker ps' o/p*

## Accessing Web Interface
Load <hostname:port> or default *127.0.0.1:9080*, with which container started (as used in 'docker run' command above) and start interacting with available options.


## References

- https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/best-practices-for-creating-an-oauth-app
- https://docs.github.com/en/rest?apiVersion=2022-11-28
- https://pkg.go.dev/github.com/google/go-github
- https://github.com/google/go-github


