#Build stage
FROM golang:1.23-alpine3.20 AS BuildStage

WORKDIR /app/cli

COPY . .

RUN go mod download

RUN go build -o /app/bin /cm/cli/main.go

# Executable stage

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /app/cli /app/cli


