FROM debian:latest

RUN mkdir /app
WORKDIR /app

ADD config.yml /app/
ADD dmon /app/

CMD ["./dmon"]
