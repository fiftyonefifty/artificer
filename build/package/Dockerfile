# Stage 1: Build
FROM golang:1.13 AS builder

RUN mkdir /dist
COPY . /dist/
WORKDIR /dist/
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app ./artificer-cobra

# Stage 2: Final Image
FROM alpine as production
RUN apk update && apk add bash
WORKDIR /root/
ENV AZURE_CLIENT_ID=ddecd689-3d6c-4311-b354-45ab6fa7e48c
ENV AZURE_CLIENT_SECRET={{SECRET}}
ENV AZURE_SUBSCRIPTION_ID=39ac48fb-fea0-486a-ba84-e0ae9b06c663
ENV AZURE_TENANT_ID=3b217a9b-6c58-428b-b022-5ad741ce2016

COPY --from=builder /dist/app .
# Serves up the UI portion of the web app.
 
COPY --from=builder /dist/config ./config
EXPOSE 9000
ENTRYPOINT ["./app"]
