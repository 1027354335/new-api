FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget \
    && update-ca-certificates

COPY alipay-bridge /alipay-bridge
RUN chmod +x /alipay-bridge

EXPOSE 3001
WORKDIR /data

ENTRYPOINT ["/alipay-bridge"]
