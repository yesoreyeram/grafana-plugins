# Contributing Guide

## Pre-requisites

- Git
- NodeJS
- Yarn
- Go 1.19
- Mage
- Docker

## How to run this locally

- Clone the repo locally `git clone https://github.com/yesoreyeram/grafana-plugins yesoreyeram-grafana-plugins` and `cd yesoreyeram-grafana-plugins`
- Install packages `yarn install`
- Build frontend and backend `yarn build:full`
- Start Grafana `docker compose up`
- Open Grafana `http://localhost:3000` and login with credentials `grafana:grafana`
- Visit `http://localhost:3000/datasources` to see the list of datasources
