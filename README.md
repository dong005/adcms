# ADCMS - 广告内容管理系统

基于 Go + Vue 3 的多租户内容管理系统。

## 技术栈

### 后端
- **Go** + **Gin** Web 框架
- **GORM** ORM
- **MySQL** 数据库
- **Redis** 缓存
- **JWT** + **TOTP** 双因素认证

### 前端
- **Vue 3** + **TypeScript**
- **vue-vben-admin** 中后台框架（Monorepo）
- **Ant Design Vue** UI 组件库
- **Vite** 构建工具

## 功能特性

- 多租户隔离
- RBAC 权限管理（用户、角色、部门、菜单）
- 内容管理（文章、分类、标签）
- 媒体管理（图片/视频/文档上传）
- 富文本编辑器（Quill）
- TOTP 谷歌验证器
- 操作日志 / 登录日志
- 通知系统
- 系统配置管理

## 项目结构

```
adcms/
├── api-backend/          # Go 后端
│   ├── cmd/              # 入口
│   ├── internal/         # 业务逻辑
│   │   ├── handler/      # 控制器
│   │   ├── middleware/    # 中间件
│   │   ├── model/        # 数据模型
│   │   ├── repository/   # 数据访问层
│   │   └── router/       # 路由
│   └── pkg/              # 公共包
│       ├── database/     # 数据库
│       ├── redis/        # Redis
│       ├── storage/      # 文件存储
│       └── utils/        # 工具函数
├── api-frontend/         # Vue 前端 (Monorepo)
│   ├── apps/web-antd/    # Ant Design Vue 应用
│   └── packages/         # 共享包
└── start.sh              # 一键启动脚本
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 20+
- pnpm 9+
- MySQL 8.0+
- Redis 6+

### 配置

1. 创建 MySQL 数据库 `adcms`
2. 修改后端配置 `api-backend/config.yaml`

### 启动

```bash
# 一键启动前后端
bash start.sh
```

- 后端 API：http://localhost:8004
- 前端页面：http://localhost:3004
- 默认账号：`admin` / `admin123`

## License

MIT
