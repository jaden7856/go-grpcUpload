# gRPC - Upload

client에서 Server로 대용량파일을 보내고 Server쪽에서 파일을 저장

<br>

<br>

## Usage

Server : start the server( default destination of files is /tmp) :

- windows

```
$ server-win serve --a <ip:port> -d <destination folder>
Eg server-win serve -a localhost:9191 -d /media/
```

- linux

```
$ ./server-linux serve --a <ip:port> -d <destination folder>
Eg ./server-linux serve -a localhost:9191 -d /media/
```



Client : Upload all files in the specified directory to the server :

- windows

```
$ client-win upload  -a <ip:port> -d <folder containing files to upload>   
Eg  client-win upload -a localhost:9191 -d /home/
```

- linux

```
$ ./client-linux upload  -a <ip:port> -d <folder containing files to upload>   
Eg  ./client-linux upload -a localhost:9191 -d /home/
```



### 참조

https://github.com/rickslick/grpcUpload

