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

```bash
<TYPE> <PAYLOAD>\n
```

```bash
<TYPE> = ["JOIN", "NAME, "MSG"]
```

## Example

### Join server

```bash
JOIN Tushar
```

### Send Message

```bash
MGS Hello, this is Tushar
```

### Change Display Name

```bash
NAME Ameria
```
