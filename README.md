# ByteDanceDemo

![Go Version](https://img.shields.io/badge/Go%20Version-1.20+-00ADD8FF?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-00ADD8FF?style=for-the-badge)
![Build Status](https://img.shields.io/badge/Build-passing-brightgreen?style=for-the-badge)
![Coverage](https://img.shields.io/badge/Coverage-75%25-orange?style=for-the-badge)

ByteDanceDemo 是一个简化版的抖音（Douyin）后端服务，使用 Go (Golang) 构建，提供视频分享、社交互动和实时消息功能。

## ✨ 特性

- 🚀 **用户认证与授权**
  - JWT 令牌认证
  - Casbin RBAC 权限控制
  - bcrypt 密码加密
  
- 📹 **视频分享功能**
  - 视频上传和发布
  - 视频列表和分页
  - 视频封面生成
  
- 💬 **社交互动**
  - 评论系统
  - 点赞/收藏功能
  - 用户关注/粉丝
  - 好友关系
  
- 💬 **实时消息**
  - 私信发送和接收
  - 消息列表查询
  - 未读消息统计
  
- 🔒 **安全防护**
  - SQL 注入防护
  - XSS 攻击防护
  - CSRF 保护
  - 速率限制

## 🚀 快速入门

### 前置条件

- Go 1.20 或更高版本
- MySQL 8.0 或更高版本
- Redis 6.0 或更高版本
- RabbitMQ 3.9 或更高版本

### 安装步骤

1. 克隆项目仓库

```bash
git clone https://github.com/yourusername/ByteDanceDemo.git
cd ByteDanceDemo
```

2. 配置环境变量（国内用户）

```bash
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
```

3. 复制配置文件

```bash
cp config/settings.yml.template config/settings.yml
```

4. 安装依赖

```bash
go mod download
go mod tidy
```

5. 创建数据库

```bash
mysql -u root -p < config/init.sql
```

6. 运行迁移

```bash
make run-migrate
```

7. 启动服务

```bash
make run-api
```

8. 验证安装

```bash
curl http://localhost:8080/douyin/feed/
```

## 📚 文档

- [架构文档](docs/ARCHITECTURE.md) - 系统架构和设计模式
- [配置参考](docs/CONFIGURATION.md) - 完整的配置选项
- [部署指南](docs/DEPLOYMENT.md) - 开发和生产部署
- [贡献指南](docs/CONTRIBUTING.md) - 开发工作流程
- [安全实践](docs/SECURITY.md) - 安全最佳实践
- [性能优化](docs/PERFORMANCE.md) - 性能调优指南
- [故障排除](docs/TROUBLESHOOTING.md) - 常见问题解决
- [API 文档](docs/api/openapi-spec.yaml) - OpenAPI 规范
- [变更日志](CHANGELOG.md) - 版本更新历史
- [常见问题](docs/FAQ.md) - FAQ 常见问题

## 🛠️ 技术栈

| 类别 | 技术 | 用途 |
|-------|------|------|
| Web 框架 | [Gin](https://gin-gonic.com/) | RESTful API |
| 数据库 | [MySQL](https://www.mysql.com/) | 数据持久化 |
| 缓存 | [Redis](https://redis.io/) | 会话和缓存 |
| 消息队列 | [RabbitMQ](https://www.rabbitmq.com/) | 异步任务处理 |
| ORM | [GORM](https://gorm.io/) | 数据库操作 |
| 认证 | [JWT](https://github.com/golang-jwt/jwt) | 令牌认证 |
| 授权 | [Casbin](https://casbin.org/) | 权限控制 |
| 日志 | [Zap](https://github.com/uber-go/zap) | 结构化日志 |
| 配置 | [Viper](https://github.com/spf13/viper) | 配置管理 |

## 📁 项目结构

```
ByteDanceDemo/
├── cmd/                    # 命令行接口
├── config/                 # 配置文件
├── controller/              # HTTP 请求处理
├── service/                 # 业务逻辑层
├── repository/              # 数据访问层
├── model/                  # 数据模型
├── middleware/              # 中间件
├── router/                 # 路由配置
├── database/               # 数据库连接
├── utils/                  # 工具函数
├── test/                   # 测试文件
├── public/                 # 静态资源
├── docs/                   # 文档
├── scripts/                # 脚本工具
├── main.go                # 应用入口
├── go.mod                 # 模块定义
└── Makefile                # 构自动动
```

## 🔌 API 端点

### 基础 API

| 端点 | 方法 | 描述 |
|-------|------|------|
| `/douyin/feed/` | GET | 获取视频流 |
| `/douyin/user/` | GET | 获取用户信息 |
| `/douyin/user/register/` | POST | 用户注册 |
| `/douyin/user/login/` | POST | 用户登录 |
| `/douyin/publish/action/` | POST | 发布视频 |
| `/douyin/publish/list/` | GET | 获取用户发布的视频 |

### 扩展 API

| 端点 | 方法 | 描述 |
|-------|------|------|
| `/douyin/favorite/action/` | POST | 点赞/取消点赞 |
| `/douyin/favorite/list/` | GET | 获取点赞列表 |
| `/douyin/comment/action/` | POST | 发布评论 |
| `/douyin/comment/list/` | GET | 获取评论列表 |
| `/douyin/relation/action/` | POST | 关注/取消关注 |
| `/douyin/relation/follow/list/` | GET | 获取关注列表 |
| `/douyin/relation/follower/list/` | GET | 获取粉丝列表 |
| `/douyin/relation/friend/list/` | GET | 获取好友列表 |
| `/douyin/message/chat/` | GET | 获取聊天消息 |
| `/douyin/message/action/` | POST | 发送消息 |

更多 API 详情请查看 [OpenAPI 规范](docs/api/openapi-spec.yaml)。

## 🧪 测试

项目包含完整的测试套件，包括单元测试、集成测试和性能基准测试。

### 运行所有测试

```bash
make test
```

### 运行特定类型测试

```bash
make test-unit          # 单元测试
make test-integration   # 集成测试
make test-benchmark     # 性能基准测试
make test-coverage      # 生成覆盖率报告
```

### 查看测试文档

- [测试文档](TEST_DOCUMENTATION.md) - 详细测试指南
- [测试目录 README](test/README.md) - 测试快速开始

## 🐳 Docker 支持

项目支持 Docker 容器化部署。

### 使用 Docker Compose

```bash
docker-compose up -d
```

### 构建镜像

```bash
docker build -t bytedancedemo:latest .
```

更多 Docker 部署详情请查看 [部署指南](docs/DEPLOYMENT.md)。

## 🔧 开发指南

### 开发环境设置

```bash
# 安装开发依赖
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 启动开发服务
make run-api

# 运行测试
make dev
```

### 代码规范

- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 遵循 Go 社区代码规范

### 提交流程

1. 创建功能分支
2. 编写代码和测试
3. 运行测试和代码检查
4. 提交 Pull Request
5. 代码审查和合并

更多详情请查看 [贡献指南文档](docs/CONTRIBUTING.md)。

## 📦 构建和运行

### 使用 Makefile

```bash
make build           # 构建项目
make run-api         # 启动 API 服务
make run-parallel    # 并行启动迁移和 API 服务
make ci             # 完整 CI 流程工作流
```

### 直接运行

```bash
go build -o bin/simple-demo
./bin/simple-demo server -c config/settings.yml -m debug
```

## 🔍 安全

项目实施了多层安全防护：

- ✅ SQL 注入防护（参数化查询）
- ✅ XSS 攻击防护（输入清理）
- ✅ CSRF 保护（令牌验证）
- ✅ 速率限制（Redis 支持）
- ✅ JWT 认证（令牌过期）
- ✅ 密码加密（bcrypt）
- ✅ 权限控制（Casbin RBAC）

更多安全细节请查看 [安全文档](docs/SECURITY.md)。

## ⚡ 性能优化

项目包含多项性能优化：

- 🔹 数据库连接池
- 🔹 Redis 缓存策略
- 🔹 异步消息处理（RabbitMQ）
- 🔹 批量查询优化
- 🔹 并发处理优化

更多性能细节请查看 [性能文档](docs/PERFORMANCE.md)。

## 📊 监控和日志

### 查看日志

```bash
# 应用日志
tail -f logs/bytedancedemo.log

# 错误日志
tail -f logs/error.log
```

### 健康检查

```bash
curl http://localhost:8080/health
```

### 使用诊断脚本

```bash
./scripts/health-check.sh
./scripts/quick-check.sh
```

## 🤝 贡献

欢迎贡献！请遵循以下步骤：

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

更多详情请查看 [贡献指南文档](docs/CONTRIBUTING.md)。

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📞 联系方式

- 项目主页: [GitHub](https://github.com/yourusername/ByteDanceDemo)
- 问题反馈: [Issues](https://github.com/yourusername/ByteDanceDemo/issues)
- 文档: [docs/](docs/)

## 🙏 致谢

感谢所有为 ByteDanceDemo 做出贡献的开发者！

---

**最后更新**: 2026-04-21
