package tui

import "devascent/internal/ticket"

// authAccountsTickets is the Auth & Accounts epic's gradeable content — small,
// deterministic Python fixes around login, sessions, lockout, and RBAC. Each is
// graded by GradeRepo (starter fails → fix passes). Scheduled into Sprint 1 after
// the Asset Pipeline epic (later days).
func authAccountsTickets() []*ticket.Ticket {
	const epic = "Auth & Accounts"
	return []*ticket.Ticket{
		{Key: "PXF-220", Type: ticket.Bug, Title: "Login is case-sensitive on email", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 1, Assignee: "you", Reporter: "Priya (QA)",
			Labels: []string{"backend", "auth"}, Epic: epic, AssignedDay: 7,
			Desc:       "A studio signed up as `Alice@Studio.com` but can't log in as `alice@studio.com`. Normalize emails (trim + lowercase) so login is case-insensitive.",
			Acceptance: crit("'Alice@X.COM ' → 'alice@x.com'", "already-normal emails unchanged"),
			Grading: pyGrading("login.py",
				"def normalize_email(e):\n    return e  # BUG: case-sensitive / untrimmed\n",
				"from login import normalize_email\n\nassert normalize_email('Alice@X.COM ') == 'alice@x.com', normalize_email('Alice@X.COM ')\nassert normalize_email('bob@y.com') == 'bob@y.com'\nprint('OK')\n")},

		{Key: "PXF-221", Type: ticket.Story, Title: "Reject weak passwords at signup", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 2, Assignee: "you", Reporter: "Sam (manager)",
			Labels: []string{"backend", "auth"}, Epic: epic, AssignedDay: 7,
			Desc:       "Signup only checks password length. Require at least 8 chars AND at least one letter AND one digit.",
			Acceptance: crit("'abc12345' accepted", "'12345678' and 'allletters' rejected", "'short1' rejected"),
			Grading: pyGrading("passwords.py",
				"def is_strong(pw):\n    return len(pw) >= 8  # BUG: only checks length\n",
				"from passwords import is_strong\n\nassert is_strong('abc12345')\nassert not is_strong('short1')\nassert not is_strong('allletters')\nassert not is_strong('12345678')\nprint('OK')\n")},

		{Key: "PXF-222", Type: ticket.Bug, Title: "Expired sessions are still accepted", Status: ticket.ToDo,
			Priority: ticket.PCritical, Points: 2, Assignee: "you", Reporter: "Priya (QA)",
			Labels: []string{"backend", "security"}, Epic: epic, AssignedDay: 8,
			Desc:       "A session token exactly at its expiry time is treated as valid (off-by-one). A token is valid only while `now < exp`.",
			Acceptance: crit("now<exp valid", "now==exp NOT valid", "now>exp NOT valid"),
			Grading: pyGrading("session.py",
				"def is_valid(exp, now):\n    return now <= exp  # BUG: accepts exactly-expired\n",
				"from session import is_valid\n\nassert is_valid(100, 50)\nassert not is_valid(100, 100)\nassert not is_valid(100, 150)\nprint('OK')\n")},

		{Key: "PXF-223", Type: ticket.Bug, Title: "Account lockout triggers one attempt too late", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 1, Assignee: "you", Reporter: "Sam (manager)",
			Labels: []string{"backend", "security"}, Epic: epic, AssignedDay: 9,
			Desc:       "Brute-force protection locks the account after `failed > limit`, so it allows one attempt past the limit. Lock at `failed >= limit`.",
			Acceptance: crit("2 of 3 → not locked", "3 of 3 → locked", "4 of 3 → locked"),
			Grading: pyGrading("lockout.py",
				"def locked(failed, limit):\n    return failed > limit  # BUG: off-by-one\n",
				"from lockout import locked\n\nassert not locked(2, 3)\nassert locked(3, 3)\nassert locked(4, 3)\nprint('OK')\n")},

		{Key: "PXF-224", Type: ticket.Story, Title: "RBAC: viewers must not delete assets", Status: ticket.ToDo,
			Priority: ticket.PMajor, Points: 2, Assignee: "you", Reporter: "Sam (manager)",
			Labels: []string{"backend", "security"}, Epic: epic, AssignedDay: 10,
			Desc:       "Permission checks only verify the role exists, not the action — so a viewer can delete. Implement `can(role, action)`: owner = read/write/delete, editor = read/write, viewer = read.",
			Acceptance: crit("owner can delete", "editor can write but not delete", "viewer can read but not write"),
			Learn:      "perms = {'owner': {'read','write','delete'}, 'editor': {'read','write'}, 'viewer': {'read'}}\nreturn action in perms.get(role, set())",
			Grading: pyGrading("rbac.py",
				"def can(role, action):\n    return role in ('owner', 'editor', 'viewer')  # BUG: ignores the action\n",
				"from rbac import can\n\nassert can('owner', 'delete')\nassert can('editor', 'write')\nassert not can('editor', 'delete')\nassert not can('viewer', 'write')\nassert can('viewer', 'read')\nprint('OK')\n")},

		{Key: "PXF-225", Type: ticket.Task, Title: "Mask emails in audit logs", Status: ticket.ToDo,
			Priority: ticket.PMinor, Points: 1, Assignee: "you", Reporter: "Priya (QA)",
			Labels: []string{"backend", "security"}, Epic: epic, AssignedDay: 11,
			Desc:       "Audit logs print full email addresses. Implement `mask(email)` → first char of the local part + `***` + `@domain` (e.g. `alice@example.com` → `a***@example.com`).",
			Acceptance: crit("'alice@example.com' → 'a***@example.com'", "'bob@x.io' → 'b***@x.io'"),
			Grading: pyGrading("mask.py",
				"def mask(email):\n    return email  # TODO: mask the local part\n",
				"from mask import mask\n\nassert mask('alice@example.com') == 'a***@example.com', mask('alice@example.com')\nassert mask('bob@x.io') == 'b***@x.io'\nprint('OK')\n")},
	}
}
