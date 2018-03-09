# JSON-RPC 2.0 for Go

## Feature
* Bi-directional
* Asynchronous
* Batch call support

## Example

From TCP Connection

```Go
con,err := net.Dial("tcp", "127.0.0.1:10000")
if err!=nil{
  log.Fatal(err)
}
rpc := jsonrpc.NewRPC( jsonrpc.NewJSONCodec(con) )
```

Serve method "sum"

```Go
rpc.OnRequest("sum", func(data json.RawMessage) *jsonrpc.JSONMessage {
  var numbers []int
  if err:=json.Unmarshal(data,&numbers); err!=nil{
    return &jsonrpc.JSONMessage{
      Error: &jsonrpc.ErrorObject{
        Message: err.Error(),
      },
    }
  }else{
    var result int
    result = 0
    for _,number := range numbers {
      result = result + number
    }
    msg,_ := jsonrpc.NewJSONMessage(result)
    return msg
  }
})
```

Call Method "sum"

```Go
if result,err := rpc.Call("sum", []int{1,3,6}); err!=nil{
  log.Print(err)
}else{
  log.Print(string(result.Result))
}
```
