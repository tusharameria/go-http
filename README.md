# go-http

### custom http server in go

## Run the TCP server

```bash
go run cmd/main.go
```

## Connect with the server

```bash
nc 127.0.0.1 8082
```

## Mini Protocol

```
<TYPE> <PAYLOAD>\n
```

```
<TYPE> = ["JOIN", "NAME, "MSG"]
```

### eg

```
JOIN Tushar
```
