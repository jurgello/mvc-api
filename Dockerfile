
## Build
FROM golang:latest AS build

RUN mkdir /app
WORKDIR /app

COPY assets /app/assets
COPY templates /app/templates

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /myapp

## Deploy
FROM golang:latest 

WORKDIR /

COPY --from=build /myapp /myapp

EXPOSE 8080

ENTRYPOINT ["./myapp"]

