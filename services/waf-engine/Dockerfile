# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/waf-engine .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /bin/waf-engine /usr/local/bin/waf-engine
COPY rules/core/baseline.conf /rules/baseline.conf
EXPOSE 8082
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/waf-engine"]