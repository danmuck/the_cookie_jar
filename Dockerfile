FROM golang:latest

ENV HOME=/root
WORKDIR /root

COPY . .

# https://docs.docker.com/build/cache/optimize/

RUN --mount=type=cache,target=GOCACHE \
    go build -o the_cookie_jar

EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait 
RUN chmod +x /wait

CMD /wait && ./the_cookie_jar