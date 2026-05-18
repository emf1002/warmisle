#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIST="$ROOT_DIR/backend/frontend/dist"

echo "=== 构建暖屿 ==="

# 1. 清空前端历史构建产物
echo ">>> 清空历史构建产物..."
rm -rf "$BACKEND_DIST"
mkdir -p "$BACKEND_DIST"

# 2. 安装前端依赖（如需）
if [ ! -d "$ROOT_DIR/frontend/node_modules" ]; then
  echo ">>> 安装前端依赖..."
  cd "$ROOT_DIR/frontend" && npm install
fi

# 3. 构建前端
echo ">>> 构建前端..."
cd "$ROOT_DIR/frontend" && npm run build -- --emptyOutDir

# 4. 创建输出目录
mkdir -p "$ROOT_DIR/dist"

# 5. 构建 Go 后端二进制
echo ">>> 构建后端..."
cd "$ROOT_DIR/backend" && go build -o "$ROOT_DIR/dist/warmisle" .

echo "=== 构建完成 ==="
echo "二进制文件: $ROOT_DIR/dist/warmisle"
