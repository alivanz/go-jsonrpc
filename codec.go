package jsonrpc

import (
  "encoding/json"
  "io"
  "sync"
)

type JSONMessage struct{
  Version string            `json:"jsonrpc,omitempty"`
  Id      *uint64           `json:"id"`
  Method  string            `json:"method,omitempty"`
  Params  json.RawMessage   `json:"params,omitempty"`
  Result  json.RawMessage   `json:"result,omitempty"`
  Error   *ErrorObject      `json:"error,omitempty"`
}

type ErrorObject struct{
  Code    int64           `json:"code"`
  Message string          `json:"message"`
  Data    json.RawMessage `json:"data,omitempty"`
}

type JSONCodec struct{
  rmutex  sync.Mutex
  wmutex  sync.Mutex
  decoder *json.Decoder
  encoder *json.Encoder
}

// Create new JSONMessage that Result contain JSON Marshaled context
func NewJSONMessage(content interface{}) (*JSONMessage,error) {
  if b,err:=json.Marshal(content); err!=nil{
    return nil,err
  }else{
    return &JSONMessage{
      Result: b,
      },nil
  }
}

func NewJSONCodec(con io.ReadWriteCloser) Codec {
  return &JSONCodec{
    decoder: json.NewDecoder(con),
    encoder: json.NewEncoder(con),
  }
}
func (c *JSONCodec) WriteJSON(message interface{}) error {
  c.wmutex.Lock()
  defer c.wmutex.Unlock()
  return c.encoder.Encode(message)
}
func (c *JSONCodec) ReadJSON(message interface{}) error {
  c.rmutex.Lock()
  defer c.rmutex.Unlock()
  return c.decoder.Decode(message)
}

func (e ErrorObject) MarshalJSON() ([]byte,error) {
  var x [3]json.RawMessage
  var err error
  if x[0],err=json.Marshal(e.Code); err!=nil{
    return nil,err
  }
  if x[1],err=json.Marshal(e.Message); err!=nil{
    return nil,err
  }
  if x[2],err=json.Marshal(e.Data); err!=nil{
    return nil,err
  }
  return json.Marshal(x)
}
func (e *ErrorObject) UnmarshalJSON(b []byte) error {
  var x [3]json.RawMessage
  if err:=json.Unmarshal(b,&x); err!=nil{
    return err
  }
  if err:=json.Unmarshal(x[0],&e.Code); err!=nil{
    return err
  }
  if err:=json.Unmarshal(x[1],&e.Message); err!=nil{
    return err
  }
  e.Data = x[2]
  return nil
}
func (e *ErrorObject) Error() string {
  return e.Message
}
