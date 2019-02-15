# Hitched API

## Quickstart

```shell
# Launch API
./run.sh
```

# Developing

```shell
# Updating vendor files
govendor add +external

# Create a local environment file before running the API locally
cp ./.env.sample ./.env

# Running the API
./run.sh

# Deploying to Heroku (after committing changes)
./deploy.sh

# Tailing the logs
heroku logs --tail
```
