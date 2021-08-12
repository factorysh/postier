# Postier

Save and serve all POST requests received on a defined HTTP endpoint (WIP).

## Build

```shell
make
```

## Run

```shell
./bin/postier
2021/08/12 17:09:36 Warning, no LISTEN_URL provided, using default one (0.0.0.0:8042)
2021/08/12 17:09:36 Warning, no HISTORY_ENDPOINT provided, using default one (/postier-history)
2021/08/12 17:09:36 Info, history url /postier-history
```

## Use

Make post requests

```shell
curl -d "test" -v http://localhost:8042/test
*   Trying 127.0.0.1:8042...
* Connected to localhost (127.0.0.1) port 8042 (#0)
> POST /test HTTP/1.1
> Host: localhost:8042
> User-Agent: curl/7.78.0
> Accept: */*
> Content-Length: 4
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Thu, 12 Aug 2021 15:14:48 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

Fetch history

```shell
curl http://localhost:8042/postier-history | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   287  100   287    0     0   161k      0 --:--:-- --:--:-- --:--:--  280k
{
  "requests": [
    {
      "timestamp": "2021-08-12T17:14:48.165194467+02:00",
      "url": "/test",
      "headers": {
        "Accept": [
          "*/*"
        ],
        "Content-Length": [
          "4"
        ],
        "Content-Type": [
          "application/x-www-form-urlencoded"
        ],
        "User-Agent": [
          "curl/7.78.0"
        ]
      },
      "method": "POST",
      "body": "test"
    }
  ],
  "non_post_requests": 0,
  "read_body_errors": 0
}
```
