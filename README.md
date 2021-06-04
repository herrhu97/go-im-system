# go-im-system

Go-im-system is a simple im-system, almost 500 lines. It can not be used for production, but it is helpful for go programming learning.

## Feature

* Support change user name
* Support both public chat and private chat
* Timeout automaticly delete user

## Usage

### Run Server

 ~~~bash
cd go-im-system
go build -o server server.go main.go user.go
./server
 ~~~

### Client

~~~
go build -o client
./client
~~~



