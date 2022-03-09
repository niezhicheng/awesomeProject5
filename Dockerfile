FROM golang:alpine

# 为我们的镜像设置必要的环境变量(其中GOPROXY设置代理很重要,不然很多依赖无法下载)
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

# 移动到容器的工作目录：/build
WORKDIR /build

# 将代码复制到容器中(前面的点表示主服务器的当前位置也就是项目目录,后面的点表示容器的当前位置,也就是/build)
COPY . .

# 将我们的代码编译成二进制可执行文件app,这个可执行文件现在在容器的/build/app
RUN awesomeProject5

# 移动到用于存放生成的二进制文件的 /dist 目录
WORKDIR /dist

# 将二进制文件从 /build 目录复制到这里
RUN cp /build/app .

# 声明服务端口(并无实际作用,相当于一个注释给别人看的,所以即使项目中使用的端口不是8889也能跑起来)
EXPOSE 8080

# 启动容器时运行的命令
CMD ["/dist/app"]

