@echo off

set OUT_DIR=%~dp0..\bin
set CMD_DIR=%~dp0..\cmd

echo Build into %OUT_DIR%
echo Sources from %CMD_DIR%

go build -o %OUT_DIR%\spa-server.exe %CMD_DIR%\spa_server\main.go
go build -o %OUT_DIR%\builder.exe %CMD_DIR%\builder\main.go