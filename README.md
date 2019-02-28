# Hitched API

- [Getting started](#getting-started)
- [Running the API](#running-the-api)

# Getting started

Make sure [Go](https://golang.org) is installed on your machine:

```shell
go version
```

Then clone the repository into your `$GOPATH` and install the dependencies:

```shell
# Prepare the kwcay directory inside your GOPATH.
mkdir --parents $GOPATH/src/github.com/kwcay && cd $GOPATH/src/github.com/kwcay

# Clone the Git repository.
git clone git@github.com:kwcay/hitched-api.git
cd hitched-api

# Install the project dependencies.
go get -v
```

Finally, create a local environment file:

```shell
cp .env.sample .env
```

# Running and deploying the API

```shell
# Use the run script to launch API on your machine.
./run.sh

# Updating the vendor files...
govendor add +external

# Deploying to Heroku (after committing changes)...
./deploy.sh

# Tailing the logs from Heroku...
heroku logs --tail
```
