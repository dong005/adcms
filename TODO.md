# ADCMS 开发任务清单

## 第一阶段：项目初始化 ✅

- [x] 1.1.1 创建项目根目录
- [x] 1.1.2 创建 api-backend/ 后端目录
- [x] 1.1.3 克隆 vue-vben-admin 到 api-frontend/
- [x] 1.2.1 执行 go mod init adcms
- [x] 1.2.2 创建后端目录结构 (cmd/internal/pkg)
- [x] 1.2.3 创建 config.yaml 配置文件
- [x] 1.2.4 创建 cmd/main.go 入口文件
- [x] 1.2.5 安装核心依赖 (gin/gorm/redis/jwt/totp/viper/zap)
- [x] 1.3.1 创建配置结构体
- [x] 1.3.2-1.3.6 配置模块完成
- [x] 1.4.1-1.4.5 数据库连接模块完成
- [x] 1.5.1-1.5.4 日志模块完成
- [x] 1.6.1 创建 api-backend/start.sh
- [x] 1.6.2 创建 api-frontend/start.sh
- [x] 1.6.3 创建根目录 start.sh

## 第二阶段：数据库设计 ✅

- [x] 2.1 租户表 (tenants)
- [x] 2.2 用户表 (users) - 含 totp_secret
- [x] 2.3 角色表 (roles)
- [x] 2.4 权限表 (permissions)
- [x] 2.5 用户角色关联表 (user_roles)
- [x] 2.6 角色权限关联表 (role_permissions)
- [x] 2.7 菜单表 (menus)
- [x] 2.8 角色菜单关联表 (role_menus)
- [x] 2.9 分类表 (categories)
- [x] 2.10 文章表 (articles)
- [x] 2.11 标签表 (tags)
- [x] 2.12 文章标签关联表 (article_tags)
- [x] 2.13 媒体表 (media)
- [x] 2.14 系统配置表 (system_configs)
- [x] 2.15 操作日志表 (operation_logs)
- [x] 2.16 登录日志表 (login_logs)
- [x] 2.17 数据库迁移 + 初始化数据

## 第三阶段：认证系统 ✅

- [x] 3.1 密码工具 (bcrypt)
- [x] 3.2 JWT工具 (GenerateToken/ParseToken/RefreshToken)
- [x] 3.3 TOTP工具 (GenerateSecret/GenerateQRCode/ValidateCode)
- [x] 3.4 认证中间件 (JWT验证/黑名单)
- [x] 3.5 租户中间件
- [x] 3.6 CORS中间件
- [x] 3.7 POST /api/auth/login 登录接口
- [x] 3.8 POST /api/auth/verify-totp TOTP验证
- [x] 3.9 POST /api/auth/totp/generate 生成二维码
- [x] 3.9 POST /api/auth/totp/bind 绑定TOTP
- [x] 3.9 POST /api/auth/totp/disable 解绑TOTP
- [x] 3.10 GET /api/auth/user-info 用户信息
- [x] 3.10 PUT /api/auth/user-info 更新信息
- [x] 3.11 PUT /api/auth/password 修改密码
- [x] 3.12 POST /api/auth/logout 登出

## 第四阶段：菜单管理（后端接管）✅

- [x] 4.1 菜单 Repository (CRUD/BuildTree/FindByRoleIDs)
- [x] 4.2 菜单 Handler
- [x] 4.3 POST/PUT/DELETE/GET /api/menus
- [x] 4.3 GET /api/menus/tree
- [x] 4.3 GET /api/menus/user 当前用户菜单
- [x] 4.4 PUT/GET /api/roles/:id/menus 角色菜单

## 第五阶段：用户管理 ✅

- [x] 5.1 用户 Repository
- [x] 5.2 用户 Handler
- [x] 5.3 POST/PUT/DELETE/GET /api/users
- [x] 5.3 PUT /api/users/:id/status
- [x] 5.3 PUT /api/users/:id/reset-password
- [x] 5.3 PUT /api/users/:id/roles

## 第六阶段：角色权限管理 ✅

- [x] 6.1 角色 CRUD
- [x] 6.2 权限管理
- [x] 6.3 PUT/GET /api/roles/:id/permissions

## 第七阶段：租户管理 ✅

- [x] 7.1 租户 CRUD
- [x] 7.2 创建租户时自动创建管理员/角色/菜单

## 第八阶段：CMS 内容管理 ✅

- [x] 8.1 分类管理 CRUD
- [x] 8.2 标签管理 CRUD
- [x] 8.3 文章管理 CRUD + 发布/草稿
- [x] 8.4 媒体管理 上传/删除/列表

## 第九阶段：系统管理 ✅

- [x] 9.1 系统配置 GET/PUT /api/configs
- [x] 9.2 操作日志中间件 + GET /api/logs/operation
- [x] 9.3 登录日志 GET /api/logs/login

## 第十阶段：前端对接

- [x] 10.1.1 修改 .env.development API地址
- [x] 10.1.2 修改端口为 3004
- [x] 10.2 登录页面对接
  - [x] 10.2.1 修改登录 API 请求地址 (/api/auth/login)
  - [x] 10.2.2 添加租户字段到登录表单
  - [x] 10.2.3 实现 TOTP 二维码显示功能
  - [x] 10.2.4 添加 TOTP 验证步骤
  - [x] 10.2.5 处理登录响应和 token 存储
  - [x] 10.2.6 实现自动跳转到首页
- [x] 10.3 菜单动态加载对接
  - [x] 10.3.1 修改菜单获取 API (/api/menus/user)
  - [x] 10.3.2 处理菜单树数据格式转换
  - [x] 10.3.3 实现菜单权限控制
  - [x] 10.3.4 添加路由守卫验证
  - [x] 10.3.5 实现菜单刷新功能
- [x] 10.4 用户中心对接
  - [x] 10.4.1 用户信息 API 对接 (/api/auth/user-info)
  - [x] 10.4.2 修改密码功能对接 (/api/auth/password)
  - [x] 10.4.3 TOTP 绑定/解绑功能对接
  - [x] 10.4.4 个人信息编辑功能
  - [x] 10.4.5 头像上传功能
- [x] 10.5 管理页面对接
  - [x] 10.5.1 用户管理页面 CRUD 对接
  - [x] 10.5.2 角色管理页面 CRUD 对接
  - [x] 10.5.3 菜单管理页面 CRUD 对接
  - [x] 10.5.4 租户管理页面 CRUD 对接
  - [x] 10.5.5 文章管理页面 CRUD 对接
  - [x] 10.5.6 分类标签管理对接
  - [x] 10.5.7 媒体管理对接
  - [x] 10.5.8 系统配置对接
  - [x] 10.5.9 操作日志查看对接
  - [x] 10.5.10 登录日志查看对接
- [x] 10.6 其他功能对接
  - [x] 10.6.1 全局错误处理优化
  - [x] 10.6.2 Loading 状态优化
  - [x] 10.6.3 权限按钮控制
  - [x] 10.6.4 多租户切换功能
  - [x] 10.6.5 退出登录功能优化

## 默认账号

- **用户名**: admin
- **密码**: admin123
- **租户**: default
