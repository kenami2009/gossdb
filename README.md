#ssdb协议（同redis使用RESP）

##RESP协议
```
For Simple Strings, the first byte of the reply is "+"
For Errors, the first byte of the reply is "-"
For Integers, the first byte of the reply is ":"
For Bulk Strings, the first byte of the reply is "$"
For Arrays, the first byte of the reply is "*"
In RESP, different parts of the protocol are always terminated with "\r\n" (CRLF).
```
##请求
```
* <参数数量> CR LF
$ <参数 1 的字节数量> CR LF
<参数 1 的数据> CR LF
... 
$ <参数 N 的字节数量> CR LF
<参数 N 的数据> CR LF
```
比如命令：set key value
```go
//字符串表示 set key value
"*3\r\n$3\r\nset\r\n$3key\r\n$5value\r\n"
```

##协议内容
```
//错误
-ERR Unknown Command: cmd
//成功
+OK
//空
+-1
//成功
:0


```

- 
