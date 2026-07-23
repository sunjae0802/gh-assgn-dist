# gh-assgn-dist

A [GitHub CLI](https://cli.github.com/) extension that's a minimal reproduction of GitHub Classroom: create a classroom, distribute an assignment as private per-student repos, and keep local clones up to date.

## Install

```bash
gh extension install sunjae0802/gh-assgn-dist
```

## Concepts

- **Classroom file** (`CLASSROOM.yaml`) — class name, GitHub org, roster path, and the list of assignments created so far.
- **Roster file** — a CSV with columns `name,email,github`. Withdrawn students can stay in the file; they just won't get new repos.
- **Student repos** are named `CLASSROOM-ASSGN-USERNAME`, e.g. `github.com/witcomp1000/witcomp1000-fall26-a1-alice`.

## Usage

### Create a classroom

```bash
gh assgn-dist new --org ORG --roster roster.csv CLASSROOM
```

Verifies the org exists and that you have admin access, then writes `CLASSROOM.yaml`.

### Distribute an assignment

```bash
gh assgn-dist create --template OWNER/REPO --classroom CLASSROOM.yaml ASSGN
```

Verifies the template repo exists, then for each student in the roster creates a private repo from the template and adds the student as an outside collaborator with `write` access.

- `--template` defaults to the assignment short name (`ASSGN`) if omitted.
- `--classroom` defaults to the single `.yaml` file in the current directory if omitted.
- `--dry-run` prints the `gh api` commands instead of running them.

### Clone or update student repos

```bash
gh assgn-dist student-repos ASSGN
```

Clones any repo that doesn't exist locally yet; runs `git pull` on any that do.

## How it works

All GitHub operations go through `gh api` (see the [reference](https://cli.github.com/manual/gh_api)). No GitHub token handling or HTTP client setup is needed beyond what `gh` already provides.

## Development

```bash
go build ./...
```

Project layout:

```
main.go                  # cobra root command
cmd/new.go               # `new` subcommand
cmd/create.go            # `create` subcommand
cmd/student_repos.go     # `student-repos` subcommand
internal/classroom.go    # CLASSROOM.yaml read/write
internal/roster.go       # CSV roster parsing
```
