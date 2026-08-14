# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`assgn-dist` is GitHub CLI extension that is a minimal reproduction of GitHub Classroom CLI.

## Key Domain Concepts

### Files

- **Classroom file**: `classroom.yaml` by default (override with `--classroom`) — stores class name, GitHub org, and roster filename
- **Roster file**: CSV with columns: name, email, github username

### Repo naming convention

Student repos follow: `github.com/ORG/CLASSROOM-ASSGN-USERNAME`

Example: `github.com/witcomp1000/witcomp1000-fall26-a1-alice`

## CLI Subcommands

```bash
# Create a new classroom (verifies org exists and user has admin access)
assgn-dist new --org ORG --roster CSVFILE CLASSROOM

# Create an assignment (verifies template repo exists, creates student repos, adds students as outside collaborators with write role)
assgn-dist create --template REPO --classroom CLASSROOM ASSGN

# Clone or update student repos for an assignment
assgn-dist student-repos ASSGN
```

## GitHub API Usage

All GitHub operations go through `gh api`. Reference: https://cli.github.com/manual/gh_api

Key behaviors:
- `new`: verify org exists + user has admin access via `gh api`
- `create`: verify template repo exists, then create per-student repos and add each student as outside collaborator with `write` role
- `student-repos`: clone repos if missing, pull if they already exist

## Notes

- Withdrawn students remain in the roster file
- The `--classroom` flag is optional for `new`, `create`, and `student-repos`; defaults to `classroom.yaml`
- The `--template` flag defaults to the assignment short name if omitted
- The `--dry-run` flag prints the `gh api` commands that will be used
