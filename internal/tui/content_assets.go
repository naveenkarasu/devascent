package tui

import "devascent/internal/ticket"

// assetPipelineTickets is the first gradeable epic's content — small, deterministic
// Python fixes on Pixel Forge's asset upload/serve pipeline. Each is graded by
// GradeRepo (starter fails → the fix passes), like the seed tickets. They're
// appended to Sprint 1 and scheduled across its days (see seedSprint1).
func assetPipelineTickets() []*ticket.Ticket {
	const epic = "Asset Pipeline"
	return []*ticket.Ticket{
		{Key: "PXF-210", Type: ticket.Bug, Title: "Filename sanitizer allows path traversal", Status: ticket.ToDo,
			Priority: ticket.PCritical, Points: 2, Assignee: "you", Reporter: "Priya (QA)",
			Labels: []string{"backend", "security"}, Epic: epic, AssignedDay: 2,
			Desc:       "Uploaded filenames are stored as-is, so `../../etc/passwd` escapes the asset directory. Return just the safe basename (strip any path/traversal); empty → \"asset\".",
			Acceptance: crit("'../../etc/passwd' → 'passwd'", "'a/b/c.jpg' → 'c.jpg'", "empty → 'asset'"),
			Comments:   []ticket.Comment{{Author: "Priya (QA)", Body: "Security found this in a pen-test — please prioritise."}},
			Grading: pyGrading("names.py",
				"def safe_name(name):\n    return name  # BUG: no sanitization\n",
				"from names import safe_name\n\nassert safe_name('../../etc/passwd') == 'passwd', safe_name('../../etc/passwd')\nassert safe_name('a/b/c.jpg') == 'c.jpg'\nassert safe_name('logo.png') == 'logo.png'\nassert safe_name('') == 'asset'\nprint('OK')\n")},

		{Key: "PXF-211", Type: ticket.Task, Title: "Thumbnails must preserve aspect ratio", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 2, Assignee: "you", Reporter: "Dev (frontend)",
			Labels: []string{"backend", "images"}, Epic: epic, AssignedDay: 2,
			Desc:       "Thumbnails come out squashed. Implement `thumb_size(w, h, max_side)`: scale so the longest side is `max_side` (floor), preserving aspect ratio; if it already fits, return (w, h).",
			Acceptance: crit("1000x500 @200 → (200,100)", "500x1000 @200 → (100,200)", "100x100 @200 → (100,100)"),
			Learn:      "scale = max_side / max(w, h)\nnew = (int(w * scale), int(h * scale))",
			Grading: pyGrading("thumb.py",
				"def thumb_size(w, h, max_side):\n    return (max_side, max_side)  # TODO: keep aspect ratio\n",
				"from thumb import thumb_size\n\nassert thumb_size(1000, 500, 200) == (200, 100), thumb_size(1000, 500, 200)\nassert thumb_size(500, 1000, 200) == (100, 200)\nassert thumb_size(100, 100, 200) == (100, 100)\nprint('OK')\n")},

		{Key: "PXF-212", Type: ticket.Bug, Title: "Content-type missing for .webp/.svg uploads", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 1, Assignee: "you", Reporter: "Priya (QA)",
			Labels: []string{"backend", "api"}, Epic: epic, AssignedDay: 3,
			Desc:       "WebP and SVG assets are served as `application/octet-stream`, so browsers download them instead of displaying. Map those extensions to the correct MIME type.",
			Acceptance: crit("a.webp → image/webp", "a.svg → image/svg+xml", "unknown → application/octet-stream"),
			Grading: pyGrading("mime.py",
				"def content_type(filename):\n    ext = filename.rsplit('.', 1)[-1].lower()\n    m = {'png': 'image/png', 'jpg': 'image/jpeg', 'jpeg': 'image/jpeg', 'gif': 'image/gif'}\n    return m.get(ext, 'application/octet-stream')  # BUG: webp/svg missing\n",
				"from mime import content_type\n\nassert content_type('a.png') == 'image/png'\nassert content_type('a.webp') == 'image/webp', content_type('a.webp')\nassert content_type('a.svg') == 'image/svg+xml'\nassert content_type('a.bin') == 'application/octet-stream'\nprint('OK')\n")},

		{Key: "PXF-213", Type: ticket.Story, Title: "Human-readable asset sizes in the dashboard", Status: ticket.ToDo,
			Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Dev (frontend)",
			Labels: []string{"backend"}, Epic: epic, AssignedDay: 4,
			Desc:       "The dashboard shows raw byte counts. Implement `humansize(n)`: bytes → 'N B' / 'N.N KB' / 'N.N MB' / 'N.N GB' (1024-based, one decimal for KB and up).",
			Acceptance: crit("512 → '512 B'", "1536 → '1.5 KB'", "2*1024*1024 → '2.0 MB'"),
			Grading: pyGrading("humansize.py",
				"def humansize(n):\n    return f'{n} B'  # TODO: KB/MB/GB\n",
				"from humansize import humansize\n\nassert humansize(512) == '512 B'\nassert humansize(1536) == '1.5 KB', humansize(1536)\nassert humansize(2 * 1024 * 1024) == '2.0 MB'\nprint('OK')\n")},

		{Key: "PXF-214", Type: ticket.Bug, Title: "Same filename across studios collides in storage", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 1, Assignee: "you", Reporter: "Sam (manager)",
			Labels: []string{"backend", "storage"}, Epic: epic, AssignedDay: 5,
			Desc:       "Two studios uploading `logo.png` overwrite each other — the storage key isn't namespaced. Implement `asset_key(project, filename)` → `'<project>/<filename>'`.",
			Acceptance: crit("('acme','logo.png') → 'acme/logo.png'", "keys differ across projects for the same filename"),
			Grading: pyGrading("keys.py",
				"def asset_key(project, filename):\n    return filename  # BUG: not namespaced by project\n",
				"from keys import asset_key\n\nassert asset_key('acme', 'logo.png') == 'acme/logo.png'\nassert asset_key('beta', 'logo.png') == 'beta/logo.png'\nassert asset_key('acme', 'logo.png') != asset_key('beta', 'logo.png')\nprint('OK')\n")},

		{Key: "PXF-215", Type: ticket.TechDebt, Title: "Normalize asset tags (case/space/dupes)", Status: ticket.ToDo,
			Priority: ticket.PTrivial, Points: 1, Assignee: "you", Reporter: "Dev (frontend)",
			Labels: []string{"backend"}, Epic: epic, AssignedDay: 6,
			Desc:       "Asset tags are stored raw, so 'Art', ' art ' and '' all coexist. Implement `normalize(tags)`: trim, lowercase, drop empties, dedupe keeping first-seen order.",
			Acceptance: crit("['Art',' art ','3D',''] → ['art','3d']", "order preserved", "empty list → []"),
			Grading: pyGrading("tags.py",
				"def normalize(tags):\n    return [t.lower() for t in tags]  # BUG: no trim/dedupe/empty\n",
				"from tags import normalize\n\nassert normalize(['Art', ' art ', '3D', '']) == ['art', '3d'], normalize(['Art', ' art ', '3D', ''])\nassert normalize([]) == []\nprint('OK')\n")},
	}
}
