# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/control-api .

# Runtime stage (distroless for minimal attack surface)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /bin/control-api /usr/local/bin/control-api
EXPOSE 8443
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/control-api"]
