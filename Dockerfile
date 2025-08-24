FROM node:22-alpine AS build-env
WORKDIR /app/frontend/
COPY frontend/package*json /app/frontend/
RUN npm install
COPY frontend/ /app/frontend/
RUN npm run build

FROM golang:1.24-bookworm AS build-go

WORKDIR /usr/src/app/

COPY go.mod go.sum  /usr/src/app/
RUN go mod download

COPY fs.go     /usr/src/app/
COPY cmd/      /usr/src/app/cmd/
COPY internal/ /usr/src/app/internal/

COPY public/ /usr/src/app/public/
COPY --from=build-env /app/public/gion.js /usr/src/app/public/

RUN go build -o /app/gion       ./cmd/app
RUN go build -o /app/insertuser ./cmd/insertuser/

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app/
COPY --from=build-go --chown=nonroot:nonroot /app/gion /app/insertuser /app/

EXPOSE 8080

