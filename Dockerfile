FROM ubuntu:groovy

RUN mkdir -p /usr/share/man/man1

RUN apk update && apk add --no-cache git ca-certificates

RUN curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz

RUN tar xvf go1.6.linux-amd64.tar.gz

RUN apt-get -yq install plantuml graphviz git fonts-ipafont fonts-ipaexfont && rm -rf /var/lib/apt/lists/*

COPY . .

ENV GOPATH=""

RUN go mod download

RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/app

ENTRYPOINT ["/bin/app"]
