# gRPC - Upload

client에서 Server로 대용량파일을 보내고 Server쪽에서 파일을 저장

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

<br>

<br>

# gRPC, Socket 벤치마킹

gRPC와 Socket의 속도 측정

코드는 [AsyngRPC](async_rpc/) 파일에 있으며, 코드의 내용은 파일의 사이즈를 정하고 전달 횟수를 설정 한 후 송수신까지의 시간을 측정하는 코드입니다.

속도 테스트는 linux 서버에서 진행을 했고 일반 개인 컴퓨터가아닌 회사 서버에서 진행하였습니다.

다만, N개의 Client는 하나의 VM에서 N번의 명령어를 `&`(백그라운드 실행)으로 여러번 실행함을 뜻합니다.

<br>

## gRPC
### Server 1 : 1 Client

- 1:1 1KB를 100만번 보내 총 1GB의 용량 전송 속도

| gRPC | Socket | 차이 |
| ---- | ------ | ---- |
| 17s  | 16s    | 1s   |

<br>

- 1:1 16KB를 100만번 보내 총 16GB의 용량 전송 속도

| gRPC  | Socket | 차이 |
| ----- | ------ | ---- |
| 4m37s | 4m39s  | 2s   |

<br>

- 1:1 32KB를 100만번 보내 총 32GB의 용량 전송 속도

| gRPC   | Socket | 차이 |
| ------ | ------ | ---- |
| 18m34s | 18m35s | 1s   |

<br>



### Server 1 : 12 Client

- 1 : 12 1KB를 100만번 보내 총 1GB의 용량 전송 속도

|      | gRPC  | Socket | 차이  |
| ---- | ----- | ------ | ----- |
| 1    | 1m21s | 1m55s  | 34s   |
| 2    | 1m22s | 2m3s   | 41s   |
| 3    | 1m24s | 2m11s  | 47s   |
| 4    | 1m26s | 2m19s  | 53s   |
| 5    | 1m25s | 2m27s  | 1m2s  |
| 6    | 1m25s | 2m35s  | 1m10s |
| 7    | 1m25s | 2m43s  | 1m18s |
| 8    | 1m48s | 2m51s  | 1m3s  |
| 9    | 1m48s | 2m59s  | 1m11s |
| 10   | 1m47s | 3m8s   | 1m21s |
| 11   | 1m52s | 3m16s  | 1m24s |
| 12   | 1m48s | 3m24s  | 1m36s |

<br>

- 12 : 1 16KB를 100만번 보내 총 16GB의 용량 전송 속도

|      | gRPC   | Socket | 차이   |
| ---- | ------ | ------ | ------ |
| 1    | 19m28s | 30m4s  | 10m36s |
| 2    | 19m26s | 32m35s | 13m9s  |
| 3    | 19m26s | 34m53s | 15m27s |
| 4    | 19m30s | 37m12s | 17m42s |
| 5    | 27m25s | 39m31s | 12m6s  |
| 6    | 27m32s | 41m50s | 14m18s |
| 7    | 27m55s | 44m9s  | 16m14s |
| 8    | 27m55s | 46m28s | 18m32s |
| 9    | 28m41s | 48m46s | 20m5s  |
| 10   | 28m55s | 51m5s  | 22m10s |
| 11   | 29m54s | 53m24s | 23m30s |
| 12   | 30m27s | 55m42s | 25m15s |

<br>

- 12 : 1 32KB를 100만번 보내 총 32GB의 용량 전송 속도

|      | gRPC     | Socket   | 차이     |
| ---- | -------- | -------- | -------- |
| 1    | 1h36m58s | 2h0m48s  | 23m50s   |
| 2    | 1h38m26s | 2h10m4s  | 31m38s   |
| 3    | 1h40m15s | 2h19m21s | 39m6s    |
| 4    | 1h44m10s | 2h28m38s | 44m28s   |
| 5    | 1h44m39s | 2h37m54s | 53m14s   |
| 6    | 1h48m4s  | 2h47m11s | 59m7s    |
| 7    | 1h52m53s | 2h56m27s | 1s3m34s  |
| 8    | 1h55m15s | 3h5m44s  | 1s10m29s |
| 9    | 1h58m46s | 3h15m0s  | 1h16m14s |
| 10   | 1h58m52s | 3h24m17s | 1h25m25s |
| 11   | 1h59m15s | 3h33m33s | 1h34m18s |
| 12   | 1h59m27s | 3h42m50  | 1h43m23s |

