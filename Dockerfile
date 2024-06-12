# 设置基础镜像为golang的指定版本
FROM golang:1.22-alpine
# golang:1.22.3-alpine

RUN apk add --no-cache --update git build-base

# 将当前工作目录设置为/workdir
WORKDIR /workdir

# 设置 Go Module 的代理为 goproxy.cn
ENV GOPROXY=https://goproxy.cn

# 拷贝整个项目到/workdir目录
COPY . .

# 下载依赖（go.mod和go.sum文件必须在root目录）
RUN go mod download

# 构建项目（如果是在CGO环境下编译，确保相应的依赖项已安装）
RUN CGO_ENABLED=1 go build -o /workdir/postbee main.go

# 设置容器的默认命令
CMD ["/workdir/postbee"]