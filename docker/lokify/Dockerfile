#######################
## STAGE: ui-builder ##
#######################
FROM node:14.8.0-stretch as ui-builder

RUN mkdir -p /opt/lokify/ui
WORKDIR /opt/lokify/ui

COPY ui/package.json ./
COPY ui/package-lock.json ./
RUN npm ci

COPY ui ./
RUN npm run build

#######################
## STAGE: api-builder ##
#######################
FROM golang:1.16.4-buster

WORKDIR /go/src/github.com/flared/lokify
RUN go get -d -v golang.org/x/net/html  
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

RUN mkdir -p /opt/lokify
WORKDIR /opt/lokify

COPY --from=ui-builder /opt/lokify/ui ui
