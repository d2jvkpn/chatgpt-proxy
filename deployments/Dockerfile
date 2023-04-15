FROM golang:1-alpine AS builder
LABEL stage=chatgpt-proxy_builder

ARG ldflags=${ldflags}
ENV TZ="Asia/Shanghai"
WORKDIR /opt/chatgpt-proxy

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache update &> /dev/null && \
  apk --no-cache upgrade &> /dev/null

COPY ./main.go ./go.mod ./go.sum ./project.yaml  /opt/chatgpt-proxy/
COPY ./internal /opt/chatgpt-proxy/internal
# COPY ./pkg      /opt/chatgpt-proxy/pkg
# COPY ./vendor   /opt/chatgpt-proxy/vendor

# in alpine, date doesn't parse %:z
RUN go env -w GOPROXY="https://goproxy.cn,direct" && \
  go build -o main -ldflags="-w -s ${ldflags}" main.go

####
FROM alpine

ENV TZ="Asia/Shanghai"
WORKDIR /opt/chatgpt-proxy

COPY --from=builder /opt/chatgpt-proxy/main /opt/chatgpt-proxy/main

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache update &> /dev/null && \
  apk --no-cache upgrade &> /dev/null && \
  apk --no-cache add tzdata &> /dev/null

EXPOSE 3020
CMD ["./main", "--addr=:3020", "--config=configs/prod.yaml", "--release"]
