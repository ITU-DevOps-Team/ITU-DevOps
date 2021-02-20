# ITU-DevOps

Course repository

## Requirements

- Go (1.15)
- Docker
- Docker-compose
- Sqlite3

## Installing the depencies

```bash
  go get ./itu-minitwit-go/
```

## Running the server

```bash
  # build the executable
  go build ./itu-minitwit-go -o minitwit

  # running the executable
  ./minitwit
```

## Creating a Docker image

```bash
  docker build -t minitwit:latest -f .deploy/itu-minitwit-go/Dockerfile .
```

## Running Docker containers with `docker-compose`

```bash
  # initialise all the services
  docker-compose up

  # intialise minitwit service
  docker-compose up minitwit
```
