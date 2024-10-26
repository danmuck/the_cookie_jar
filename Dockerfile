FROM golang:alpine AS builder

ENV HOME=/root
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# https://docs.docker.com/build/cache/optimize/
COPY . .

RUN --mount=type=cache,target=GOCACHE \
     go build -o the_cookie_jar

# RUN --mount=type=cache,target=GOCACHE \
#     go build cmd/client/server.go

FROM alpine
WORKDIR /root
COPY --from=builder /build .

EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait 
RUN chmod +x /wait


CMD /wait && ./the_cookie_jar
