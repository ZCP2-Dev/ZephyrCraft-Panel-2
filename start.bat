@echo off
echo ========================================
echo    ZephyrCraft Panel 2 启动脚本
echo ========================================
echo.

echo 正在启动后端服务器...
cd Server
start "ZephyrCraft Panel Backend" cmd /k "go run ."
cd ..

echo.
echo 等待后端启动...
timeout /t 3 /nobreak >nul

echo.
echo 正在启动前端应用...
cd Client\Tauri
start "ZephyrCraft Panel Frontend" cmd /k "npm run dev"
cd ..\..

echo.
echo ========================================
echo    启动完成！
echo ========================================
echo.
echo 后端地址: ws://localhost:1145
echo 前端地址: http://localhost:1420
echo.
echo 如果前端没有自动打开，请手动访问:
echo http://localhost:1420
echo.
pause 