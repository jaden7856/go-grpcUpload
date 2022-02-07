# gRPC - Upload

client에서 Server로 대용량파일을 보내고 Server쪽에서 파일을 저장

<br>

<br>

## Usage

Server : start the server( default destination of files is /tmp) :

- windows

```
$ server serve --a <ip:port> -d <destination folder>
Eg server serve -a localhost:8080 -d /home/
```

- linux

```
$ ./server-linux serve --a <ip:port> -d <destination folder>
Eg ./server-linux serve -a localhost:8080 -d /home/
```

<br>

Client : Upload all files in the specified directory to the server :

- windows

```
$ client upload  -a <ip:port> -d <folder containing files to upload>   
Eg  client upload -a localhost:8080 -d /home/
```

- linux

```
$ ./client-linux upload  -a <ip:port> -d <folder containing files to upload>   
Eg  ./client-linux upload -a localhost:8080 -d /home/
```



### 참조

https://github.com/rickslick/grpcUpload

