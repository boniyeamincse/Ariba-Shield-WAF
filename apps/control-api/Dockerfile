# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/control-api .

# Runtime stage (Alpine — contains wget for healthcheck)
FROM alpine:3.20
RUN apk add --no-cache wget ca-certificates \
    && addgroup -S nonroot && adduser -S nonroot -G nonroot
COPY --from=build /bin/control-api /usr/local/bin/control-api
EXPOSE 8443
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/control-api"]
