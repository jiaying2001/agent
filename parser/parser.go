package parser

import (
	"fmt"
	"regexp"
)

type Parser interface {
	Parse(log string) (string, string)
}

type RFC5424Parser struct {
}

func (p *RFC5424Parser) Parse(log string) (string, string) {
	pattern := `(\w+)\[(\d+)\]`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(log)
	if len(matches) == 3 {
		appName := matches[1]
		pid := matches[2]
		return appName, pid
	} else {
		fmt.Println("No match found!")
		return ``, ``
	}
}
