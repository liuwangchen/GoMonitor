# GoMonitor
EasyMonitor

![Image text](https://raw.githubusercontent.com/liuwangchen/GoMonitor/master/Image/monitor.png)

主页：http://localhost:8080/    会返回个uuid，用于其他api
<br/>cpu接口：ws://localhost:8080/monitorCpu?uuid=fdsfds
<br/>cpu排序接口：http://localhost:8080/cpuSort?uuid=fdsfds&propertyName=Used&ad=asc
<br/>property支持项：Used     ad支持项：asc desc
<br/>
<br/>net接口：ws://localhost:8080/monitorNet?uuid=fdsfds
<br/>net排序接口：http://localhost:8080/netSort?uuid=fdsfds&propertyName=Used&ad=asc
<br/>property支持项：BytesSent,BytesRecv,PacketsSent,PacketsRecv,Errin,Errout,Dropin,Dropout,Fifoin,Fifoout     ad支持项：asc desc
<br/>
<br/>process接口：ws://localhost:8080/monitorProcess?uuid=fdsfds
<br/>process排序接口：http://localhost:8080/processSort?uuid=fdsfds&propertyName=Used&ad=asc
<br/>property支持项：MemoryPercent,Id,CPUPercent,    ad支持项：asc desc