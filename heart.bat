@echo off
tasklist | find  /c "influxd"
:endFail
exit 0

:end
exit 1