FROM golang:1.15.0-alpine3.12
COPY . /app
WORKDIR /app
RUN apk add --no-cache git
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOROOT/bin
ENV GOBIN=$GOPATH/bin
ENV GIN_MODE=release
RUN echo ${GOPATH} && echo ${GOROOT} && echo ${PATH} && pwd && go version
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN go get ./...
RUN go build
RUN printf "#!/bin/sh\n/app/app" > ./entrypoint.sh
RUN chmod 777 ./entrypoint.sh
RUN pwd && ls
WORKDIR /
ENTRYPOINT ["/app/entrypoint.sh"]
