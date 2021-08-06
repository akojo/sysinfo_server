# Sysinfo Server

A simple example server providing an HTTP API for `systemd-analyze`.

`sysinfo_server` reads startup kernel and user startup times provided by
`systemd-analyze` and returns them either as plaintext or formatted as JSON,
depending on request `Content-Type` parameter.

## Usage

Start the server:

```sh
$ go run .
```

The server is now running at port 8080. You can change the port number by
defining an alternative port in `PORT` environment variable.

Request startup times:

```sh
curl http://localhost:8080/duration
```

Request startup times as JSON:

```sh
curl -H 'Content-Type: application/json'  http://localhost:8080/duration
```
