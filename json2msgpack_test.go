package json2msgpack

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestJSONParse(t *testing.T) {
	res := EncodeJSON([]byte(`{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`))
	str := fmt.Sprintf("%#x", res)
	t.Logf("Base example result: %v", str)
	if str[2:] != "82a7636f6d70616374c3a6736368656d619281a46e616d65a469676f7281a5666e616d65a57a696e696e" {
		t.Fatal("mismatch expected base example")
	}
	// [ 1, "auth", {"tok": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfdmVyc2lvbiI6MSwiZGV2aWNlX2lkIjoiZDE6VTUzODlmMTcwM2VjMTQ1ZjkiLCJkZXZpY2VfaW5kZXgiOjEsImV4cCI6MTUxMDA2MjYzNCwianRpIjoiMDFCWUJCMllTTVpTU0FWRjhWNUU5MlFES0giLCJwbGF0Zm9ybSI6ImFuZHJvaWQiLCJ1c2VyX2lkIjoiVTUzODlmMTcwM2VjMTQ1ZjkifQ.EcIKJd2cCxSpupF0-tP9bfinndMi275MmUTwWVS-1bE"}, null ]
	str = fmt.Sprintf("%#x",
		EncodeJSON([]byte(`[ 1, "auth", {"tok": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfdmVyc2lvbiI6MSwiZGV2aWNlX2lkIjoiZDE6VTUzODlmMTcwM2VjMTQ1ZjkiLCJkZXZpY2VfaW5kZXgiOjEsImV4cCI6MTUxMDA2MjYzNCwianRpIjoiMDFCWUJCMllTTVpTU0FWRjhWNUU5MlFES0giLCJwbGF0Zm9ybSI6ImFuZHJvaWQiLCJ1c2VyX2lkIjoiVTUzODlmMTcwM2VjMTQ1ZjkifQ.EcIKJd2cCxSpupF0-tP9bfinndMi275MmUTwWVS-1bE"}, null ]`)))

	if str[2:] != "9401a46175746881a3746f6bda013b65794a68624763694f694a49557a49314e694973496e523563434936496b705856434a392e65794a6a62476c6c626e5266646d567963326c76626949364d5377695a47563261574e6c58326c6b496a6f695a4445365654557a4f446c6d4d5463774d32566a4d5451315a6a6b694c434a6b5a585a70593256666157356b5a5867694f6a4573496d5634634349364d5455784d4441324d6a597a4e437769616e5270496a6f694d44464357554a434d6c6c545456705455304657526a68574e5555354d6c4645533067694c434a77624746305a6d397962534936496d46755a484a76615751694c434a316332567958326c6b496a6f695654557a4f446c6d4d5463774d32566a4d5451315a6a6b6966512e4563494b4a64326343785370757046302d7450396266696e6e644d693237354d6d5554775756532d316245c0" {
		t.Fatal("mismatch expected case_0")
	}
	// [ 2, "cj", {"c": "Pmygroup"}, null ]
	str = fmt.Sprintf("%#x",
		EncodeJSON([]byte(`[ 2, "cj", {"c": "Pmygroup"}, null ]`)))

	if str[2:] != "9402a2636a81a163a8506d7967726f7570c0" {
		t.Fatal("mismatch expected case_1")
	}
	// [ 3, "m", {"i": "5191747051.1", "r": "Pmygroup"}, "\"hello kitty!\"" ]
	str = fmt.Sprintf("%#x",
		EncodeJSON([]byte(`[ 3, "m", {"i": "5191747051.1", "r": "Pmygroup"}, "\"hello kitty!\"" ]`)))

	if str[2:] != "9403a16d82a169ac353139313734373035312e31a172a8506d7967726f7570ae2268656c6c6f206b697474792122" {
		t.Fatal("mismatch expected case_2")
	}
	// [ 5, "m", {"f": null, "i": "5191747051.1", "r": "Pmygroup", "s": "U5389f1703ec145f9", "t": 1510061749068}, "\"hello kitty!\"" ]
	str = fmt.Sprintf("%#x",
		EncodeJSON([]byte(`[ 5, "m", {"f": null, "i": "5191747051.1", "r": "Pmygroup", "s": "U5389f1703ec145f9", "t": 1510061749068}, "\"hello kitty!\"" ]`)))

	if str[2:] != "9405a16d85a166c0a169ac353139313734373035312e31a172a8506d7967726f7570a173b15535333839663137303365633134356639a174cf0000015f96b1b34cae2268656c6c6f206b697474792122" {
		t.Fatal("mismatch expected case_3")
	}
}

func TestMessagePackInterop(t *testing.T) {
	origStr := `{"compact":true,"schema":[{"name": "igor"}, {"fname": "zinin"}]}`
	src := bytes.NewBuffer([]byte(EncodeJSON([]byte(origStr))))

	var js bytes.Buffer
	_, err := msgp.CopyToJSON(&js, src)
	if err != nil {
		t.Fatalf("Cannot convert MessagePack to JSON: %v", err)
	}
	t.Logf("Original JSON: \n\t%s", origStr)
	t.Logf("Converted string JSON --> MsgPack --> JSON: \n\t%s", js.String())
}
