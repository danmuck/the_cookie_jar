FROM golang:latest

ENV HOME=/root
WORKDIR /root

COPY . .

# https://docs.docker.com/build/cache/optimize/

# RUN --mount=type=cache,target=GOCACHE \
#     go build -o the_cookie_jar

RUN --mount=type=cache,target=GOCACHE \
    go build cmd/client/server.go

EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait 
RUN chmod +x /wait

CMD /wait && ./server
