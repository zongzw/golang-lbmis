# lbmis

This project implements a simple mis program to store relations between loadbalancer and bigip.

## Usage

   ```
    $ ./lbmis --help
    Usage of ./lbmis:
    -dbpath string
            The database path. (default "/path/to/program/dir/mapping.db")
    -port string
            The port to listen. (default "8080")
   ```

   It uses RESTful API to provide abilities:

   ```
    GET Requests:
        GET /mapping?loadbalancer=<id>
        GET /mapping?bigip=<id>
        GET /mapping?loadbalancer=<id>&bigip=<id>
    GET Response: 200
        [
            {
                <gorm.models>
                loadbalancer: <id>
                bigip: <id>
            }
        ]

    Response Example:
        {
            "error": "",
            "result": [
                {
                    "ID": 1001,
                    "CreatedAt": "2020-09-19T20:05:17.960326+08:00",
                    "UpdatedAt": "2020-09-19T20:05:17.960326+08:00",
                    "DeletedAt": {
                        "Time": "0001-01-01T00:00:00Z",
                        "Valid": false
                    },
                    "loadbalancer": "204bdf84-af71-4e5b-962e-5456f06e70c9",
                    "bigip": "204bdf84-af71-4e5b-962e-5458f26e72c5"
                }
            ]
        }
   ```

   ```
    POST Request:
        POST /mapping
        body: 
            {
                lbid: <lbid>
                bipid: <bipid>
            }
    POST Response: 201
        {
            "error": ""
        }
   ```

   ```
    DELETE Request:
        DELETE /mapping?lbid=<lbid>&bipid=<bipid>
    DELETE Response: 202
        {
            "error": ""
        }
   ```

## Performance

| amount | GET | POST | DELETE |
| ------ | --- | ---- | ------ |
| 0 | 325.21µs | 3.46ms | 495.91µs |
| 1k | 474.38µs | 3.39ms | 644.876µs |
| 100k | 32.92ms  | 3.60ms  | 20.93ms  |
| 1m | 254.96ms  | 3.08ms  | 177.78ms  |

## DevOps

Run `go build .` to build the binary, `go mod` will handle the dependencies.
```
$ go build . 

or # or CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
```

If it prompts error message like
```
# github.com/mattn/go-sqlite3
exec: "gcc": executable file not found in $PATH
```

You need to install gcc

```
$ yum install -y gcc
```

The binary `lbmis` has some dependencies:

```
# ldd lbmis
	linux-vdso.so.1 =>  (0x00007ffda3bb5000)
	libdl.so.2 => /lib64/libdl.so.2 (0x00007f4f2c976000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007f4f2c75a000)
	libc.so.6 => /lib64/libc.so.6 (0x00007f4f2c38b000)
	/lib64/ld-linux-x86-64.so.2 (0x00007f4f2cb81000)
```

However, they are common .so files available for all linux platforms.


