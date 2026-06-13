# Security Policy

## Threat model — please read

DevAscent **compiles and runs code using your own installed toolchain** (Python,
Go, Rust, etc.). It is **not a sandbox.** The protections around grading are for
*accidents* (a runaway loop, a crash) — a timeout, process-tree termination, and a
scoped temporary working directory — **not** for containing malicious code.

Practical implications:

- The code you write and grade is your own, so this is normally fine — it's no
  different from running your own program.
- **Do not** paste in and grade solution code from an untrusted source; it would
  run with your user's privileges, just like any script you execute yourself.

If you need strong isolation, run DevAscent inside a container or VM.

## Supported versions

DevAscent is pre-1.0 and moves fast. Only the latest `main` is supported; please
reproduce issues against the current build before reporting.

## Reporting a vulnerability

Please **do not** open a public issue for a security problem.

Use GitHub's private reporting: the repository's **Security → "Report a
vulnerability"** tab (GitHub Private Vulnerability Reporting). If that's
unavailable, contact the maintainer via their GitHub profile
(<https://github.com/naveenkarasu>).

Please include repro steps, your OS, and the affected version/commit. I'll
acknowledge as soon as I can — this is a solo project, so response times are
best-effort.
