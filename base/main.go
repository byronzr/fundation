package base

import (
	"encoding/json"
	"fmt"
)

// Dump dump everything
func Dump(o ...interface{}) {
	buf := make([]byte, 0)
	for idx, s := range o {
		out, err := json.MarshalIndent(s, "", "\t")
		if err != nil {
			panic(err)
		}
		buf = append(buf, []byte(fmt.Sprintf("\n\033[33m----%02d----------------------------------------------\033[0m\n%s", idx, string(out)))...)
	}
	fmt.Println(string(buf))
}

func Find(ss []string, s string) int {
	for idx, str := range ss {
		if str == s {
			return idx
		}
	}
	return -1
}
