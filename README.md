# gofetch

A minimal HTTP client (similar to `curl`). Fetches a URL over GET and streams the response body to stdout.

## Installation

```bash
go install github.com/robert-watkin/gofetch@latest
```

Or build from source:

```bash
git clone https://github.com/robert-watkin/gofetch
cd gofetch
go build .
```

## Usage

```
gofetch [flags] URL
```

### Flags

- `-header value` (repeatable): Add a request header in `Name: value` form
- `-timeout duration`: Request timeout (default `30s`)
- `-json`: Only print the body for `application/json` responses; exit non-zero otherwise
- `-verbose`: Print request and response headers to stderr (body stays on stdout)

## Examples

Basic fetch:

```bash
gofetch https://example.com
gofetch https://httpbin.org/get | head -c 200
```

Repeatable headers:

```bash
gofetch --verbose \
  --header "Accept: application/json" \
  --header "X-Test: hello" \
  https://httpbin.org/headers
```

JSON mode (exits 1 for non-JSON):

```bash
gofetch --json https://httpbin.org/json
gofetch --json https://example.com; echo "exit=$?"
```

Timeout:

```bash
gofetch --timeout 1s https://httpbin.org/delay/3; echo $?
gofetch --timeout 5s https://httpbin.org/delay/1; echo $?
```

## Exit codes

- `0` — Success (2xx response and body written)
- `1` — HTTP error (4xx/5xx) or `--json` validation failure
- `2` — Network / timeout / usage error

## Development

```bash
go build .
go test -v -race ./...
```