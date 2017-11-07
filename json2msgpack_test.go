package json2msgpack

import (
	"flag"
	"testing"
)

func TestJSONParse(t *testing.T) {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	res := EncodeJSON([]byte(`{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`))
	t.Logf("Final result: %#x", res)
}
