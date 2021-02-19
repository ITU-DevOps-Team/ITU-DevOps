FROM golang:1.13 AS build

WORKDIR /src
# COPY go.mod /src/
# COPY go.sum /src/
# COPY deps.go /src/
# COPY minitwit.go /src/

COPY ./itu-minitwit-go/go.mod /src/
COPY ./itu-minitwit-go/go.sum /src/
COPY ./itu-minitwit-go/deps.go /src/
COPY ./itu-minitwit-go/minitwit.go /src/

RUN mkdir -p /db_backup/
COPY ./db_backup/minitwit.db /db_backup/

RUN go get -d -v /src/...
RUN go install -v /src/...
RUN go build -o /src/minitwit
COPY ./itu-minitwit-go/ /src/

EXPOSE 8080
CMD [ "./minitwit" ]
