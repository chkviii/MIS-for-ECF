# MyPage Backend API

# Ignore this part please

基于 Go + Fiber + MySQL 构建的博客评论系统后端 API。

## 路由说明

### 认证相关路由
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录  
- `GET /api/auth/verify` - 验证JWT token（需要认证）

### 评论相关路由
- `GET /api/comments` - 获取评论列表（支持分页和按文章筛选）
- `POST /api/comments` - 创建评论（需要认证）
- `DELETE /api/comments/:id` - 删除评论（需要认证，只能删除自己的评论）

## 启动方式

### Windows
```bash
start.bat
```

### Linux/Mac
```bash
chmod +x start.sh
./start.sh
```

### 手动启动
```bash
go mod tidy
go run cmd/server/main.go
```

## 环境配置 TODO

创建 `.env` 文件并配置以下参数：
```env
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=mypage
JWT_SECRET=your-super-secret-jwt-key
```

## 数据库

系统会自动创建所需的数据表，TODO

# Main content starts here

## Dirctories
- backend
    - internal
        - config
        - handler
        - middleware
        - repo: database io
            - database.go: table structure and connection
            - ameta.go
            - astatistics.go
            - comment.go
            - user.go
        - service
        - util
    - servver
    - test

