#!/bin/bash

export OUTPUT_DIR="/go/src/github.com/d-sparks/gravy/scrapers/alphavantage/output"
export DOCKER_HOST=tcp://0.0.0.0:2375
/usr/bin/docker run \
  --hostname=${DOCKER_HOST} \
  --volume=/root/alphavantage:${OUTPUT_DIR} \
  --env-file=/root/secrets \
  --rm \
  --name=alphavantage-scrape \
  alphavantage-scraper
