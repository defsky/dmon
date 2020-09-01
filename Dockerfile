FROM debian:latest

# set local timezone
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# init work dir
RUN mkdir /app
WORKDIR /app

ADD config.yml /app/
ADD dmon /app/

CMD ["./dmon"]
