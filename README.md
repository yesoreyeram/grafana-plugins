# Grafana Plugins

My experiments with [Grafana](https://grafana.com) [plugins](https://grafana.com/plugins) eco-system. This is a mono-repo of following plugins and grafana related packages.

## Included plugins

- [yesoreyeram-blank-datasource](./plugins/yesoreyeram-blank-datasource/)
- [yesoreyeram-hello-datasource](./plugins/yesoreyeram-hello-datasource/)
- [yesoreyeram-petstore-datasource](./plugins/yesoreyeram-petstore-datasource/)
- [yesoreyeram-vercel-datasource](./plugins/yesoreyeram-vercel-datasource)

## Included go packages

- [macros](./lib/go/macros/)
- [anyframer](./lib/go/anyframer/)
- [restds](./lib/go/restds/)

## Pre-Requisites

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
