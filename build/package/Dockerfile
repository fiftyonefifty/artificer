# Stage 1: Build
FROM golang:1.12.9 AS builder

RUN mkdir /dist
COPY . /dist/
WORKDIR /dist/
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app ./cmd/artificer

# Stage 2: Final Image
FROM scratch as production
WORKDIR /root/
COPY --from=builder /dist/app .
# Serves up the UI portion of the web app.
COPY --from=builder /dist/static ./static

EXPOSE 8000
ENTRYPOINT ["./app"]
