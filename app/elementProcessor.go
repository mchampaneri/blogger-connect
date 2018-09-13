package main

import (
	"fmt"
)

func ListItem(class, id, data string) string {
	return fmt.Sprintf(`<li class="%s" id="%s" >%s<li>`, class, id, data)
}
