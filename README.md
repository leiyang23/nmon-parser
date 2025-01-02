## 项目
解析 nmon 导出数据，以图表的形式进行展示，同时支持导出 Excel 文件。

## 使用
### 源码编译
将代码下载到本地后，执行 `go build -o nmon-parser main.go`，然后运行 `nmon-parser web --port=8081`，在浏览器打开地址：http://localhost:2233  

### 或则直接下载编译好的可执行文件：[release](https://gitee.com/code_butcher/nmon-parser/releases/tag/1.0)
Linux 系统：`nmon-parser web --port=8081`  
Windows 系统：`nmon-parser.exe web --port=8081`
