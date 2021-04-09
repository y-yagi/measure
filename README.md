# measure

measure is a simple tool for measuring file size. Also supports http protocol, you can measure the size before downloading the file.

```
$ measure testdata/test.zip
testdata/test.zip: 8.5K
$ measure https://github.com/y-yagi/measure/archive/master.zip
Redirectd to https://codeload.github.com/y-yagi/measure/zip/master
https://github.com/y-yagi/measure/archive/master.zip: 9.4K
```

## Installation

```
go get -u github.com/y-yagi/measure
```
