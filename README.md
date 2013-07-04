Batch Download Tools 
===
##Install
```
go get github.com/Centny/BatchDownloader/src/org.cny.bdown/bdown
```

##Usage
```
bdown <URL Pattern> <Folder Pattern> <File Name Pattern> [Batch number or char]
```
##Example
```
bdown "http://www.baidu.com/\1/\2" "\1" "\2.html" a-z:3 1-100:3
bdown "http://www.baidu.com/\1/\2" "\1" "\2.html" A-Z:3 1-100:3
bdown "http://www.baidu.com/\1/\2/\3" "\1" "\3.html" A-Z:3 1-100:3 1-2
bdown "http://www.baidu.com/\1/\3/\2" "\1/\3" "\2.html" A-Z:3 1-100:3 1-2
```