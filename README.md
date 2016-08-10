# simpletpl
a simple tpl render

# Usage:
```
package main

import (
	"github.com/YEXINGZHE54/simpletpl"
)

func main() {
	tpl,err := simpletpl.NewTemplate("{{channel:1}}P10{{ID:10}}", "")
	if err != nil {
		panic(err)
	}
	data := make(map[string]interface{})
	data["channel"] = "1"
	data["ID"] = 100
	id,err := tpl.Render(data)
	if err != nil {
		panic(err)
	}
	println(id)
}
```
