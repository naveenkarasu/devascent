package guiapi

// ── Primers (the "Learn" surface) ─────────────────────────────────────────────
// 15 categories × 8 languages of authored refreshers already live in the
// content catalog; this exposes one for the bench's Learn drawer.

// PrimerOpView is one labeled snippet.
type PrimerOpView struct {
	Label string `json:"label"`
	Code  string `json:"code"`
}

// PrimerSectionView is one grouped page of ops.
type PrimerSectionView struct {
	Title string         `json:"title"`
	Ops   []PrimerOpView `json:"ops"`
}

// PrimerView is a category primer rendered for one language.
type PrimerView struct {
	Found    bool                `json:"found"`
	Category string              `json:"category"`
	Title    string              `json:"title"`
	Summary  string              `json:"summary"`
	Sections []PrimerSectionView `json:"sections"`
	Example  string              `json:"example"`
}

// PrimerFor returns the primer for a bench category in lang (python fallback
// handled by the catalog).
func (e *Engine) PrimerFor(category, lang string) PrimerView {
	p, ok := e.cat.PrimerByCategoryAndLang(category, lang)
	if !ok {
		return PrimerView{Found: false, Category: category}
	}
	v := PrimerView{
		Found: true, Category: p.Category, Title: p.Title,
		Summary: p.Summary, Example: p.Example,
	}
	for _, s := range p.Sections {
		sv := PrimerSectionView{Title: s.Title}
		for _, op := range s.Ops {
			sv.Ops = append(sv.Ops, PrimerOpView{Label: op.Label, Code: op.Code})
		}
		v.Sections = append(v.Sections, sv)
	}
	return v
}
