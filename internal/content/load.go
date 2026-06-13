package content

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed all:data
var dataFS embed.FS

// Load reads all embedded YAML content into an ordered Catalog.
func Load() (Catalog, error) {
	var c Catalog
	err := fs.WalkDir(dataFS, "data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		b, err := dataFS.ReadFile(path)
		if err != nil {
			return err
		}
		switch {
		case strings.Contains(path, "/intro/"):
			var list []Diagnostic
			if err := yaml.Unmarshal(b, &list); err == nil && len(list) > 0 {
				c.Intro = append(c.Intro, list...)
			} else {
				var di Diagnostic
				if err := yaml.Unmarshal(b, &di); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				c.Intro = append(c.Intro, di)
			}
		case strings.Contains(path, "/diagnostics/"):
			// A file may hold a single Diagnostic OR a YAML list of variants.
			var list []Diagnostic
			if err := yaml.Unmarshal(b, &list); err == nil && len(list) > 0 {
				c.Diagnostics = append(c.Diagnostics, list...)
			} else {
				var di Diagnostic
				if err := yaml.Unmarshal(b, &di); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				c.Diagnostics = append(c.Diagnostics, di)
			}
		case strings.Contains(path, "/lessons/"):
			var l Lesson
			if err := yaml.Unmarshal(b, &l); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			c.Lessons = append(c.Lessons, l)
		case strings.Contains(path, "/devliteracy/"):
			// A file may hold a single DevTask OR a YAML list (the command bank).
			var list []DevTask
			if err := yaml.Unmarshal(b, &list); err == nil && len(list) > 0 {
				c.DevTasks = append(c.DevTasks, list...)
			} else {
				var dt DevTask
				if err := yaml.Unmarshal(b, &dt); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				c.DevTasks = append(c.DevTasks, dt)
			}
		case strings.Contains(path, "/bench/"):
			var list []Problem
			if err := yaml.Unmarshal(b, &list); err == nil && len(list) > 0 {
				c.Problems = append(c.Problems, list...)
			} else {
				var p Problem
				if err := yaml.Unmarshal(b, &p); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				c.Problems = append(c.Problems, p)
			}
		case strings.Contains(path, "/advanced/"):
			// Stage-2 Advanced Topics: one topic per file.
			var at AdvancedTopic
			if err := yaml.Unmarshal(b, &at); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			c.AdvancedTopics = append(c.AdvancedTopics, at)
		case strings.Contains(path, "/install/"):
			// One install guide per language file (ADR-0007 BYO runtime).
			var ig InstallGuide
			if err := yaml.Unmarshal(b, &ig); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			c.InstallGuides = append(c.InstallGuides, ig)
		case strings.Contains(path, "/primers/"):
			var list []Primer
			if err := yaml.Unmarshal(b, &list); err == nil && len(list) > 0 {
				c.Primers = append(c.Primers, list...)
			} else {
				var pr Primer
				if err := yaml.Unmarshal(b, &pr); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				c.Primers = append(c.Primers, pr)
			}
		}
		return nil
	})
	if err != nil {
		return c, err
	}
	// Back-compat: the original primers predate the per-language Lang field, so an
	// empty Lang means the Python primer. Normalize once at load so every lookup
	// (PrimerByCategoryAndLang / PrimerLangs) can assume Lang is set.
	for i := range c.Primers {
		if c.Primers[i].Lang == "" {
			c.Primers[i].Lang = "python"
		}
		// Migration shim: a legacy primer with a flat `ops:` list (and no
		// `sections:`) becomes a single "Basic operations" section, so old and
		// new files both render through the section pager during the rollout.
		if len(c.Primers[i].Sections) == 0 && len(c.Primers[i].Ops) > 0 {
			c.Primers[i].Sections = []PrimerSection{{Title: "Basic operations", Ops: c.Primers[i].Ops}}
		}
	}
	sort.SliceStable(c.Diagnostics, func(i, j int) bool { return c.Diagnostics[i].Order < c.Diagnostics[j].Order })
	sort.SliceStable(c.Lessons, func(i, j int) bool { return c.Lessons[i].Order < c.Lessons[j].Order })
	sort.SliceStable(c.DevTasks, func(i, j int) bool { return c.DevTasks[i].Order < c.DevTasks[j].Order })
	if err := validateSlots(c.Diagnostics); err != nil {
		return c, err
	}
	return c, nil
}

// validateSlots enforces the slot invariant: every diagnostic has a slot, and
// all variants within a slot share the same Kind and Measures. A stray mismatch
// would desync the player flow, so we fail loudly at load time.
func validateSlots(ds []Diagnostic) error {
	type sig struct{ kind, measures string }
	seen := map[string]sig{}
	for _, d := range ds {
		if d.Slot == "" {
			return fmt.Errorf("diagnostic %q has no slot", d.ID)
		}
		s := sig{d.Kind, d.Measures}
		if prev, ok := seen[d.Slot]; ok {
			if prev != s {
				return fmt.Errorf("slot %q mixes kinds/measures: %+v vs %+v (item %q)", d.Slot, prev, s, d.ID)
			}
		} else {
			seen[d.Slot] = s
		}
	}
	return nil
}

// LessonByID returns the lesson with the given id, or false.
func (c Catalog) LessonByID(id string) (Lesson, bool) {
	for _, l := range c.Lessons {
		if l.ID == id {
			return l, true
		}
	}
	return Lesson{}, false
}
