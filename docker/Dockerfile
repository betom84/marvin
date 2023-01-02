FROM golang:1.19 as build

WORKDIR /build

COPY . .

RUN go mod download && go mod verify
RUN CGO_ENABLED=0 go build -v -o marvin

## Deploy
FROM alpine:latest

WORKDIR /opt/marvin
ADD build/ .

COPY --from=build /build/marvin .

RUN chmod u+x ./marvin

EXPOSE 8080 6443

#VOLUME /opt/marvin/config

#ENTRYPOINT ["/bin/sh"]
ENTRYPOINT ["/opt/marvin/marvin"]
CMD ["-config", "config/config.json"]