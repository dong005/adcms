# 构建阶段
FROM node:20-alpine AS builder

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web-antd/package.json apps/web-antd/
COPY packages/ packages/
COPY internal/ internal/

RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm --filter @vben/web-antd build

# 运行阶段 - 使用 Nginx 托管静态文件
FROM nginx:alpine

COPY --from=builder /app/apps/web-antd/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
