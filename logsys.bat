@echo off

@REM echo zookeeper service start
@REM start /min "zookeeper" "D:\zookeeper\apache-zookeeper-3.6.1-bin\bin\zkServer.cmd" 

echo etcd service start
start /min "etcd" "D:\etcd-v3.2.32-windows-amd64\etcd.exe"

echo influxdb service start 
start /min "influxdb" "D:\influxdb-1.8.4_windows_amd64\influxdb-1.8.4-1\influxd.exe"

echo grafana service start 
start /min "grafana" "D:\grafana\grafana\bin\grafana-server.exe"

echo elasticsearch service start 
start /min "elasticsearch" "D:\Elasticsearch\elasticsearch-7.11.2\bin\elasticsearch.bat"

echo kibana service start 
start /min "kibana" "D:\kibana\kibana-7.11.2-windows-x86_64\bin\kibana.bat"

@REM set SLEEP=ping 127.0.0.1 
@REM %SLEEP% 3
@REM set ENV_HOME="D:\kafka\kafka_2.12-2.7.0"

@REM D:
@REM color 0a
@REM cd %ENV_HOME%
@REM echo kafka service start
@REM .\bin\windows\kafka-server-start.bat config\server.properties
pause