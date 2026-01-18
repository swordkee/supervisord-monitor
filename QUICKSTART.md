# 快速开始指南

## 1. 克隆或复制项目

```bash
cd E:\work\php\supervisord-monitor-go
```

## 2. 安装依赖

### Windows

```bash
# 安装 Go 依赖
go mod download

# 安装 Node 依赖
cd frontend
npm install
cd ..
```

## 3. 配置服务器

编辑 `config.yaml` 文件，修改 Supervisord 服务器配置：

```yaml
supervisor_servers:
  your-server-name:
    url: "http://your-server-ip/RPC2"
    port: "9001"
    username: "yourusername"
    password: "yourpassword"
```

## 4. 构建项目

### Windows

```bash
build.bat
```

### 或使用 Make（推荐）

```bash
make windows
```

## 5. 运行

```bash
supervisord-monitor.exe
```

或指定配置文件：

```bash
supervisord-monitor.exe -config config.yaml
```

## 6. 访问

在浏览器中打开：`http://localhost:8080`

## 开发模式

### 同时运行前端和后端

```bash
make dev
```

或分别运行：

```bash
# 终端 1 - 前端开发服务器
cd frontend
npm run dev

# 终端 2 - 后端服务器
go run main.go
```

## 常见问题

### 1. Go 依赖下载失败

设置 Go 代理（中国大陆用户）：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 2. Node 依赖下载失败

设置 npm 镜像：

```bash
npm config set registry https://registry.npmmirror.com
```

### 3. 构建失败

确保已安装：
- Go 1.21 或更高版本
- Node.js 16 或更高版本

### 4. 无法连接到 Supervisord

- 确保 Supervisord 的 `inet_http_server` 已启用
- 检查防火墙设置
- 验证用户名和密码是否正确

## 目录说明

```
supervisord-monitor-go/
├── config/          # 配置管理代码
├── handlers/        # API 处理器
├── services/        # Supervisord 客户端
├── webembed/        # 静态文件嵌入
├── frontend/        # Vue 前端代码
├── sounds/          # 声音文件
├── config.yaml      # 配置文件
├── main.go          # 主程序
├── build.bat        # Windows 构建脚本
├── build.sh         # Linux/Mac 构建脚本
└── Makefile         # Make 构建脚本
```

## 更多信息

请查看 [README.md](README.md) 获取详细文档。
