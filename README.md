# Hitched API

## Quickstart

```shell
# Launch API
go run main.go router.go
```

# Developing

```shell
# Updating vendor files
govendor add +external

# Deploying to Heroku (after committing changes)
heroku accounts:set kc
git push heroku master
```
