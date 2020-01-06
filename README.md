# go-requests

# Requests

Requests is an elegant and simple HTTP library for Go.

```go
import "github.com/fanjindong/go-requests"

func main(){
    resp, err := requests.Get("http://example.com/ping", requests.Params{"name": "fjd"})
    if err!=nil{
        fmt.Printf("Get err: %v", err)
    }
    fmt.Println(resp.Text)

    var rMap map[string]interface{}
    var rStruct struct{Code int `json:"code"`}
    
    err = resp.Json(&rMap)
    if err!=nil{
        fmt.Printf("resp.Json to map err: %v \n", err)
    }else {
        fmt.Printf("resp.Json to map: %v \n", rMap)
    }
    
    err = resp.Json(&rStruct)
    if err!=nil{
        fmt.Printf("resp.Json to struct err: %v \n", err)
    }else {
        fmt.Printf("resp.Json to struct: %+v \n", rStruct)
    }
}
```

Out:
```shell script
{"code":0,"message":"pong"}
resp.Json to map: map[code:0 message:pong] 
resp.Json to struct: {Code:0} 
```

