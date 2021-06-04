# go-im-system

Go-im-system is a simple im-system, almost 500 lines. It can not be used for production, but it is helpful for go programming learning.

![image](https://user-images.githubusercontent.com/74177355/120818427-324a8980-c585-11eb-81b8-38026a8bab33.png)


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

## TODO

+ [ ] File transfer function


