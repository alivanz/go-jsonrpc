package jsonrpc

var ParseError = JSONMessage{
  Version: "2.0",
  Error: &ErrorObject{
    Code: -32700,
    Message: "An error occurred on the server while parsing the JSON text.",
  },
}
var MethodNotFound = JSONMessage{
  Version: "2.0",
  Error: &ErrorObject{
    Code: -32601,
    Message: "The method does not exist / is not available.",
  },
}
var InvalidParams = JSONMessage{
  Version: "2.0",
  Error: &ErrorObject{
    Code: -32602,
    Message: "Invalid method parameter(s).",
  },
}
var InternalError = JSONMessage{
  Version: "2.0",
  Error: &ErrorObject{
    Code: -32603,
    Message: "Internal JSON-RPC error.",
  },
}
