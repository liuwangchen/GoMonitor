# GoMonitor
EasyMonitor

![Image text](https://raw.githubusercontent.com/liuwangchen/GoMonitor/master/Image/monitor.png)

主页：http://localhost:8080/
<br/>cpu接口：ws://localhost:8080/monitorCpu
<br/>cpu排序接口：http://localhost:8080/cpuSort?propertyName=Used&ad=asc
<br/>property支持项：Used     ad支持项：asc desc
<br/>
<br/>net接口：ws://localhost:8080/monitorNet
<br/>net排序接口：http://localhost:8080/netSort?propertyName=Used&ad=asc
<br/>property支持项：BytesSent,BytesRecv,PacketsSent,PacketsRecv,Errin,Errout,Dropin,Dropout,Fifoin,Fifoout     ad支持项：asc desc
<br/>
<br/>process接口：ws://localhost:8080/monitorProcess
<br/>process排序接口：http://localhost:8080/processSort?propertyName=Used&ad=asc
<br/>property支持项：MemoryPercent,Id,CPUPercent,    ad支持项：asc desc