# 前端项目Docker部署指南

本目录包含用于将前端静态文件部署到Docker容器中的配置文件。

## 文件说明

- `Dockerfile`: 用于构建Docker镜像的配置文件
- `nginx.conf`: Nginx服务器的配置文件，支持SPA应用的路由

## 使用方法

### 前提条件

- 已安装Docker
- 已构建好前端项目，生成静态文件

### 构建镜像

在前端项目构建完成后，将生成的静态文件放入此目录，然后执行以下命令构建Docker镜像：

```bash
# 进入build目录
cd /path/to/your/project/build

# 构建Docker镜像
docker build -t frontend-app:latest .
```

### 运行容器

构建完成后，可以使用以下命令运行容器：

```bash
# 运行容器，将容器的80端口映射到主机的8080端口
docker run -d -p 8080:80 --name frontend-app frontend-app:latest
```

现在可以通过访问 `http://localhost:8080` 来访问你的前端应用。

### 使用自定义配置

如果需要使用自定义的Nginx配置，可以取消Dockerfile中的注释：

```dockerfile
# 复制自定义的nginx配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

## 注意事项

- 确保构建前端项目时，所有的资源路径都是相对路径
- 如果应用需要与后端API通信，可能需要在Nginx配置中添加反向代理设置