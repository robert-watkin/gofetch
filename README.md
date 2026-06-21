# gofetch

A small Go CLI project that fetches a URL and writes it to stdout. Various options are available to control the output format. Similar to `curl`.

## Installation

```bash
go get github.com/robert-watkin/gofetch
```

## Usage

```bash
gofetch [flags] URL
```

### Flags

- `-header` or `-H`: Add a header to the request. Can be specified multiple times.
- `-json`: Process the response body as JSON.
- `-timeout`: Set the request timeout.
- `-verbose`: Output request and response headers.

### Example

```bash
$ gofetch -json -verbose -header "X-API-Key: 1234567890" https://api.example.com/v1/users
