package analyzer

import (
	"os"
	"regexp"
)

type PHPClass struct {
	ClassName   string
	Methods     []string
	Constructor bool
	Code        string
}

func ParsePhpFiles(path string) (*PHPClass, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	code := string(data)

	classRegex := regexp.MustCompile(`class\s+(\w+)`)
	methodRegex := regexp.MustCompile(`function\s+(\w+)`)

	classMatches := classRegex.FindStringSubmatch(string(data))
	methodmatches := methodRegex.FindAllStringSubmatch(string(data), -1)

	var methods []string
	hasConstructor := false
	for _, m := range methodmatches {
		methods = append(methods, m[1])
		if m[1] == "__construct" {
			hasConstructor = true
		}
	}

	if len(classMatches) < 2 {
		return nil, nil
	}

	return &PHPClass{
		ClassName:   classMatches[1],
		Methods:     methods,
		Constructor: hasConstructor,
		Code:        code,
	}, nil

}
