FROM golang:1.15.0-alpine3.12
COPY . /app
WORKDIR /app
RUN apk add --no-cache git
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOROOT/bin
ENV GOBIN=$GOPATH/bin
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN go get ./...
RUN go build -o ./dist main.go
RUN printf "#!/bin/sh\n/app/dist/main" > ./entrypoint.sh
RUN chmod 777 ./entrypoint.sh
WORKDIR /
ENTRYPOINT ["/app/entrypoint.sh"]
