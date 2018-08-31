# Tag Service

[![CircleCI](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/tag-service.svg?style=svg&circle-token=30939c85b1ffb3af814b26f073ceee88b7899956)](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/tag-service)
[![codecov](https://codecov.io/gh/BetaProjectWave/tag-service/branch/master/graph/badge.svg?token=eRaLXqZqxI)](https://codecov.io/gh/BetaProjectWave/tag-service)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/bcddc93c687546cc85aaccd61be64f8d)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=BetaProjectWave/tag-service&amp;utm_campaign=Badge_Grade)

## Checkout the project

In GO, it matters where you set your projects - see GOPATH explanation below.
1. Create the following folder in your <WAVE_PROJECT_HOME>:  

```bash
mkdir -p <$WAVE_PROJECT_HOME>/go/src/github.com
```
2. Checkout this project in the above folder.

## Setting up GO environment and package management tool

1. Install GO: `brew install go`
2. Install DEP `brew install dep`
3. Set the GOPATH env variable and add it to the path in your bash profile (.profile, .zshrc): 
```bash
export GOPATH=$WAVE_PROJECT_HOME/go
export PATH=$PATH:$GOPATH/bin
```
    
A good explanation of how to set up GO projects can be found [here](https://golang.org/doc/code.html#Workspaces)
More explanation in setting up GOPATH env var can be found [here](https://golang.org/doc/code.html#GOPATH)

## Build tools

We are using **DEP** tool to manage go dependencies. You can read more about it [here](https://golang.github.io/dep)  

- Adding the project for the first time:    
  ```bash 
  dep ensure -vendor-only 
  ```

- Adding a new dependency 
  ```bash 
  dep ensure -add <dependecy path> 
  ```

## Build   


### Go   

```bash
go build -v
go install -v
```

  ### Docker  

```bash
docker build -t projectwave/tag-service .
```

## Run the tests

```bash
go test ./... -v
```

  ### Generate test coverage

```bash
./scripts/run_test.sh
```

  ## Run the server  

### Go  

```bash
go run main.go
```

### Docker 

```bash
docker run -p 8080:8080 projectwave/tag-service
```
## Format the code

GO comes out with formatting tool out of the box. There are three options:

1. To run the tool manually: `go fmt <name of the file>` or `gofmt -w -l .` for all files.
2. Using git pre commit hook that will flags all the unformatted files for you. Here are the steps: 
   * Create file named: `pre-commit` in `.git/hooks` of the current projects  with the content below  
   * Run: `chmod +x ./git/hooks/pre-commit`

    ```bash
    #!/bin/sh
    gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
    [ -z "$gofiles" ] && exit 0
    
    unformatted=$(gofmt -l $gofiles)
    [ -z "$unformatted" ] && exit 0
        
    echo >&2 "Go files must be formatted with gofmt. Running go fmt..."
    for fn in $unformatted; do
     gofmt -w $PWD/$fn
    done
    
    exit 1
    ```

3. If you use IntelliJ, there is a plugin called **File Watchers** that handles the format of different languages such as
 **Go** or **terraform**. See https://www.jetbrains.com/help/idea/using-file-watchers.html
 
## Swagger

We are using gin-swagger (https://github.com/swaggo/gin-swagger) - see the comments added in each endpoint handler (`api/handler.go`)

To generate a newer version of the operations to be exposed. Run these steps:
1. Install the swag tool using this script: `./scripts/install_swag.sh`
2. Run `swag init`. This will generate a new swagger docs file see - `./docs/docs.go` 

You can access swagger documentation: http://localhost:8080/swagger/index.html

## Mocking
We are using Go Mock module for generating mocks. A good introduction of Go Mock can be found [here](https://blog.codecentric.de/en/2017/08/gomock-tutorial/).


