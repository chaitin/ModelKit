# 使用官方的nginx镜像作为基础镜像，指定x86架构
FROM --platform=linux/amd64 nginx:stable-alpine

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 删除默认的nginx静态文件
RUN rm -rf ./*

# 复制构建好的前端静态文件到nginx的默认静态文件目录
COPY ./dist /usr/share/nginx/html

# 复制自定义的nginx配置文件
COPY ./build/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露80端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]