# WORK IN PROGRESS
# Github Client using OAuth2.0

A simple github web app for performing git operations, built over go-github package utilizing github REST APIs, OAuth2.0 for authroization. Idea is to demo list of repos for authenticated user, create a new branch, create a commit, using which create a new PR in newly created branch.

To authorize user with GitHub using OAuth2.0, register the application and generate client_id, client_secret and keep them safely, specially the secret. No where we are going to store the client id or secret in the code.

Refer [Authorizing OAuth Apps](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow)

Register new OAuth Apps [here](https://github.com/settings/applications/new)


## Temporarily some manual changes needed
Currently App doesnt have enough frontend html to take input on repo name, and modified file to consider for a new commit, so need to update same in code [('sourceRepo','prRepo'), ('fileForAutoCommit')] to make it work for any user (i know its dirty, though need to take care of this in future only). Until then, following git diff should help to make necessary changes quickly. 
```
    cmd$ git diff
    diff --git a/cmd/githubOperations.go b/cmd/githubOperations.go
    index ab72e60..197012e 100644
    --- a/cmd/githubOperations.go
    +++ b/cmd/githubOperations.go
    @@ -21,7 +21,7 @@ type Package struct {
            LastUpdatedBy string
    }
    
    -var fileForAutoCommit = "commitTimeStamps.txt"
    +var fileForAutoCommit = "tmp.txt"
    
    func editSampleFile(file string, data string) bool {
            fmt.Printf("\nEditing file for commit: %s\n", file)
    diff --git a/cmd/main.go b/cmd/main.go
    index bbd6602..d6d9b37 100644
    --- a/cmd/main.go
    +++ b/cmd/main.go
    @@ -58,9 +58,9 @@ func callbackHandler(w http.ResponseWriter, r *http.Request) {
            //fmt.Println(*user)
    
            *sourceOwner = *user.Login
    -       *sourceRepo = "gomods" //DEFAULT, need to take as i/p ideally
    +       *sourceRepo = "dummyrepo" //DEFAULT, need to take as i/p ideally
            *prRepoOwner = *user.Login
    -       *prRepo = "gomods" //DEFAULT
    +       *prRepo = "dummyrepo" //DEFAULT
            *authorName = *user.Name
            *authorEmail = "dummy@gmail.com"
```

### Additional care
Also make sure to ***keep Repo description populated*** (About section in repo page), or app may run into errors.

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


