@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

set ROOT_DIR=%~dp0
set BACKEND_DIST=%ROOT_DIR%backend\frontend\dist

echo ==============================
echo   Build WarmIsle
echo ==============================

echo [1/5] Clean frontend dist...
if exist "%BACKEND_DIST%" (
    rmdir /s /q "%BACKEND_DIST%"
)
mkdir "%BACKEND_DIST%"

echo [2/5] Check frontend deps...
if not exist "%ROOT_DIR%frontend\node_modules" (
    pushd "%ROOT_DIR%frontend"
    call npm install
    if !errorlevel! neq 0 (
        echo npm install failed! Error: !errorlevel!
        popd
        goto :fail
    )
    popd
)

echo [3/5] Build frontend...
pushd "%ROOT_DIR%frontend"
call npm run build -- --emptyOutDir
if !errorlevel! neq 0 (
    echo Frontend build failed! Error: !errorlevel!
    popd
    goto :fail
)
popd

echo [4/5] Create output dir...
if not exist "%ROOT_DIR%dist" (
    mkdir "%ROOT_DIR%dist"
)

echo [5/5] Build backend...
pushd "%ROOT_DIR%backend"
go build -o "%ROOT_DIR%dist\warmisle.exe" .
if !errorlevel! neq 0 (
    echo Backend build failed! Error: !errorlevel!
    popd
    goto :fail
)
popd

echo ==============================
echo   Build complete
echo   Binary: %ROOT_DIR%dist\warmisle.exe
echo ==============================

explorer "%ROOT_DIR%dist"

:fail
pause
endlocal
