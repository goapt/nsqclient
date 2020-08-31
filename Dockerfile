FROM nsqio/nsq
LABEL maintainer="caixudong <fifsky@gmail.com>"

ENTRYPOINT ["/nsqd"]