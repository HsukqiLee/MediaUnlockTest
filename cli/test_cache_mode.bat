@echo off
echo 测试缓存模式 (串行地区，启用缓存)...
echo.
.\main.exe -cache -m 4 -conc 10
echo.
echo ========================================
echo.
echo 测试非缓存模式 (并行所有测试，不启用缓存)...
echo.
.\main.exe -m 4 -conc 10
pause
