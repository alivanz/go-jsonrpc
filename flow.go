package jsonrpc

import (
  "encoding/json"
  "sync"
)

type rpc struct{
  codec   Codec
  methods map[string] func(params json.RawMessage) *JSONMessage
  methodnotfound func(method string, params json.RawMessage) *JSONMessage
  pending map[uint64](chan JSONMessage)
  wmutex  sync.Mutex
  seq     uint64
}
type BatchRequest struct{
  Method  string
  Params  interface{}
}

func NewRPC(c Codec) RPC {
  return &rpc{
    codec   : c,
    methods : make( map[string] func(params json.RawMessage) *JSONMessage ),
    pending : make( map[uint64] (chan JSONMessage) ),
    methodnotfound: func(method string, params json.RawMessage) *JSONMessage {
      return &MethodNotFound
    },
  }
}
func (r *rpc) Serve() error {
  for{
    var msg JSONMessage
    var bmsg []JSONMessage
    if err := r.codec.ReadMessage(&msg); err==nil{
      // Single
      resp := r.incoming(msg)
      if resp!=nil{
        if err:=r.codec.WriteMessage(resp); err!=nil{
          return err
        }
      }
    }else if err := r.codec.ReadMessage(&bmsg); err==nil{
      // Batch
      for _,msg := range bmsg{
        resp := r.incoming(msg)
        if resp!=nil{
          if err:=r.codec.WriteMessage(resp); err!=nil{
            return err
          }
        }
      }
    }else{
      resp := JSONMessage{
        Version: "2.0",
        Error: &ErrorObject{
          Code: -32700,
          Message: err.Error(),
        },
      }
      if err:=r.codec.WriteMessage(resp); err!=nil{
        return err
      }
    }
  }
}
func (r *rpc) OnNotFound(callback func(method string, params json.RawMessage) *JSONMessage) {
  if callback != nil{
    r.methodnotfound = callback
  }
}
func (r *rpc) incoming(msg JSONMessage) *JSONMessage {
  if msg.Method!=""{
    // notify / call
    if msg.Id != nil{
      r.wmutex.Lock()
      r.seq = r.seq+1
      r.wmutex.Unlock()
    }
    if fx,found := r.methods[msg.Method]; !found {
      resp := r.methodnotfound(msg.Method, msg.Params)
      //resp := JSONMessage{
      //  Id: msg.Id,
      //  Version: "2.0",
      //  Error: &ErrorObject{
      //    Code: -32601,
      //    Message: "method not found",
      //  },
      //}
      if msg.Id == nil{
        return nil
      }else if resp==nil{
        return &InternalError
      }
      resp.Id = msg.Id
      resp.Version = "2.0"
      return resp
    }else if resp := fx(msg.Params); resp!=nil{
      resp.Id = msg.Id
      resp.Version = "2.0"
      return resp
    }
  }else{
    // result
    if msg.Id==nil{
      // Id nil, no method
      //log.Print(msg)
      //return nil
      resp := JSONMessage{
        Version: "2.0",
        Error: &ErrorObject{
          Code: -32600,
          Message: "id and method are empty",
        },
      }
      return &resp
    }else if cmsg,found := r.pending[*msg.Id]; !found {
      return &IDNotFound
    }else{
      delete(r.pending, *msg.Id)
      cmsg <- msg
    }
  }
  return nil
}
func (r *rpc) ACall(method string, params interface{}) (chan JSONMessage,error) {
  jparams,err := json.Marshal(params)
  if err!=nil{
    return nil,err
  }
  r.wmutex.Lock()
  id := r.seq
  r.seq = r.seq+1
  r.wmutex.Unlock()
  msg := JSONMessage{
    Version : "2.0",
    Id      : &id,
    Method  : method,
    Params  : json.RawMessage(jparams),
  }
  if err:=r.codec.WriteMessage(msg); err!=nil{
    return nil,err
  }
  cmsg := make(chan JSONMessage)
  r.pending[id] = cmsg
  return cmsg, nil
}
func (r *rpc) Call(method string, params interface{}) (*JSONMessage,error) {
  if cmsg,err:=r.ACall(method,params); err!=nil{
    return nil,err
  }else{
    msg := <-cmsg
    return &msg, nil
  }
}
func (r *rpc) BatchCall(batch []BatchRequest) ([](chan JSONMessage),error) {
  var requests []JSONMessage
  var chans [](chan JSONMessage)
  r.wmutex.Lock()
  defer r.wmutex.Unlock()
  for _,req := range batch{
    jparams,err := json.Marshal(req.Params)
    if err!=nil{
      return nil,err
    }
    id := r.seq
    r.seq = r.seq+1
    msg := JSONMessage{
      Version : "2.0",
      Id      : &id,
      Method  : req.Method,
      Params  : json.RawMessage(jparams),
    }
    cmsg := make(chan JSONMessage)
    r.pending[id] = cmsg
    requests = append(requests, msg)
    chans = append(chans, cmsg)
  }
  if err:=r.codec.WriteMessage(requests); err!=nil{
    return nil,err
  }
  return chans, nil
}
func (r *rpc) Notify(method string, params interface{}) error {
  jparams,err := json.Marshal(params)
  if err!=nil{
    return err
  }
  msg := JSONMessage{
    Version : "2.0",
    Id      : nil,
    Method  : method,
    Params  : json.RawMessage(jparams),
  }
  return r.codec.WriteMessage(msg)
}

func (r *rpc) OnRequest(method string, callback func(params json.RawMessage) *JSONMessage) {
  r.methods[method] = callback
}
