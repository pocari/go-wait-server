# go-wait-server

任意の時間まってレスポンス返す

# 使い方

## ビルド

```
$ go get https://github.com/pocari/go-wait-server
$ go-wait-server
```

## 別端末で

```
curl localhost:8080/wait?time=5
root@51421ff38733:/go# curl localhost:8080/wait?time=5
{"status":200,"message":"success"}
```

time=X で指定した秒数待ってレスポンスが返ってくる(0 <= X <= 30)
