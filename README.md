# json2msgpack

## This is implementation of generic JSON serialiser to [MessagePack](https://msgpack.org/) in Golang
It complements [MessagePack message generator]("github.com/tinylib/msgp")

### The problem it solves
[MessagePack message generator]("github.com/tinylib/msgp") does not support JSON --> messagePack conversion. It allows JSON iteroperability by just providing [data memory to JSON conversion](https://godoc.org/github.com/tinylib/msgp/msgp#CopyToJSON) but not backward (see `msgp.CopyToJSON() and msgp.UnmarshalAsJSON()`)
In some cases you do not need message generator for each and every message type you use in your project but rather translation your JSON data to MessagePack. It is useful in your project MessagePack API testing

### How to use it 
```go
package main

import (
  "fmt"
  json2msgpack "github.com/izinin/json2msgpack"
)

func main(){
  res := json2msgpack.EncodeJSON([]byte(`{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`))
  fmt.Printf("MessagePack format: %#x", test)
}

```
it outputs MessagePack representation: 
`0x82a7636f6d70616374c3a6736368656d619281a46e616d65a469676f7281a5666e616d65a57a696e696e`


### How to test 
You run Golang test by command line `go test -v`

MessagePack data conversion for the test criterias are taken from [MessagePack console, Try! menu item](https://msgpack.org/)
