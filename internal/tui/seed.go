package tui

import "devascent/internal/ticket"

// seedSprint1 is the first real apprenticeship sprint (#62, re-authored for the
// Company Sim): you've just joined **Pixel Forge** — a SaaS where indie studios
// upload, optimize, and serve game assets via an API + CDN (Python/FastAPI). This
// is the Graduate stage (Solo→MVP): onboarding chores + small, fully-deterministic
// fixes on the real product. The day starts at Day 0, everything is in To Do, and
// tickets are revealed on the day they're assigned (scheduled reveal).
func seedSprint1() (*ticket.Project, *ticket.Sprint) {
	proj := &ticket.Project{Key: "PXF", Name: "Pixel Forge"}
	sp := &ticket.Sprint{
		Number: 1, Day: 0, Capacity: 12, Goal: "Onboarding & first fixes",
		Tickets: []*ticket.Ticket{
			// ── Onboarding (ungraded chores; complete with [d] mark done) ──
			{Key: "PXF-100", Type: ticket.Task, Title: "Set up your local dev environment", Status: ticket.ToDo,
				Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Sam (manager)",
				Labels: []string{"onboarding"}, AssignedDay: 0,
				Desc: "Get Pixel Forge running locally so you can ship. Clone the repo, create a virtualenv, install deps, copy `.env.example` → `.env`, bring up Postgres + Redis with `docker compose up`, run the migrations, and confirm the test suite is green.",
				Subtasks: []ticket.Subtask{
					{Title: "git clone + python -m venv .venv"},
					{Title: "pip install -r requirements.txt"},
					{Title: "cp .env.example .env"},
					{Title: "docker compose up -d  (Postgres + Redis)"},
					{Title: "alembic upgrade head"},
					{Title: "pytest  → all green"},
				},
				Comments: []ticket.Comment{{Author: "Sam (manager)", Body: "Welcome aboard! Knock this out first — ping me if Docker gives you trouble."}},
			},
			{Key: "PXF-110", Type: ticket.Task, Title: "Read the README & architecture brief", Status: ticket.ToDo,
				Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Sam (manager)",
				Labels: []string{"onboarding", "docs"}, AssignedDay: 0,
				Desc:     "Pixel Forge is a SaaS for indie game studios to upload, optimize, version, and serve art/sprite assets via an API + CDN. Stack (MVP today): FastAPI backend, React dashboard, Postgres (+Alembic), Redis + a worker for image jobs, S3 for assets. We work the board: pick a ticket → fix it → tests pass → review → merge. Skim the README and the `/docs/architecture.md` so the asset/upload/billing modules make sense.",
				Learn:    "# Pixel Forge modules\n# api/        FastAPI routers (assets, auth, billing)\n# workers/    image optimize + thumbnail jobs (Redis queue)\n# models/     SQLModel + Alembic migrations\n# web/        React dashboard",
				Comments: []ticket.Comment{{Author: "Sam (manager)", Body: "The asset pipeline is the heart of it — that's where most of your tickets will land."}},
			},

			// ── First fixes (graded; the proven Python fixtures, re-themed) ──
			{Key: "PXF-201", Type: ticket.Bug, Title: "Asset list API drops the first item per page", Status: ticket.ToDo,
				Priority: ticket.PMajor, Points: 2, Assignee: "you", Reporter: "Priya (QA)",
				Labels: []string{"backend", "api"}, AssignedDay: 0,
				Desc:       "The studio dashboard's asset list is paginated, but page 1 starts at index `size` instead of 0 — every page loses its first asset. Pages are 1-indexed; fix the slice without breaking the last (partial) page.",
				Acceptance: crit("page 1 returns assets 1..size", "no assets overlap or go missing between pages", "the hidden tests pass"),
				Learn:      "def paginate(items, page, size):\n    start = (page - 1) * size   # 1-indexed pages\n    return items[start:start + size]",
				Comments:   []ticket.Comment{{Author: "Priya (QA)", Body: "Repro: open the assets tab with >1 page; the first asset on each page is gone."}},
				Grading: pyGrading("paginate.py",
					"def paginate(items, page, size):\n    start = page * size  # BUG: page 1 should start at index 0\n    return items[start:start + size]\n",
					"from paginate import paginate\n\nitems = list(range(1, 11))\nassert paginate(items, 1, 3) == [1, 2, 3], paginate(items, 1, 3)\nassert paginate(items, 2, 3) == [4, 5, 6], paginate(items, 2, 3)\nassert paginate(items, 4, 3) == [10], paginate(items, 4, 3)\nprint('OK')\n"),
			},
			{Key: "PXF-202", Type: ticket.Bug, Title: "Studio signup rejects valid + tag emails", Status: ticket.ToDo,
				Priority: ticket.PMajor, Points: 2, Assignee: "you", Reporter: "Sam (manager)",
				Labels: []string{"backend", "auth"}, AssignedDay: 1,
				Desc:       "A studio tried to sign up with `team+billing@studio.com` (plus-tagging) and got rejected. Accept plus-tags while still rejecting genuinely malformed addresses.",
				Acceptance: crit("team+billing@studio.com is accepted", "bad@ and @studio.com are rejected", "the hidden tests pass"),
				Grading: pyGrading("validate.py",
					"def is_valid(email):\n    if '+' in email:        # BUG: plus-tags are valid\n        return False\n    parts = email.split('@')\n    if len(parts) != 2:\n        return False\n    local, domain = parts\n    return bool(local) and '.' in domain\n",
					"from validate import is_valid\n\nassert is_valid('a+b@x.com')\nassert is_valid('user@mail.com')\nassert not is_valid('bad@')\nassert not is_valid('@x.com')\nassert not is_valid('no-at-sign')\nprint('OK')\n"),
			},
			{Key: "PXF-203", Type: ticket.Task, Title: "Slugify project names for asset URLs", Status: ticket.ToDo,
				Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Dev (frontend)",
				Labels: []string{"backend"}, AssignedDay: 1,
				Desc:       "Asset URLs use the project name as a slug. Implement `slugify`: lowercase, runs of non-alphanumeric → a single '-', trimmed.",
				Acceptance: crit("'Hello, World!' → 'hello-world'", "repeated separators collapse to one '-'", "the hidden tests pass"),
				Learn:      "# str.isalnum() + str.lower(); track whether the previous char was already a '-'.",
				Grading: pyGrading("slug.py",
					"def slugify(title):\n    return title  # TODO: implement\n",
					"from slug import slugify\n\nassert slugify('Hello, World!') == 'hello-world', slugify('Hello, World!')\nassert slugify('  A  B  ') == 'a-b', slugify('  A  B  ')\nassert slugify('Already-slug') == 'already-slug'\nprint('OK')\n"),
			},
			{Key: "PXF-204", Type: ticket.Bug, Title: "Storage usage total drops the last asset", Status: ticket.ToDo,
				Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Priya (QA)",
				Labels: []string{"backend", "billing"}, AssignedDay: 2,
				Desc:       "The per-studio storage total (used for billing) is short by one asset — the loop stops one early. Fix the range so every asset is counted.",
				Acceptance: crit("total([1,2,3]) == 6", "total([]) == 0 and total([5]) == 5", "the hidden tests pass"),
				Grading: pyGrading("total.py",
					"def total(nums):\n    s = 0\n    for i in range(len(nums) - 1):  # BUG: skips the last element\n        s += nums[i]\n    return s\n",
					"from total import total\n\nassert total([1, 2, 3]) == 6, total([1, 2, 3])\nassert total([]) == 0\nassert total([5]) == 5\nprint('OK')\n"),
			},
			{Key: "PXF-205", Type: ticket.Story, Title: "Avg asset size crashes for empty projects", Status: ticket.ToDo,
				Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Sam (manager)",
				Labels: []string{"backend"}, AssignedDay: 2,
				Desc:       "The analytics panel calls `average([])` for a brand-new project and 500s (ZeroDivisionError). Return 0.0 for an empty project; otherwise the mean asset size.",
				Acceptance: crit("average([]) == 0.0 (no crash)", "average([2,4]) == 3.0", "the hidden tests pass"),
				Grading: pyGrading("avg.py",
					"def average(nums):\n    return sum(nums) / len(nums)  # BUG: crashes on []\n",
					"from avg import average\n\nassert average([2, 4]) == 3.0\nassert average([5]) == 5.0\nassert average([]) == 0.0\nprint('OK')\n"),
			},
			{Key: "PXF-206", Type: ticket.TechDebt, Title: "Dedupe asset tags but keep order", Status: ticket.ToDo,
				Priority: ticket.PTrivial, Points: 1, Assignee: "you", Reporter: "Dev (frontend)",
				Labels: []string{"backend", "tech-debt"}, AssignedDay: 3,
				Desc:       "Asset tag de-duplication uses a set, so it drops duplicates but scrambles tag order in the UI. Keep the first occurrence of each tag in its original position.",
				Acceptance: crit("dedupe([3,1,3,2,1]) == [3,1,2]", "dedupe([]) == []", "the hidden tests pass"),
				Grading: pyGrading("dedupe.py",
					"def dedupe(items):\n    return list(set(items))  # BUG: loses order\n",
					"from dedupe import dedupe\n\nassert dedupe([3, 1, 3, 2, 1]) == [3, 1, 2], dedupe([3, 1, 3, 2, 1])\nassert dedupe([]) == []\nprint('OK')\n"),
			},

			// ── Backlog (the Pixel Forge product roadmap; Backlog view, grouped by epic) ──
			{Key: "PXF-300", Type: ticket.Story, Title: "Bulk asset upload (zip)", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Sam (manager)", Labels: []string{"backend", "api"}, Epic: "Asset Pipeline"},
			{Key: "PXF-301", Type: ticket.Story, Title: "Asset versioning & rollback", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Priya (QA)", Labels: []string{"backend"}, Epic: "Asset Pipeline"},
			{Key: "PXF-310", Type: ticket.Feature, Title: "OAuth (GitHub/Google) login", Status: ticket.Backlog, Priority: ticket.PMajor, Reporter: "Sam (manager)", Labels: []string{"auth"}, Epic: "Auth & Accounts"},
			{Key: "PXF-311", Type: ticket.Story, Title: "Team invites & roles", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Sam (manager)", Labels: []string{"auth"}, Epic: "Auth & Accounts"},
			{Key: "PXF-320", Type: ticket.Feature, Title: "Stripe billing integration", Status: ticket.Backlog, Priority: ticket.PMajor, Reporter: "Sam (manager)", Labels: []string{"billing"}, Epic: "Billing"},
			{Key: "PXF-321", Type: ticket.Story, Title: "Usage-based metering & invoices", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Sam (manager)", Labels: []string{"billing"}, Epic: "Billing"},
			{Key: "PXF-330", Type: ticket.Feature, Title: "Serve assets via CloudFront CDN", Status: ticket.Backlog, Priority: ticket.PMajor, Reporter: "Dev (frontend)", Labels: []string{"infra", "cdn"}, Epic: "CDN & Performance"},
			{Key: "PXF-340", Type: ticket.TechDebt, Title: "Move image jobs to a real queue", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Dev (frontend)", Labels: []string{"infra"}, Epic: "Reliability"},
			{Key: "PXF-341", Type: ticket.Story, Title: "Job retries + dead-letter queue", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Priya (QA)", Labels: []string{"infra"}, Epic: "Reliability"},
			{Key: "PXF-350", Type: ticket.Feature, Title: "Read replicas + connection pooling", Status: ticket.Backlog, Priority: ticket.PMinor, Reporter: "Sam (manager)", Labels: []string{"data", "scale"}, Epic: "Scale"},
			{Key: "PXF-360", Type: ticket.Story, Title: "RBAC + audit log", Status: ticket.Backlog, Priority: ticket.PMajor, Reporter: "Sam (manager)", Labels: []string{"security"}, Epic: "Security & Compliance"},
		},
	}
	// Derive each ticket's SLA due day from its priority + assigned day (real-world
	// SLA mapping; see the Company Sim design doc).
	// Map each sprint ticket to its product epic (backlog items carry their own).
	epics := map[string]string{
		"PXF-100": "Foundations", "PXF-110": "Foundations",
		"PXF-201": "Asset Pipeline", "PXF-203": "Asset Pipeline", "PXF-206": "Asset Pipeline",
		"PXF-202": "Auth & Accounts", "PXF-204": "Billing", "PXF-205": "Reliability",
	}
	for _, t := range sp.Tickets {
		if t.Epic == "" {
			t.Epic = epics[t.Key]
		}
		if t.Status == ticket.Backlog {
			continue // backlog items aren't scheduled yet
		}
		if t.DueDay == 0 {
			t.DueDay = ticket.DueDayFor(t.Priority, t.AssignedDay)
		}
		t.CreatedDay = t.AssignedDay
	}
	return proj, sp
}

// pyGrading builds a Python repo-grade spec: the player edits editFile (seeded
// with start), graded by `python check.py` (ticket-owned).
func pyGrading(editFile, start, check string) *ticket.Grading {
	return &ticket.Grading{
		Lang:        "python",
		Command:     []string{"python", "check.py"},
		EditPaths:   []string{editFile},
		StartFiles:  map[string]string{editFile: start},
		HiddenFiles: map[string]string{"check.py": check},
	}
}

// crit turns plain strings into acceptance criteria.
func crit(texts ...string) []ticket.Criterion {
	out := make([]ticket.Criterion, len(texts))
	for i, t := range texts {
		out[i] = ticket.Criterion{Text: t}
	}
	return out
}
