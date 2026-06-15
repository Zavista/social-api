
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /bin/api ./cmd/api


FROM --platform=$TARGETPLATFORM alpine:3.21

# ca-certificates: needed for outbound HTTPS (SendGrid, RDS TLS, etc.)
RUN apk add --no-cache ca-certificates

COPY --from=build /bin/api /bin/api

EXPOSE 8080

ENTRYPOINT ["/bin/api"]
