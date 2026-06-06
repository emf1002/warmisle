@echo off
setlocal

set "E2E_DIR=%~dp0"

if not exist "%E2E_DIR%..\dist\warmisle.exe" (
    echo [test.bat] Binary not found, building...
    pushd "%E2E_DIR%..\frontend"
    call npm run build -- --emptyOutDir
    if errorlevel 1 (
        echo [test.bat] Frontend build failed.
        popd
        exit /b 1
    )
    popd
    pushd "%E2E_DIR%..\backend"
    call go build -o ..\dist\warmisle .
    if errorlevel 1 (
        echo [test.bat] Backend build failed.
        popd
        exit /b 1
    )
    popd
)

cd /d "%E2E_DIR%"
npx playwright test %*

endlocal
