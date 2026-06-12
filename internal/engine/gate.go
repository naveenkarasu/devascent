package engine

// The graduation gate (Track A3): the design-spec apprenticeship exit, as
// opposed to the Step-0 starter milestone in bench.go. To graduate you need
// ≥49 fully banked Blind-75 problems, a minimum in every category that has
// Blind-75 coverage, and a handful of mandatory classics. "Fully banked"
// means solved AND write-up complete (A1) — the gate counts explained solves,
// not green checkmarks. Counting is by canonical slug so problem variants
// can't be double-counted.

import (
	"sort"

	"devascent/internal/content"
)

// Blind75Target is the design-spec solve count for graduation.
const Blind75Target = 49

// gateMandatory are the can't-skip classics (canonical slugs, all blind75).
var gateMandatory = []string{
	"two-sum",
	"valid-parentheses",
	"best-time-to-buy-and-sell-stock",
	"reverse-linked-list",
	"merge-intervals",
}

// gateCatRequired derives a category's minimum from its Blind-75 pool size:
// half the pool, floor, at least one. Self-maintaining as content evolves.
func gateCatRequired(available int) int {
	r := available / 2
	if r < 1 {
		r = 1
	}
	return r
}

// GateCategory is one category row of the progress view.
type GateCategory struct {
	Category  string
	Done      int // fully banked slugs
	Required  int
	Available int
}

// GateItem is one mandatory problem row.
type GateItem struct {
	Slug  string
	Title string
	Done  bool
}

// GateProgress is the full graduation-gate accounting for one save.
type GateProgress struct {
	Full        int // fully banked Blind-75 slugs (solved + write-up)
	Provisional int // solved but write-up pending
	Target      int
	Categories  []GateCategory
	Mandatory   []GateItem
	CountMet    bool
	CatsMet     bool
	MandatoryOK bool
	Met         bool
}

// Blind75Progress computes the gate from the catalog and the save's solved
// set; fullyBanked reports whether a solved problem's write-up is complete.
func Blind75Progress(problems []content.Problem, solved map[string]bool, fullyBanked func(id string) bool) GateProgress {
	type slugState struct {
		category    string
		title       string
		full        bool
		provisional bool
	}
	slugs := map[string]*slugState{}
	for _, p := range problems {
		inB75 := false
		for _, l := range p.Lists {
			if l == "blind75" {
				inB75 = true
				break
			}
		}
		if !inB75 {
			continue
		}
		slug := p.CanonicalSlug()
		s := slugs[slug]
		if s == nil {
			s = &slugState{category: p.Category, title: p.Title}
			slugs[slug] = s
		}
		if solved[p.ID] {
			if fullyBanked(p.ID) {
				s.full = true
			} else {
				s.provisional = true
			}
		}
	}

	g := GateProgress{Target: Blind75Target}
	catAvail := map[string]int{}
	catDone := map[string]int{}
	for _, s := range slugs {
		catAvail[s.category]++
		if s.full {
			g.Full++
			catDone[s.category]++
		} else if s.provisional {
			g.Provisional++
		}
	}

	cats := make([]string, 0, len(catAvail))
	for c := range catAvail {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	g.CatsMet = true
	for _, c := range cats {
		row := GateCategory{Category: c, Done: catDone[c], Required: gateCatRequired(catAvail[c]), Available: catAvail[c]}
		if row.Done < row.Required {
			g.CatsMet = false
		}
		g.Categories = append(g.Categories, row)
	}

	g.MandatoryOK = true
	for _, slug := range gateMandatory {
		s := slugs[slug]
		item := GateItem{Slug: slug}
		if s != nil {
			item.Title = s.title
			item.Done = s.full
		}
		if !item.Done {
			g.MandatoryOK = false
		}
		g.Mandatory = append(g.Mandatory, item)
	}

	g.CountMet = g.Full >= g.Target
	g.Met = g.CountMet && g.CatsMet && g.MandatoryOK
	return g
}
