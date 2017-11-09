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

### [MessagePack message generator]("github.com/tinylib/msgp") interoperability
You should see in debug output two JSON strings original and double converted. Note in dictionary type fields order is not guaranteed(see WTF chapter below) 
```go
import (
	"bytes"
	"fmt"

  json2msgpack "github.com/izinin/json2msgpack"
	"github.com/tinylib/msgp/msgp"
)

func main(){
  origStr := `{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`
  src := bytes.NewBuffer([]byte(json2msgpack.EncodeJSON([]byte(origStr)))) 
  var js bytes.Buffer
  _, err := msgp.CopyToJSON(&js, src)
  if err != nil {
  	panic(fmt.Sprintf("Cannot convert MessagePack to JSON: %v", err))
  }
  fmt.Printf("Original JSON: \n\t%s", origStr)
  fmt.Printf("Converted string JSON --> MsgPack --> JSON: \n\t%s", js.String())
}
```

### Golang WTF limitation
Please note JSON unmarshalling **does not keep dictionary order**, then we use dictionary sorted by key in alphabetical order just to be deterministic. [Here is the line.](https://github.com/izinin/json2msgpack/blob/968f39ee8e4d5b8225d210a86db30b8bab030ac6/json2msgpack.go#L177)

