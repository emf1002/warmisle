@echo off
setlocal enabledelayedexpansion

:: ============================================================
:: E2E 测试隔离运行器
:: 用法:
::   e2e\run-test.bat tests\auth.spec.ts
::   e2e\run-test.bat tests\ledger.spec.ts tests\todo.spec.ts
::   e2e\run-test.bat --port 8090 tests\auth.spec.ts
::   e2e\run-test.bat --all
:: ============================================================

set "E2E_DIR=%~dp0"
set "ROOT_DIR=%E2E_DIR%.."
set "BINARY=%ROOT_DIR%\dist\warmisle.exe"
set "HC_PORT=8081"
set "TEST_FILES="

:: Parse arguments
:parse_args
if "%~1"=="" goto :args_done
if "%~1"=="--port" (
    set "HC_PORT=%~2"
    shift
    shift
    goto :parse_args
)
if "%~1"=="--all" (
    set "TEST_FILES=."
    shift
    goto :parse_args
)
if defined TEST_FILES (
    set "TEST_FILES=!TEST_FILES! %~1"
) else (
    set "TEST_FILES=%~1"
)
shift
goto :parse_args
:args_done

if not defined TEST_FILES (
    echo Usage: run-test.bat [--port PORT] ^<test-file^> [test-file2 ...]
    echo        run-test.bat [--port PORT] --all
    exit /b 1
)

:: Auto-increment port if default is in use
:check_port
netstat -an 2>nul | findstr ":%HC_PORT% " | findstr "LISTENING" >nul 2>&1
if !errorlevel! equ 0 (
    set /a HC_PORT+=1
    goto :check_port
)

set "DB_PATH=%E2E_DIR%e2e-data\test-%HC_PORT%.db"
set "BASE_URL=http://localhost:%HC_PORT%"

:: Clean stale DB
if exist "%DB_PATH%" del /f "%DB_PATH%" >nul 2>&1
if exist "%DB_PATH%-wal" del /f "%DB_PATH%-wal" >nul 2>&1
if exist "%DB_PATH%-shm" del /f "%DB_PATH%-shm" >nul 2>&1

:: Build binary if missing
if not exist "%BINARY%" (
    echo [run-test] Binary not found, building...
    pushd "%ROOT_DIR%\frontend"
    call npm run build -- --emptyOutDir
    if !errorlevel! neq 0 ( popd & exit /b 1 )
    popd
    pushd "%ROOT_DIR%\backend"
    go build -o ..\dist\warmisle .
    if !errorlevel! neq 0 ( popd & exit /b 1 )
    popd
)

:: Start server
echo [run-test] Starting server on :%HC_PORT% ...
start "" /b cmd /c ""%BINARY%" >nul 2>&1"
set "SERVER_PID="

:: Wait for server ready
set WAIT_COUNT=0
:wait_ready
curl -s "%BASE_URL%/api/init/check" >nul 2>&1
if !errorlevel! equ 0 goto :server_ready
set /a WAIT_COUNT+=1
if !WAIT_COUNT! geq 75 (
    echo [run-test] ERROR: Server not ready after 15s
    exit /b 1
)
timeout /t 1 /nobreak >nul
goto :wait_ready
:server_ready
echo [run-test] Server ready on :%HC_PORT%

:: Run tests
set HC_TEST_PORT=%HC_PORT%
set HC_TEST_BASE_URL=%BASE_URL%

echo [run-test] Running: npx playwright test %TEST_FILES%
cd /d "%E2E_DIR%"
npx playwright test %TEST_FILES%
set "TEST_RESULT=!errorlevel!"

:: Cleanup - kill server on this port
for /f "tokens=5" %%p in ('netstat -aon 2^>nul ^| findstr ":%HC_PORT% " ^| findstr "LISTENING"') do (
    taskkill /F /PID %%p >nul 2>&1
)

:: Clean DB
if exist "%DB_PATH%" del /f "%DB_PATH%" >nul 2>&1
if exist "%DB_PATH%-wal" del /f "%DB_PATH%-wal" >nul 2>&1
if exist "%DB_PATH%-shm" del /f "%DB_PATH%-shm" >nul 2>&1

echo [run-test] Done (exit code: %TEST_RESULT%)
exit /b %TEST_RESULT%
