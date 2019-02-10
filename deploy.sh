#!/bin/sh

echo "Deploying Hitched API..."

# Set Heroku account
#HEROKU_ACCOUNT=$(heroku accounts:current)
heroku accounts:set kc

# Push to Heroku
git push heroku master

# Reset Heroku account
echo "TODO: reset Heroku account..."
