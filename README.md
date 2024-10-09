# Caching Proxy

Caching Proxy is a program designed to proxy HTTP requests to a specified address and cache the results. The program automatically clears the cache at a specified interval and also supports manual cache clearing via a specific API endpoint.

Project idea from <https://roadmap.sh/projects/caching-server>
## Features

- Proxy requests to another server.
- Cache responses to speed up subsequent requests.
- Automatic and manual cache clearing.
- `x-cache` header in the response to indicate whether the result was retrieved from the cache.

## Installation and Usage

### Step 1: Build the executable

Before running the program, you need to build the executable using the `make build` command:

```bash
make build
```

This will compile the program and create an executable file named `caching-proxy`.

### Step 2: Running the program

Once the program is built, you can start the server using the following command:

```bash
./caching-proxy -port=<number> -origin=<url> -auto-clear-cache=<number>
```

Where:
- `-port` is the port the server will listen on (default: 8080).
- `-origin` is the target address where requests will be proxied (default: "*").
- `-auto-clear-cache` is the interval (in minutes) after which the cache will be cleared (default: 30).

### Example

```bash
./caching-proxy -port=8081 -origin="https://api.github.com" -auto-clear-cache=10
```

In this example:
- The server will listen on port `8081`.
- Requests will be proxied to `https://api.github.com`.
- The cache will automatically be cleared every `10` minutes.

### Step 3: Making requests

Assuming the program is running on `localhost:8081`, all requests to `localhost:8081/proxy/*` will be proxied to the specified origin.

For example, making a `GET` request to:

```bash
http://localhost:8081/proxy/users/PavelNikoltsev
```

This will forward the request to:

```bash
https://api.github.com/users/PavelNikoltsev
```

### Cache Behavior

Responses from proxied requests will be cached. On subsequent requests, the cached response will be used.

You can check whether a response was served from the cache by inspecting the `x-cache` header:
- `HIT` means the response was retrieved from the cache.
- `MISS` means the response was fetched directly from the origin server.

### Manual Cache Clearing

To manually clear the cache, send a `POST` request to the following endpoint:

```
POST http://<your-server-address>/api/clear-cache
```

For example, if the server is running on `localhost:8081`:

```bash
POST http://localhost:8081/api/clear-cache
```

This will immediately clear all cached data.

