# What I've Learned So Far

## Programming

### Concorrency
The `net/http` lib create a new goroutine for each new connection. This means that concorrency is the default behaviour. This is [based on this article](https://eli.thegreenplace.net/2021/life-of-an-http-request-in-a-go-server/).


### Context

I see that the application has having timeout, this seems because the database takes long to respond. Because of this, I've impelmented context to timeout the requests.

Example of logs I'm seeing in the stdout:
```log
time=2024-09-03T19:03:35.485Z level=INFO msg="error getting person" error="dial tcp 192.168.112.2:5432: connect: connection timed out"

2024/09/03 19:02:48 http: superfluous response.WriteHeader call from go-rinha-de-backend-2023/internal/handler.(*PersonHandler).SearchPersons (handler.go:135)

2024/09/03 19:02:46 http: superfluous response.WriteHeader call from go-rinha-de-backend-2023/internal/handler.(*PersonHandler).GetPersonById (handler.go:103)
```

I believe the context will help with this. I might lower the time, by firstly it is 10 seconds.

## Database

### Postgres

#### Configuration File
This file can be retrieve in SQL QUERY `SHOW config_file`. The standard file can be then copied to my local machine and modified on demand.

Even when doing that, when I pass the file using the configuration below:
```yaml
 postgres-db:
    ...
    volumes:
    - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    - ./postgresql.conf:/etc/postgresql.conf
    command: ["postgres", "-c", "config_file=/etc/postgresql.conf"]
```
But doing the below works just fine:
```yaml
 postgres-db:
    ...
    volumes:
    - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    command: ["postgres", "-c", "max_connections=150"]
```

Maybe I might be doing something wrong. I can dig more later.

## Nginx

### Connection Pool
Nginx uses worker processes to handle incoming connections. Each worker can handle many connections simultaneously using an asynchronous, non-blocking I/O model. For upstream connections (connections to backend servers), Nginx maintains a pool of connections to these servers, which it reuses to handle multiple client requests efficiently.

The **keepalive** directive in the upstream block sets the maximum number of idle keepalive connections to upstream servers. This helps in reusing connections for multiple requests, reducing connection setup overhead.