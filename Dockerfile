FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 10001 -S wwb && \
    adduser -u 10001 -S wwb -G wwb && \
    mkdir -p /app/data && \
    chown -R wwb:wwb /app

WORKDIR /app

COPY bin/wwb /app/wwb

USER 10001:10001
EXPOSE 8080

ENTRYPOINT ["/app/wwb"]

