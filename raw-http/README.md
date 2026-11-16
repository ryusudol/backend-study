# Stage 1: Raw HTTP Server

**Rule:** don't use `net/http` in the server. Use only `net`, `bufio`, `io`, `strings`, `strconv` and other low-level packages.
The goal is to see that HTTP is just text sent over a TCP connection.

## Part A: TCP echo server

Listen on `:8080`. For each connection, read lines and write each one back.

```sh
go run .
nc localhost 8080      # type lines and see them echoed back
```

**Done when:**
- [ ] Two `nc` sessions open at once both work (one goroutine per connection)
- [ ] When a client disconnects, the server keeps running and doesn't leak the connection (`defer conn.Close()`)

## Part B: HTTP/1.1 server

Replace the echo handler with one that speaks HTTP.

### What to parse
```
POST /echo HTTP/1.1\r\n        <- request line: METHOD SP TARGET SP VERSION
Host: localhost:8080\r\n       <- headers: Name: value
Content-Length: 5\r\n
\r\n                           <- blank line = end of headers
hello                          <- body: exactly Content-Length bytes
```

### What to send back
```
HTTP/1.1 200 OK\r\n
Content-Type: text/plain\r\n
Content-Length: 5\r\n
Connection: close\r\n
\r\n
hello
```

### Routes
| Request | Response |
|---|---|
| `GET /` | `200` with body `hello` |
| `GET /echo/{text}` | `200` with body `{text}` |
| `POST /echo` | `200`, echoes the request body |
| `GET /headers` | `200`, lists the request headers, one per line |
| unknown path | `404` |
| known path, wrong method | `405` with an `Allow` header |
| malformed request | `400`, and the server stays up |

Close the connection after each response.

### Test it
```sh
curl -v localhost:8080/
curl -v localhost:8080/echo/abc
curl -v -d 'hello' localhost:8080/echo
curl -v -X DELETE localhost:8080/echo
printf 'garbage\r\n\r\n' | nc localhost 8080
```
Also open `http://localhost:8080/` in a browser.

**Done when:**
- [ ] Every row in the routes table behaves as described
- [ ] Header names are matched case-insensitively (`content-length` works the same as `Content-Length`)
- [ ] Malformed input (no blank line, bad request line, non-numeric `Content-Length`) returns `400` and never crashes the server
- [ ] You can explain each line of the `curl -v` output

## Stretch
- [ ] **Keep-alive:** handle several requests on one connection (HTTP/1.1 does this by default)
- [ ] **Read timeout:** a client that connects and sends nothing gets dropped after N seconds (`conn.SetReadDeadline`). Look up the "slowloris" attack.
- [ ] **Size limits:** reject oversized headers (`431`) and bodies (`413`)
- [ ] Write `_test.go` tests that dial your server with `net.Dial` and send raw requests

## Hints (read only when stuck)
- `net.Listen`, `ln.Accept`, `go handle(conn)`
- `bufio.NewReader(conn)`, then `ReadString('\n')`. Remember to trim the `\r\n`.
- `strings.Cut`, `strings.Fields`, `strings.ToLower`, `strconv.Atoi`
- `io.ReadFull(reader, buf)` reads exactly `Content-Length` bytes. Why not just call `Read`?
- Read the body from the **same** `bufio.Reader` you used for the headers. Why?

## Afterwards: questions to be able to answer
1. Why does HTTP need `Content-Length` (or chunked encoding)? What would happen without it?
2. What's the difference between a TCP connection and an HTTP request?
3. What does a goroutine per connection cost? What happens with 10,000 clients?
4. What does `net/http` give you that you just wrote by hand?
