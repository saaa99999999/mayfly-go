ARG BASEIMAGES=m.daocloud.io/docker.io/alpine:3.20.2

FROM ${BASEIMAGES} AS builder
ARG TARGETARCH

ARG MAYFLY_GO_VERSION
ARG MAYFLY_GO_DIR_NAME=mayfly-go-linux-${TARGETARCH}
ARG MAYFLY_GO_URL=https://gitee.com/dromara/mayfly-go/releases/download/${MAYFLY_GO_VERSION}/${MAYFLY_GO_DIR_NAME}.zip

RUN wget -cO mayfly-go.zip ${MAYFLY_GO_URL} && \
    unzip mayfly-go.zip && \
    mv ${MAYFLY_GO_DIR_NAME} /opt/mayfly-go && \
    rm -rf mayfly-go.zip

FROM ${BASEIMAGES}

ARG TZ=Asia/Shanghai
ENV TZ=${TZ}
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


# 从 builder 阶段复制完整目录
COPY --from=builder /opt/mayfly-go/bin/mayfly-go /usr/local/bin/mayfly-go

# 设置执行权限
RUN chmod +x /usr/local/bin/mayfly-go

WORKDIR /mayfly-go

EXPOSE 18888

CMD ["mayfly-go"]
