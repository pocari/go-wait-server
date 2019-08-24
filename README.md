# go-wait-server

任意の時間まってレスポンス返す

# 使い方

## ビルドしてサーバー起動

```
$ go get github.com/pocari/go-wait-server
$ go-wait-server
```

## 別端末で

```
curl localhost:8080/wait?time=5
root@51421ff38733:/go# curl localhost:8080/wait?time=5
{"status":200,"message":"success"}
```

time=X で指定した秒数待ってレスポンスが返ってくる(0 <= X <= 30)

待っている間にクライアント側で、Ctr-C とかで curl を止めたら、サーバ側の wait 処理も（まだ残り時間あってもキャンセルする
