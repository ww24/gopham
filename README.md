gopham
======

Push message Manager

It is like Google Cloud Messaging.

How to Run
----------
```
go get
go run gopham.go
```

Open (localhost:3000)[http://localhost:3000/] on your browser.

Usage
-----

### WS /subscribe
WebSocket entry point `/subscribe`.

### POST /
Request body (JSON)

```
{
  "channel": "test",
  "ttl": 0,
  "data": {"message": "json"}
}
```

#### channel [required]
チャンネル識別用文字列

現状、明示的に channel を指定しなくても全てのクライアントに broadcast されます。

#### ttl [option]
⚠まだ実装されていません。

* `ttl = 0` の時は揮発性メッセージ
* `ttl > 0` の時は `ttl` の大きさだけキューに留まります

例えば、 `ttl = 1` のメッセージは message list の末尾に追加され、末尾から数えて2以上深くなった場合に削除されます。
重要度が高く消えては困るメッセージの ttl は大きく設定します。

message list はクライアントが接続時に受け取ります。

#### data [required]
任意の JSON 形式

### Example
```
curl 'http://localhost:3000/' --verbose --request POST --header 'Content-Type: application/json' \
--data-binary '{"channel":"test", "data":{"message":"json"}}'
```
