// From TCP Connection
//    con,err := net.Dial("tcp", "127.0.0.1:10000")
//    if err!=nil{
//      log.Fatal(err)
//    }
//    rpc := jsonrpc.NewRPC( jsonrpc.NewJSONCodec(con) )
//
// Serve method
//    rpc.OnRequest("sum", func(data json.RawMessage) *jsonrpc.JSONMessage {
//      var numbers []int
//      if err:=json.Unmarshal(data,&numbers); err!=nil{
//        return &jsonrpc.JSONMessage{
//          Error: &jsonrpc.ErrorObject{
//            Message: err.Error(),
//          },
//        }
//      }else{
//        var result int
//        result = 0
//        for _,number := range numbers {
//          result = result + number
//        }
//        msg,_ := jsonrpc.NewJSONMessage(result)
//        return msg
//      }
//    })
//
// Call Method "sum"
//    if result,err := rpc.Call("sum", []int{1,3,6}); err!=nil{
//      log.Print(err)
//    }else{
//      log.Print(string(result.Result))
//    }
package jsonrpc

import (
  "encoding/json"
)

type RPC interface{
  // Blocking Call function
  Call(method string, params interface{}) (*JSONMessage,error)
  // Non-Blocking Call function
  ACall(method string, params interface{}) (chan JSONMessage,error)
  // Notify function
  Notify(method string, params interface{}) error
  // Batch Call
  BatchCall(batch []BatchRequest) ([](chan JSONMessage), error)

  // for Call Request and Notify Request
  // You SHOULD know if the request is either Call or Notify
  // If Call Request you SHOULD return the reply, or the request never replied.
  // If Notify Request you SHOULD return nil, or zombie reply occurred.
  // You dont have to set the ID, it will set later.
  OnRequest(method string, callback func(params json.RawMessage) *JSONMessage )

  // Called if method not found
  // WARNING : THIS WILL CHANGE THE "method not found" PROCEDURE!!!
  OnNotFound(callback func(method string, params json.RawMessage) *JSONMessage)

  // Serve RPC Server
  Serve() error
}

type Codec interface{
  // Write message to connection
  WriteJSON(message interface{}) error
  // Read message from connection
  // Remember the input must be a pointer to make this work
  ReadJSON(message interface{}) error
}
