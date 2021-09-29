package gconf

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(gconf(":url[https://github.com]\n"+
		":text[cd /tmp]", "url"))
}

func gconf(gconfs string, keyword string) string {

	for _, line := range strings.Split(strings.TrimRight(gconfs, "\n"), "\n") {

		if string(line[0:len(keyword)+1]) == ":"+keyword {
			text := line[:len(line)-1]
			text = text[len(keyword)+2:]
			return text
		}
	}
	panic("should never happen")

}
