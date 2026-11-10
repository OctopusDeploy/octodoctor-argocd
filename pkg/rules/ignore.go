package rules

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

// defaultConfigFile is read from the working directory when --config isn't given
const defaultConfigFile = ".octodoctor.yaml"

// IgnoreEntry accepts the failures of a rule, either all of them or only those of some files
type IgnoreEntry struct {
	Rule      string    `yaml:"rule"`
	Resources yaml.Node `yaml:"resources"` // a node, so that an empty or null list is told apart from a missing one
	Reason    string    `yaml:"reason"`

	whole    bool
	patterns []string
	regexps  []*regexp.Regexp
	hits     []int // matches per pattern, or a single counter for a whole-rule entry
}

// LoadIgnores reads and validates the ignore entries of a configuration file. With an empty path it
// reads .octodoctor.yaml from the working directory, and returns no entries if that file doesn't exist.
func LoadIgnores(path string) ([]IgnoreEntry, error) {
	explicit := path != ""
	if !explicit {
		path = defaultConfigFile
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) && !explicit {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var config struct {
		Ignore []IgnoreEntry `yaml:"ignore"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("invalid configuration file %s: %w", path, err)
	}

	known := map[string]bool{}
	for _, rule := range NewRegistry().Rules {
		known[rule.Metadata().ID] = true
	}
	wholeByRule := map[string]bool{}
	for i := range config.Ignore {
		e := &config.Ignore[i]
		if !known[e.Rule] {
			return nil, fmt.Errorf("invalid configuration file %s: ignore entry %d: unknown or missing rule %q", path, i+1, e.Rule)
		}
		e.whole = e.Resources.Kind == 0
		if !e.whole {
			if err := e.Resources.Decode(&e.patterns); err != nil || len(e.patterns) == 0 {
				return nil, fmt.Errorf("invalid configuration file %s: ignore entry %d: resources is empty, leave it out to ignore the whole rule", path, i+1)
			}
		}
		if whole, seen := wholeByRule[e.Rule]; seen && (whole || e.whole) {
			return nil, fmt.Errorf("invalid configuration file %s: %s is ignored as a whole more than once, or both as a whole and by resources", path, e.Rule)
		}
		wholeByRule[e.Rule] = e.whole
		for _, p := range e.patterns {
			e.regexps = append(e.regexps, regexp.MustCompile("^"+strings.ReplaceAll(regexp.QuoteMeta(toSlash(p)), `\*`, ".*")+"$"))
		}
		e.hits = make([]int, max(len(e.patterns), 1))
	}

	log.Printf("Using configuration file %s with %d ignore entries", path, len(config.Ignore))
	return config.Ignore, nil
}

// toSlash treats a backslash as a slash, so paths match the same way on every OS
func toSlash(path string) string {
	return strings.ReplaceAll(path, `\`, "/")
}

// applyIgnores removes the files of a failed rule that the entries ignore from exec, and returns one
// result per entry that ignored something. allIgnored reports that the rule has no failure left.
func applyIgnores(entries []IgnoreEntry, exec *analysis.RuleExecution) (ignored []analysis.IgnoredExecution, allIgnored bool) {
	byEntry := make([][]string, len(entries))
	var remaining []string
	for i := range entries {
		if e := &entries[i]; e.Rule == exec.ID && e.whole {
			e.hits[0]++
			return []analysis.IgnoredExecution{{RuleExecution: *exec, Reason: e.Reason}}, true
		}
	}
	for _, file := range exec.FilesTriggered {
		claimedBy := -1
		for i := range entries {
			if entries[i].Rule != exec.ID {
				continue
			}
			for j, re := range entries[i].regexps {
				if re.MatchString(toSlash(file)) {
					entries[i].hits[j]++
					if claimedBy < 0 {
						claimedBy = i
					}
				}
			}
		}
		if claimedBy < 0 {
			remaining = append(remaining, file)
		} else {
			byEntry[claimedBy] = append(byEntry[claimedBy], file)
		}
	}
	for i, files := range byEntry {
		if len(files) > 0 {
			result := analysis.IgnoredExecution{RuleExecution: *exec, Reason: entries[i].Reason}
			result.FilesTriggered = files
			ignored = append(ignored, result)
		}
	}
	exec.FilesTriggered = remaining
	return ignored, len(ignored) > 0 && len(remaining) == 0
}

// warnUnusedIgnores logs every whole-rule entry and every pattern that ignored nothing
func warnUnusedIgnores(entries []IgnoreEntry) {
	for _, e := range entries {
		for j, hits := range e.hits {
			if hits > 0 {
				continue
			}
			if e.whole {
				log.Printf("Ignore entry for %s had no effect, the rule didn't fail", e.Rule)
			} else {
				log.Printf("Ignore pattern %q for %s matched no failed file", e.patterns[j], e.Rule)
			}
		}
	}
}
