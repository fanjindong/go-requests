# Requests
![](./doc/images/TrakaiLithuania_ZH-CN0447602818_1920x1080.jpg)

Requests is an elegant and simple HTTP library for Go.

## Install

```shell script
go get github.com/fanjindong/go-requests
```

## QuickStart

### Make a Request

Making a request with Requests is very simple.
For this example:

```go
resp, err := requests.Get("http://example.com/ping", requests.Params{"key": "value"})
```
Now, we have a Response object called resp. We can get all the information we need from this object.

Requests’ simple API means that all forms of HTTP request are as obvious. For example, this is how you make an HTTP POST request:

```go
resp, err := requests.Post("http://example.com/ping", requests.Params{"k": "v"}, requests.Json{"key": "value"})
```

What about the other HTTP request types: PUT, DELETE, HEAD and OPTIONS? These are all just as simple:

```go
resp, err := requests.Put("http://example.com/ping", requests.Data{"key": "value"})
resp, err := requests.Delete("http://example.com/ping")
resp, err := requests.Head("http://example.com/ping")
resp, err := requests.Options("http://example.com/ping")
```
That’s all well and good, but it’s also only the start of what Requests can do.

### Passing Parameters In URLs

You often want to send some sort of data in the URL’s query string. 
If you were constructing the URL by hand, this data would be given as key/value pairs in the URL after a question mark, 
e.g. example.com/get?key=val. Requests allows you to provide these arguments as a dictionary of strings, 
using the params keyword argument. As an example, if you wanted to pass key1=value1 and key2=value2 to example.com/get, 
you would use the following code:

```go
resp, err := requests.Get("http://example.com/get", requests.Params{"key1": "value1", "key2": "value2"})
```
You can see that the URL has been correctly encoded by printing the URL:
```go
fmt.Println(resp.Request.URL)
//http://example.com/get?key2=value2&key1=value1
```

### Response Content

We can read the content of the server’s response. Consider the GitHub timeline again:

```go
resp, _ := requests.Get("https://api.github.com/events")
fmt.Println(resp.Text)
//{"code":0,"message":"pong"}
```

### JSON Response Content

There’s also a builtin JSON decoder, in case you’re dealing with JSON data:

```go
var rStruct struct{
    Code int `json:"code"`
    Message string `json:"message"`
}

err := resp.Json(&rStruct)
fmt.Printf("resp.Json to struct: %+v \n", rStruct)
// resp.Json to struct: {Code:0, Message:"success"} 
```

### Custom Headers

If you’d like to add HTTP headers to a request, simply pass in a `requests.Headers` to the headers parameter.

For example, we did not specify our user-agent in the previous example:

```go
r, err := requests.Get("https://api.github.com/some/endpoint", requests.Headers{"user-agent": "my-app/0.0.1"})
```

### More complicated POST requests

Typically, you want to send some form-encoded data — much like an HTML form. To do this, 
simply pass a `requests.Data` to the data argument. 
Your data will automatically be form-encoded when the request is made:

```go
r, err := requests.Post("https://httpbin.org/post", requests.Data{"key1": "value1", "key2": "value2"})
fmt.Println(r.Text)
//{"code":0,"message":"pong"}
```

For example, the GitHub API v3 accepts JSON-Encoded POST/PATCH data, 
you can also pass it `requests.Json` using the json parameter and it will be encoded automatically:

```go
r, err := requests.Post("https://api.github.com/some/endpoint", requests.Json{"key1": "value1", "key2": "value2"})
```

Using the `requests.Json` in the request will change the Content-Type in the header to application/json.

### POST a Multipart-Encoded File

```go
file, err := requests.FileFromPath("demo.text")

r, err := requests.Post("https://httpbin.org/post", requests.Files{"key": "value", "file": file})
```

### Response Status Codes

We can check the response status code:

```go
r, err := requests.Get("https://httpbin.org/get")
fmt.Println(r.StatusCode)
// 200
```

### Response Headers

We can view the server’s response headers:

```go
fmt.Println(r.Headers)
//map[Cache-Control:[private] Content-Type:[application/json] Set-Cookie:[QINGCLOUDELB=d9a2454c187d2875afb6701eb80e9c8761ebcf3b54797eae61b25b90f71273ea; path=/; HttpOnly]]

```
We can access the headers using Get method:

```go
r.Headers.Get("Content-Type")
//"application/json"
```

### Timeouts

You can tell Requests to stop waiting for a response after a given number of seconds with the timeout parameter. 
Nearly all production code should use this parameter in nearly all requests. 
Failure to do so can cause your program to hang indefinitely:


```go
r, err := requests.Get("https://github.com/", requests.Params{"key": "value"}, requests.Timeout(3*time.Secend))
```