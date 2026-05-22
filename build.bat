@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

set ROOT_DIR=%~dp0
set BACKEND_DIST=%ROOT_DIR%backend\frontend\dist

echo ==============================
echo   Build WarmIsle
echo ==============================

echo [1/4] Clean frontend dist...
if exist "%BACKEND_DIST%" (
    rmdir /s /q "%BACKEND_DIST%"
)
mkdir "%BACKEND_DIST%"

echo [2/4] Check frontend deps...
if not exist "%ROOT_DIR%frontend\node_modules" (
    pushd "%ROOT_DIR%frontend"
    call npm install
    popd
)

echo [3/4] Build frontend...
pushd "%ROOT_DIR%frontend"
call npm run build -- --emptyOutDir
if %errorlevel% neq 0 (
    echo Frontend build failed! Error: %errorlevel%
    popd
    exit /b %errorlevel%
)
popd

echo [4/5] Create output dir...
if not exist "%ROOT_DIR%dist" (
    mkdir "%ROOT_DIR%dist"
)

echo [5/5] Build backend...
pushd "%ROOT_DIR%backend"
go build -o "%ROOT_DIR%dist\warmisle.exe" .
if %errorlevel% neq 0 (
    echo Backend build failed! Error: %errorlevel%
    popd
    exit /b %errorlevel%
)
popd

echo ==============================
echo   Build complete
echo   Binary: %ROOT_DIR%dist\warmisle.exe
echo ==============================

explorer "%ROOT_DIR%dist"

pause
endlocal
