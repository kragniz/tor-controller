FROM alpine
RUN apk update && apk add tor && mkdir -p /run/tor/service
ADD ./tor-daemon/entrypoint.sh /
RUN chmod +x /entrypoint.sh
WORKDIR /app
ENTRYPOINT ["/entrypoint.sh"]
