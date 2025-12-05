@echo off
echo Building for linux amd64...

REM Update Go modules
go mod tidy

@REM REM Set environment variables for ARMv7l
@REM set GOOS=linux
@REM set GOARCH=arm
@REM set GOARM=7
@REM set CGO_ENABLED=0

REM Set environment variables for linux x86_64
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=1


REM Building...
go build -a -ldflags "-s -w" -o erp-backend-amd64 ../server/main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Output: erp-backend-amd64
    dir erp-backend-amd64
) else (
    echo Build failed!
    exit /b 1
)

REM 重置环境变量
set GOOS=
set GOARCH=
set GOARM=
set CGO_ENABLED=