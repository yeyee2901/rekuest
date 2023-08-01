# HTTP Request, Simplified

A simple package I made, to simplify your flow in making HTTP request for an
API call.

I have provided many common options. Have any advices? Just make a PR

# Examples

```go
package main

import (
    "github.com/yeyee2901/rekuest"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    httpStatus, err := rekuest.HTTPRequest(
        method, host+endpoint, nil,
        rekuest.WithContext(ctx),
        rekuest.WithCustomHTTPClient(new(http.Client)),

        rekuest.WithHeader("Authorization", "Bearer <some_token>"),

        // query string
        rekuest.WithQuery(key1, val1),
        rekuest.WithQuery(key2, val2),
        rekuest.WithQuery(key3, val3),
        rekuest.WithQuery(key4, val4),

        // dump request & response
        rekuest.WithResponseDump(os.Stdout),
        rekuest.WithRequestDump(os.Stdout),

        // body payload
        rekuest.WithJSON(jsonPayload),
    )

    if err != nil {
        log.Fatalln(err)
    }

    fmt.Println("Status:", httpStatus)
}
```
