package pytools

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// extractImportBodyParts accepts a python import
// line and will return a list of the individual body
// imports. This will format "from" imports to act
// as if they were standard imports. It assumes that
// the given line is a valid import line.
func extractImportBodyParts(l string) []string {
	parts := strings.Split(l, " ")
	header := parts[0]

	// Create replacers
	spaceRepl := strings.NewReplacer(" ", "")

	var bodyParts []string
	if header == "import" {
		body := strings.Join(parts[1:], " ")
		body = spaceRepl.Replace(body)
		bodyParts = strings.Split(body, ",")
	} else if header == "from" {
		base := parts[1]
		body := strings.Join(parts[3:], " ")
		body = spaceRepl.Replace(body)
		bodyParts = strings.Split(body, ",")

		newBody := []string{base}
		for _, p := range bodyParts {
			newBody = append(newBody, base+"."+p)
		}

		bodyParts = newBody
	}

	return bodyParts
}

// GetImportLines accepts a path to a script and will
// return all lines that begin with the word "import"
// or "from".
func GetImportLines(scr string) []string {
	bs, err := ioutil.ReadFile(scr)
	if err != nil {
		panic(err)
	}

	var importLines []string
	lines := strings.Split(string(bs), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		parts := strings.Split(l, " ")
		if parts[0] == "import" || parts[0] == "from" {
			importLines = append(importLines, l)
		}
	}

	return importLines
}

// StandardizeImportLine accept a python import line
// and reformats it to a standard format.
func StandardizeImportLine(line string) string {
	messyReg := regexp.MustCompile(`as [^\s,]+`)
	spaceReg := regexp.MustCompile(`\s+`)
	commaReg := regexp.MustCompile(`\s+,`)

	// Clean line
	line = messyReg.ReplaceAllString(line, "")
	line = spaceReg.ReplaceAllString(line, " ")
	line = commaReg.ReplaceAllString(line, ",")

	return line
}
