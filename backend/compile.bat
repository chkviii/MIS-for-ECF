@echo off
echo Building for ARMv7l...

REM Update Go modules
go mod tidy

REM Set environment variables for ARMv7l
set GOOS=linux
set GOARCH=arm
set GOARM=7
set CGO_ENABLED=0

REM Building...
go build -a -ldflags "-s -w" -o mypage-backend-armv7l server/main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Output: mypage-backend-armv7l
    echo File size:
    dir mypage-backend-armv7l
) else (
    echo Build failed!
    exit /b 1
)

REM 重置环境变量
set GOOS=
set GOARCH=
set GOARM=
set CGO_ENABLED=