package pytools

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rburmorrison/go-pytools/internal/verify"
)

// GetAssociatedScripts accepts the path to an entry
// python script and will return the paths to all
// the scripts that it imports, recursively. It will
// get trapped in an infinite loop if any two python
// scripts import each other.
func GetAssociatedScripts(base string, ep string) []string {
	scripts := []string{ep}

	// Standardize all lines
	il := GetImportLines(ep)
	for i, l := range il {
		il[i] = StandardizeImportLine(l)
	}

	// Create replacers
	dotRepl := strings.NewReplacer(".", string(os.PathSeparator))

	// Analyze imports
	for _, l := range il {
		bodyParts := extractImportBodyParts(l)

		// Turn all body parts into paths
		for _, p := range bodyParts {
			p = dotRepl.Replace(p)
			path := path.Join(base, p+".py")

			if verify.FilePath(path) == nil {
				// Leads to an actual file
				scripts = append(scripts, path)
				scripts = append(scripts, GetAssociatedScripts(base, path)...)
			}
		}
	}

	return scripts
}

// GetRunScriptCommand returns the command to be
// executed to run the passed python script.
func GetRunScriptCommand(interp string, path string, args ...string) *exec.Cmd {
	command := append([]string{path}, args...)
	cmd := exec.Command(interp, command...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}

// RunScript accepts a path to a script and a python
// interpreter and returns an error if it is unable
// to run it. Output is written to os.Stdout.
func RunScript(interp string, path string, args ...string) error {
	cmd := GetRunScriptCommand(interp, path, args...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
