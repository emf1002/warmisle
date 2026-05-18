# 构建暖屿 (Windows PowerShell)
$ErrorActionPreference = "Stop"

$ROOT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$BACKEND_DIST = "$ROOT_DIR\backend\frontend\dist"

Write-Host "=== 构建暖屿 ===" -ForegroundColor Cyan

# 1. 清空前端历史构建产物
Write-Host ">>> 清空历史构建产物..." -ForegroundColor Yellow
if (Test-Path $BACKEND_DIST) {
    Remove-Item -Recurse -Force $BACKEND_DIST
}
New-Item -ItemType Directory -Force -Path $BACKEND_DIST | Out-Null

# 2. 安装前端依赖（如需）
if (-not (Test-Path "$ROOT_DIR\frontend\node_modules")) {
    Write-Host ">>> 安装前端依赖..." -ForegroundColor Yellow
    Push-Location "$ROOT_DIR\frontend"
    npm install
    Pop-Location
}

# 3. 构建前端
Write-Host ">>> 构建前端..." -ForegroundColor Yellow
Push-Location "$ROOT_DIR\frontend"
npm run build -- --emptyOutDir
Pop-Location

# 4. 创建输出目录
if (-not (Test-Path "$ROOT_DIR\dist")) {
    New-Item -ItemType Directory -Path "$ROOT_DIR\dist" | Out-Null
}

# 5. 构建 Go 后端二进制
Write-Host ">>> 构建后端..." -ForegroundColor Yellow
Push-Location "$ROOT_DIR\backend"
go build -o "$ROOT_DIR\dist\warmisle.exe" .
Pop-Location

Write-Host "=== 构建完成 ===" -ForegroundColor Green
Write-Host "二进制文件: $ROOT_DIR\dist\warmisle.exe" -ForegroundColor Green
