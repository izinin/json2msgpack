package json2msgpack

import (
	"testing"
)

func TestJSONParse(t *testing.T) {
	res := EncodeJSON([]byte(`{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`))
	t.Logf("Final result: %#x", res)
}
