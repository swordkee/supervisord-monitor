# Supervisord Monitor (Go + Vue)

基于 Gin + Vue 复刻的 Supervisord 多服务器监控工具，将所有静态文件打包到单一可执行文件中。

 ## 特性

 - 🎯 实时监控多个 Supervisord 服务器
 - 🚀 启动/停止/重启单个或所有进程
 - 📊 实时显示进程状态和运行时间
 - 🔔 错误日志查看和声音警报
 - 🔇 静音功能（使用浏览器本地存储）
 - 🔄 自动刷新页面（可配置）
 - 🔐 支持认证和非认证的 Supervisord 服务器
 - 🔒 Web界面 HTTP Basic 认证（可选）
 - 📦 单一可执行文件，无需额外部署
 - ⚡ 高性能 Go 后端 + Vue 3 前端

## 技术栈

### 后端
- **Gin**: 高性能 Go Web 框架
- **XML-RPC**: 与 Supervisord RPC2 接口通信
- **Embed**: 将静态文件打包到二进制
- **Viper**: 配置文件管理

### 前端
- **Vue 3**: 渐进式 JavaScript 框架
- **Vite**: 下一代前端构建工具
- **Axios**: HTTP 客户端
- **Bootstrap 2**: 响应式 UI 框架
- **BootCDN**: 国内 CDN 加速访问

## 安装

### 前置要求

- Go 1.21+
- Node.js 16+ (仅构建时需要)
- Supervisord 服务器（需要启用 RPC2 接口）

### 配置 Supervisord 服务器

在每个 Supervisord 服务器的配置文件中启用 `inet_http_server`：

```ini
[inet_http_server]
port=*:9001
username="yourusername"
password="yourpass"
```

重启 Supervisord 服务：

```bash
sudo supervisorctl restart all
```

### 构建项目

#### Windows

```bash
build.bat
```

#### Linux/Mac

```bash
chmod +x build.sh
./build.sh
```

构建完成后会生成单一的可执行文件：
- Windows: `supervisord-monitor.exe`
- Linux/Mac: `supervisord-monitor`

## 配置

编辑 `config.yaml` 文件配置监控服务器：

```yaml
 supervisor_cols: 2      # 仪表板列数（2 或 3）
 refresh: 10             # 刷新间隔（秒），0 表示禁用自动刷新
 enable_alarm: true      # 启用/禁用警报声音
 show_host: false        # 在服务器名称后显示主机名
 timeout: 3              # RPC2 接口连接超时时间（秒）
 port: 8080             # Web 服务端口

 # Web界面 HTTP Basic 认证配置（可选）
 # 留空则禁用认证
 http_auth:
   username: "admin"    # 登录用户名
   password: "admin123"  # 登录密码

 supervisor_servers:
  server01:
    url: "http://server01.app/RPC2"
    port: "9001"
    username: "yourusername"
    password: "yourpass"
  server02:
    url: "http://server02.app/RPC2"
    port: "9001"

redmine:
  url: "http://redmine.url/path_to_new_issue_url"
  assignee_id: "69"
```

## 运行

### 基本运行

```bash
# Windows
supervisord-monitor.exe

# Linux/Mac
./supervisord-monitor
```

### 指定配置文件

```bash
supervisord-monitor.exe -config /path/to/config.yaml
```

### 指定端口

```bash
supervisord-monitor.exe -port 9000
```

### 组合使用

```bash
 supervisord-monitor.exe -config config.yaml -port 8080
 ```

 默认访问地址：`http://localhost:8080`

 ### HTTP 认证

 如果在 `config.yaml` 中配置了 `http_auth`，访问 Web 界面时会弹出登录对话框：

 - 用户名：配置文件中 `http_auth.username` 的值
 - 密码：配置文件中 `http_auth.password` 的值

 如需禁用认证，在配置文件中将 `username` 和 `password` 留空即可。

## 功能说明

### 进程控制

- **Start**: 启动单个进程
- **Stop**: 停止单个进程
- **Restart**: 重启单个进程（先停止再启动）
- **Start All**: 启动服务器上的所有进程
- **Stop All**: 停止服务器上的所有进程
- **Restart All**: 重启服务器上的所有进程

### 进程状态

- **RUNNING**: 进程正在运行（绿色）
- **STOPPED**: 进程已停止（黑色）
- **STARTING**: 进程正在启动（蓝色）
- **FATAL**: 进程启动失败（红色，触发警报）
- **EXITED**: 进程已退出（红色，触发警报）

### 错误日志

- 当进程有 stderr 日志时，会在进程名称旁显示警告图标
- 点击警告图标可查看错误日志内容
- 进程状态为 FATAL 时会自动触发警报

### 声音警报

- 当检测到 FATAL 状态或 stderr 日志时播放警报
- 点击 "Mute" 按钮可静音（保存到浏览器本地存储）
- 静音状态在刷新页面后保持

### 自动刷新

- 默认每 10 秒自动刷新页面
- 在配置文件中修改 `refresh` 选项可调整间隔
- 设置为 0 可禁用自动刷新

## API 接口

### 获取仪表板数据

```
GET /api/dashboard?mute=1
```

### 进程控制

```
POST /api/start/{server}/{worker}
POST /api/stop/{server}/{worker}
POST /api/restart/{server}/{worker}
POST /api/clear/{server}/{worker}
POST /api/startall/{server}
POST /api/stopall/{server}
POST /api/restartall/{server}
```

## 故障排查

### "Did not receive a '200 OK' response from remote server"

- 检查防火墙和网络连接
- 确保 Supervisord RPC2 接口可访问
- 验证 URL 和端口配置正确

### "HTTP/1.0 401 Unauthorized"

- 用户名或密码错误
- 检查配置文件中的认证信息

### "UNKNOWN_METHOD"

- Supervisord v3+ 需要启用 RPC 接口
- 在配置文件中添加 `[rpcinterface:supervisor]` 部分

### 无法启动可执行文件

- 确保可执行文件有执行权限（Linux/Mac）
- Windows 可能需要允许运行未签名的应用程序

## 开发

### 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 安装 Node 依赖
cd frontend
npm install
```

### 前端开发

```bash
cd frontend
npm run dev
```

访问 `http://localhost:5173`

### 后端开发

```bash
go run main.go
```

访问 `http://localhost:8080`

## 项目结构

```
supervisord-monitor-go/
├── config/
│   └── config.go          # 配置管理
├── handlers/
│   └── handlers.go        # API 处理器
├── services/
│   └── supervisord.go     # Supervisord XML-RPC 客户端
├── webembed/
│   └── embed.go           # 静态文件嵌入
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── ServerCard.vue  # 服务器卡片组件
│   │   ├── App.vue        # 主应用组件
│   │   ├── main.js        # 入口文件
│   │   └── style.css      # 全局样式
│   ├── index.html         # HTML 模板
│   ├── package.json       # Node 依赖
│   └── vite.config.js     # Vite 配置
├── sounds/
│   └── alert.mp3          # 警报声音
├── config.yaml            # 配置文件示例
├── go.mod                 # Go 模块定义
├── main.go                # 主程序入口
├── build.bat              # Windows 构建脚本
└── build.sh               # Linux/Mac 构建脚本
```

## 许可证

MIT License

## 致谢

- 原项目：[mlazarov/supervisord-monitor](https://github.com/mlazarov/supervisord-monitor)
- 基于 CodeIgniter 2.x + Bootstrap 2.x + PHP 重构
