FROM golang:1-alpine AS builder
LABEL stage=chatgpt-proxy_builder

ARG ldflags=${ldflags}
ENV TZ="Asia/Shanghai"
WORKDIR /opt/chatgpt-proxy

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache update &> /dev/null && \
  apk --no-cache upgrade &> /dev/null && \
  apk --no-cache add tzdata

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
# COPY ./pkg     ./pkg
# COPY ./vendor  ./vendor

# in alpine, date doesn't parse %:z
RUN go env -w GOPROXY="https://goproxy.cn,direct" && \
  go build -o main -ldflags="-w -s ${ldflags}" main.go

####
FROM alpine

ENV TZ="Asia/Shanghai"

# RUN useradd -m -u 1000 -s /bin/bash hello
RUN adduser -D -u 1000 hello
USER hello
ENV PATH="/home/hello/.local/bin:$PATH"
WORKDIR /home/hello/chatgpt-proxy

COPY --from=builder /opt/chatgpt-proxy/main ./main

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache update &> /dev/null && \
  apk --no-cache upgrade &> /dev/null && \
  apk --no-cache add tzdata python3 &> /dev/null

RUN python3 -m ensurepip && \
  pip3 config set global.index-url 'https://pypi.douban.com/simple/' && \
  pip3 config set install.trusted-host 'pypi.douban/simple' &&  \
  pip3 install ipython && \
  pip3 install pip.txt && \
  rm -r /home/hello/.cache/

EXPOSE 3020
CMD ["./main", "--addr=:3020", "--config=configs/prod.yaml", "--release"]
