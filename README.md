# ITU-DevOps

[![Build Status](https://travis-ci.org/ITU-DevOps-Team/ITU-DevOps.svg?branch=development)](https://travis-ci.org/ITU-DevOps-Team/ITU-DevOps)

Course repository

## Requirements

- Go (1.15)
- Docker
- Docker-compose
- Sqlite3

## Installing the depencies

```bash
  // ui app
  go get ./itu-minitwit-go/

  // api
  go get ./itu-minitwit-api/
```

## Running the server

```bash
  # ui app
  # build the executable
  go build ./itu-minitwit-go -o minitwit

  # running the executable
  ./minitwit
```

```bash
  # api
  # build the executable
  go build ./itu-minitwit-api -o minitwit-api

  # running the executable
  ./minitwit-api
```

## Creating a Docker image

```bash
  # ui app
  docker build -t minitwit:latest -f .deploy/itu-minitwit-go/Dockerfile .

  # api
  docker build -t minitwit-api:latest -f .deploy/itu-minitwit-api/Dockerfile .
```

## Running Docker containers with `docker-compose`

```bash
  # initialise all the services
  docker-compose up

  # intialise minitwit service
  docker-compose up minitwit

  # intialise minitwit-api service
  docker-compose up minitwit-api
```
