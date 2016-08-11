package simpletpl

import (
	"testing"
)

func TestRender(t *testing.T) {
	tpl,err := NewTemplate("{{channel:1}}P10{{ID:10}}", "")
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string]interface{})
	data["channel"] = "1"
	
	tcases := make(map[int]string)
	tcases[100] = "0000000100"
	tcases[0] = "0000000000"
	tcases[999999999] = "0999999999"
	tcases[1000000000] = "1000000000"
	tcases[9999999999] = "9999999999"
	tcases[-1] = "-000000001"
	tcases[10000000000] = "0000000000"
	tcases[10000000001] = "0000000001"
	tcases[-10000000001] = "-000000001"
	for i,suf := range tcases {
		data["ID"] = i
		id,_,err := tpl.Render(data)
		if err != nil {
			t.Fatal(err)
		}
		if id != "1P10" + suf {
			t.Fatalf("expecting: %s, got: %s", "1P10" + suf, id)
		}
	}
}