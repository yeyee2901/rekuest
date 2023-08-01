# HTTP Request, Simplified

A simple package I made, to simplify your flow in making HTTP request for an
API call.

I have provided many common options. Have any advices? Just make a PR

# Examples

```go
package main

import (
    "github.com/yeyee2901/httprequest"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    httpStatus, err := request.HTTPRequest(
        method, host+endpoint, nil,
        request.WithContext(ctx),
        request.WithCustomHTTPClient(new(http.Client)),

        request.WithHeader("Authorization", "Bearer <some_token>"),

        // query string
        request.WithQuery(key1, val1),
        request.WithQuery(key2, val2),
        request.WithQuery(key3, val3),
        request.WithQuery(key4, val4),

        // dump request & response
        request.WithResponseDump(os.Stdout),
        request.WithRequestDump(os.Stdout),

        // body payload
        request.WithJSON(jsonPayload),
    )

    if err != nil {
        log.Fatalln(err)
    }

    fmt.Println("Status:", httpStatus)
}
```
