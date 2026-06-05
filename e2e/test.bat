@echo off
setlocal

set "E2E_DIR=%~dp0"

if not exist "%E2E_DIR%..\dist\warmisle.exe" (
    echo [test.bat] Binary not found, building...
    pushd "%E2E_DIR%.."
    make build
    if errorlevel 1 (
        echo [test.bat] Build failed.
        popd
        exit /b 1
    )
    popd
)

cd /d "%E2E_DIR%"
npx playwright test %*

endlocal
