# Stage 1: Build
FROM golang:1.12.9 AS builder

RUN mkdir /build
COPY . /build/
WORKDIR /build/
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# Stage 2: Final Image
FROM scratch
WORKDIR /root/
COPY --from=builder /build/app .
# Serves up the UI portion of the web app.
COPY --from=builder /build/static ./static

EXPOSE 8000
ENTRYPOINT ["./app"]
