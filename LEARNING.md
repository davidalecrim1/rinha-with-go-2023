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

I believe the context will help with this. I've tried 5 and 10 seconds. With 5 was not working well and the inserted people was 50% less.

### CPU and Memory Limits

I've removed them to perform some tests. Seems I can reach about 21000 requests. I'll use this to learn and improve my code. [file:///Users/davidalecrim/Documents/Code/rinha-de-backend-2023/rinha-de-backend-2023-q3/stress-test/user-files/results/rinhabackendsimulation-20240904002934456/index.html](check the results).

### Go References

I've been just practicing without lookling at another Go code to not be tempted to copy or test the concept I've haven't seen the need for. I've listed all the rinha versions in Go to explore later in the future:
- | 2 | h4nkb31f0ng | 44270 | 15690 | [README](./stress-test/rinha-de-backend-2023-q3/participantes/h4nkb31f0ng/README.md) |
- | 6 | isadora-souza | 42612 | 57327 | [README](./stress-test/rinha-de-backend-2023-q3/participantes/isadora-souza/README.md) |
- | 8 | jrodrigues | 41193 | 44445 | [README](./stress-test/rinha-de-backend-2023-q3/participantes/jrodrigues/README.md) |
- | 23 | luanpontes100 | 21315 | 54779 | [README](./stress-test/rinha-de-backend-2023-q3/participantes/luanpontes100/README.md) |


## Database with Postgres

### Configuration File
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


### Connection Pool
When working with connection pools in PostgreSQL, it’s important to configure the maximum number of connections both in the PostgreSQL server itself and in the application using the database. Here’s how I’ve configured it:
- **SetMaxOpenConns(n int)**: This method sets the maximum number of open connections to the database that the application can have at any given time. If all connections are in use and a new request is made, that request will block until a connection is freed up. This setting should be chosen carefully based on your database’s resources and your application’s needs.
- **SetMaxIdleConns(n int)**: This method sets the maximum number of idle connections that can be retained in the connection pool. Idle connections can be reused for future queries without the overhead of opening a new connection. Setting this appropriately helps maintain a balance between resource usage and the ability to handle bursts of traffic efficiently.

By configuring these settings properly, you can ensure that your application efficiently manages database connections, balancing performance with resource usage.


### Tables
I didn't thought to create a UNIQUE column in the table, it was better then performing one query to check the data before inserting in the table. This was a nice learning.

Also, creating a column for searching using Postgres to auto generate was a nice thing I've learned too.

### General Configuration

The **shared_buffers** setting in PostgreSQL controls how much memory is allocated for caching data pages. The default value is 128 MB, but it’s often increased in production environments to 25-40% of system RAM for better performance.


### Indexes and Search

#### LIKE vs ILIKE
I was using the ILIKE operation in the new search index instead of LIKE with all the case in lower on the index, seems that this was causing a major performance issue on the database side, removing it provided a 12k increase in the persons inserted in the database.

#### Concorrency
When creating an index with the `CONCURRENTLY` option, PostgreSQL builds the index without locking the table for writes. This allows other operations to continue on the table while the index is being created. This is particularly useful for large tables or production environments where downtime needs to be minimized.


## Nginx

### Connection Pool
Nginx uses worker processes to handle incoming connections. Each worker can handle many connections simultaneously using an asynchronous, non-blocking I/O model. For upstream connections (connections to backend servers), Nginx maintains a pool of connections to these servers, which it reuses to handle multiple client requests efficiently.

The **keepalive** directive in the upstream block sets the maximum number of idle keepalive connections to upstream servers. This helps in reusing connections for multiple requests, reducing connection setup overhead.

## Docker

### Network

There seems to be an increased latency in docker when using network mode as bridge, because this creates a virtual network over the real one. For this kind of application, the host mode seems to be more effective.

But as I've found out, the host mode for network using Docker doesn't work on MacOS. I can research later how to use Linux on MacOS.

This is also the fault of something that was killing me, the **j.i.lOException Premature close**. This seemed to be the fault of nginx, but it's actually Docker's bridge mode.