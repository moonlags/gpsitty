# Project gpsitty

My own server implementation for gps trackers found on aliexpress

The original program is named "365GPS"

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### You need to have

1. Go >= 1.21 installed
2. Nodejs installed
3. MySQL database running

build the application

```bash
go build cmd/api/main.go
```

run the application

```bash
go run cmd/api/main.go > /logs/logs.log &
cd client
npm run dev
```

run the test suite

```bash
go test ./tests -v
```
