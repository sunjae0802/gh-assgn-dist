# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`assgn-dist` is GitHub CLI extension that is a minimal reproduction of GitHub Classroom CLI.

## Key Domain Concepts

### Files

- **Classroom file**: `classroom.yaml` by default (override with `--classroom`) — stores classroom name, GitHub org, and roster filename; one file per classroom
- **Roster file**: CSV with columns: name, email, github username

### Repo naming convention

Student repos follow: `github.com/ORG/CLASSROOM-ASSGN-USERNAME`

The classroom name is the prefix that keeps sections apart inside one flat org, so the same
assignment name can be reused across classes and terms without colliding.

Example: `github.com/witcomp1000/witcomp1000-fall26-a1-alice`

## CLI Subcommands

```bash
# Create a new classroom (verifies org exists and user has admin access)
assgn-dist create --name CLASSROOM --org ORG --roster CSVFILE

# Distribute an assignment (verifies template repo exists, creates student repos, adds students as outside collaborators with write role)
assgn-dist distribute --template REPO --classroom CLASSROOM.yaml ASSGN

# Clone or update student repos for an assignment (into ASSGN/CLASSROOM-ASSGN-USERNAME)
assgn-dist clone ASSGN
```

## GitHub API Usage

All GitHub operations go through `gh api`. Reference: https://cli.github.com/manual/gh_api

Key behaviors:
- `create`: verify org exists + user has admin access via `gh api`
- `distribute`: verify template repo exists, then create per-student repos and add each student as outside collaborator with `write` role; safe to re-run — skips repos that already exist, still (re-)adds the collaborator
- `clone`: clone repos if missing, pull if they already exist

## Notes

- Withdrawn students remain in the roster file
- The `--classroom` flag is optional for `create`, `distribute`, and `clone`; defaults to `classroom.yaml`. `create` writes it, the others read it
- `create` takes no positional argument; the classroom name comes from the required `--name` flag
- The `--template` flag defaults to `ORG/ASSGN` (classroom org + assignment short name) if omitted
- The `--only` flag on `distribute` takes a comma-separated list of GitHub usernames and skips every
  other student; matching is case-insensitive, and omitting it distributes to the whole roster
- The `--dry-run` flag makes no changes: `distribute` prints the `gh api` commands it would run,
  `clone` prints the `gh repo clone` / `git pull` commands, and `create` verifies the org then
  prints what it would write. Read-only verification still runs in all three, so a dry run needs
  working `gh` auth
