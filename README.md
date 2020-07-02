# Runtime Env
`HONEYCOMB_API_KEY`: Needs to be set to the honeycomb api key


# Building for ARM (raspberry pi 4)
```
GOOS=linux GOARCH=arm GOARM=7 go build
```