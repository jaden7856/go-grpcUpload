# Error 정리

### Network

- Client에서 Server로 데이터를 보낼때 Server에서 연결을 끊을 때

  - `failed by the remote host.erver: read tcp 127.0.0.1:55017->127.0.0.1:8080: wsarecv: An existing connection was forcibly close`

    - Client 자체 에러

    <br>

  - `EOF: failed to send`

    - 파일이 정상적으로 보내지 못했다는 커스텀 에러

<br>

<br>

### Files

- 클라이언트나 서버에서 정상적으로 파일을 송수신을 못했을 때
  - `failed to receive upstream status response`
    - client - `stream.CloseAndRecv()`, server - `stream.SendAndClose()`  각각 이 메서드들을 통해서 정상적으로 완료가 되었는지 아닌지 메세지를 전달

<br>

- 클라이언트에서 서버로 보낼때 오래걸리면 timeout 발생
  - `Error uploading file`
    - 보내는 도중 보내지 못한 파일 이름과 함께 출력

<br>

- 서버에서 파일을 받지 못할때
  - `failed while reading from stream`

<br>

- 서버에서 파일에 data들을 쓸 때 오류가 발생
  - `Unable to write chunk of filename`
    - 전송 도중 오류가 발생하면 서버에서 데이터를 받아오는 것을 실패하여 발생하는 오류

