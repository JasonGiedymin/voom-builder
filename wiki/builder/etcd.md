# ETCD

## Base etcd config

```
etcdctl set /app/config/supervisor "{\"taskqueue\":\"localhost\"}"
```

## Requires etcd

Project requires ETCD which can be `go get` and requires the etcd client
to be running (for testing).

```bash
go get github.com/coreos/go-etcd/etcd
```

Then:
 - Install etcd (on a non dev machine must use go 1.2) `brew install etcd`
 - Install etcdctl (command line tool) `brew install etcdctl`

 To run:
 ```bash
 $> etcd
 ```

 ## Use `etcdctl`

 For more in-depth usage see the [repo](https://github.com/coreos/etcdctl).

 Setting keys

 ```bash
 etcdctl set /app/config/supervisor "{\"taskqueue\":\"localhost\"}"
 ```