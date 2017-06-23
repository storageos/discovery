
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY discovery /bin/discovery

ENTRYPOINT ["/bin/discovery"]

EXPOSE 8081