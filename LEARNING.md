# What I've Learned So Far

## Programming

### Concorrency
The `net/http` lib create a new goroutine for each new connection. This means that concorrency is the default behaviour. This is [based on this article](https://eli.thegreenplace.net/2021/life-of-an-http-request-in-a-go-server/).

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