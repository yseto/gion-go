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
COPY app/           /usr/src/app/app/
COPY cmd/           /usr/src/app/cmd/
COPY config/        /usr/src/app/config/
COPY db/            /usr/src/app/db/
COPY internal/      /usr/src/app/internal/
ENV CGO_ENABLED=0
RUN go build -o /app/gion       ./cmd/app
RUN go build -o /app/queueing   ./cmd/queueing/
RUN go build -o /app/worker     ./cmd/worker/
RUN go build -o /app/insertuser ./cmd/insertuser/

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app/
COPY --from=build-go --chown=nonroot:nonroot /app/gion /app/queueing /app/worker /app/insertuser /app/
COPY public/ /app/public/
COPY --from=build-env --chown=nonroot:nonroot /app/public/gion.js /app/public/

EXPOSE 8080

