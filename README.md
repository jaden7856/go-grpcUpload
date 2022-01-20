# gRPC - Upload

client에서 Server로 대용량파일을 보내고 Server쪽에서 파일을 저장

<br>

<br>

## Usage

Server : start the server( default destination of files is /tmp) :

```
$./grpcUploadServer serve --a <ip:port> -d <destination folder>
Eg ./UploadClient serve -a localhost:9191 -d /media/
```



Client : Upload all files in the specified directory to the server :

```
$ ./UploadClient upload  -a <ip:port> -d <folder containing files to upload>   
Eg  ./UploadClient upload -a localhost:9191 -d /home/
```



### 참조

https://github.com/rickslick/grpcUpload

