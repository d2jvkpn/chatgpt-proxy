FROM golang:1 AS builder
LABEL stage=chatgpt-proxy_builder

ARG ldflags=${ldflags}
ENV TZ="Asia/Shanghai"
WORKDIR /opt/chatgpt-proxy

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
# COPY ./pkg     ./pkg
# COPY ./vendor  ./vendor

# in alpine, date doesn't parse %:z
RUN go env -w GOPROXY="https://goproxy.cn,direct" && \
  go build -o main -ldflags="-w -s ${ldflags}" main.go

####
FROM python:3.10

ENV TZ="Asia/Shanghai"

# RUN apt update && apt -y upgrade

# ubuntu
RUN useradd -m -u 1000 -s /bin/bash hello
# alpine
# RUN adduser -D -u 1000 hello
USER hello
ENV PATH="/home/hello/.local/bin:$PATH"
WORKDIR /home/hello/chatgpt-proxy

ADD ./deployments/pip.txt  ./
COPY --from=builder /opt/chatgpt-proxy/main ./main

RUN pip3 config set global.index-url 'https://pypi.douban.com/simple/' && \
  pip3 config set install.trusted-host 'pypi.douban/simple' &&  \
  pip3 install --no-cache -r pip.txt

EXPOSE 3020
CMD ["./main", "--addr=:3020", "--config=configs/prod.yaml", "--release"]
