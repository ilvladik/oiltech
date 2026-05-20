FROM golang:1.25-alpine AS build
ARG APP
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${APP}

FROM alpine:3.20
RUN adduser -D appuser
USER appuser
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]
