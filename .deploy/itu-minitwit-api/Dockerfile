FROM golang:1.15

WORKDIR /src
COPY ./itu-minitwit-api/go.* /src/
COPY ./itu-minitwit-api/*.go /src/

RUN mkdir -p /db_backup/
COPY ./db_backup/minitwit.db /db_backup/

RUN go get -d -v /src/...
RUN go install -v /src/...
RUN go build -o /src/minitwit-api
COPY ./itu-minitwit-api/ /src/

EXPOSE 8080
CMD [ "./minitwit-api" ]
