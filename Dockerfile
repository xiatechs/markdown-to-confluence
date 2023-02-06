FROM golang:1.18-alpine AS builder 

ENV PLANTUML_VERSION 1.2020.14
ENV LANG en_US.UTF-8

# Install git.
# Git is required for fetching the dependencies. - also, install java!
RUN apk update && apk add --no-cache git ca-certificates && apk add openjdk11

# Install plantuml dependancies
RUN apk add --no-cache graphviz font-droid font-droid-nonlatin curl \
    && apk add --no-cache \
        --repository https://nl.alpinelinux.org/alpine/edge/testing \
    && mkdir /app \
    && curl -L https://sourceforge.net/projects/plantuml/files/plantuml.${PLANTUML_VERSION}.jar/download -o /app/plantuml.jar \
    && apk del curl
   
COPY . .

# We don't need GOPATH so unset it. 
# If GOPATH is set go mod download will fail as it thinks there's a go.mod in the GOPATH
ENV GOPATH=""

#how to run the plantuml java app
#RUN java -jar /app/plantuml.jar -h

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/app

ENTRYPOINT ["/bin/app"]
