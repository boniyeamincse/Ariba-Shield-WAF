FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/event-ingestor .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /bin/event-ingestor /usr/local/bin/event-ingestor
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/event-ingestor"]