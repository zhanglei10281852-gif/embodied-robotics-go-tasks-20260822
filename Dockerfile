FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/robotics ./cmd/server
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/robotics /app/robotics
COPY migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/app/robotics"]
