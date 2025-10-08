@echo off
setlocal enabledelayedexpansion

REM ================================
REM CONFIGURATION
REM ================================
set APP_NAME=key
set OUTPUT_DIR=bin

REM Ensure bin directory exists
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

REM ================================
REM TARGET LIST
REM ================================
set TARGETS=windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

echo Building %APP_NAME% for selected targets...
echo.

for %%T in (%TARGETS%) do (
    for /f "tokens=1,2 delims=/" %%A in ("%%T") do (
        set GOOS=%%A
        set GOARCH=%%B
        set EXT=
        if "!GOOS!"=="windows" set EXT=.exe

        echo Building for !GOOS! / !GOARCH! ...
        set OUTFILE=%OUTPUT_DIR%\%APP_NAME%-!GOOS!-!GOARCH!!EXT!
        go build -o "!OUTFILE!" .
        if errorlevel 1 (
            echo ❌ Build failed for !GOOS! / !GOARCH!
        ) else (
            echo ✅ Done: !OUTFILE!
        )
        echo.
    )
)

echo ================================
echo All selected builds complete!
echo Output directory: %OUTPUT_DIR%
echo ================================
endlocal
