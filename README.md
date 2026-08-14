# gh-assgn-dist

A [GitHub CLI](https://cli.github.com/) extension that's a minimal reproduction of GitHub Classroom:
create a classroom, distribute an assignment as private per-student repos, and keep local clones up
to date.

## How it works

All GitHub operations go through `gh api` (see the
[reference](https://cli.github.com/manual/gh_api)). No GitHub token handling or HTTP client setup is
needed beyond what `gh` already provides.

## Install

```bash
gh extension install sunjae0802/gh-assgn-dist
```

## Concepts

- **Roster file** — a CSV with columns `name,email,github`. If a student is added later, just add a
  new row; if a student withdraws, just delete the row
- **Classroom name** — The name of the classroom is used as a prefix. An example would be
  `witcomp1000-fall26`
- **Classroom file** (`classroom.yaml` by default; override with `--classroom`) — created by this
  extension, and contains class name, GitHub org, and path to roster file.
- **Student repos** are named `github.com/ORG/CLASSROOM-ASSGN-USERNAME`. For example, if the ORG is
  `witcomp1000`, and classroom name is `witcomp1000-fall26`, assignment name is `a1`, and username
  is `alice`, then the created repo is `github.com/witcomp1000/witcomp1000-fall26-a1-alice`.

## Usage

### Create a classroom

Verifies the org exists and that you have admin access, then writes the classroom file.

- `--classroom` sets the output path; defaults to `classroom.yaml`.

```bash
$ gh assgn-dist new --org ORG --roster roster.csv CLASSROOM
# Example:
$ gh assgn-dist new --org witcomp1000 --roster roster.csv witcomp1000-fall26
```

### Distribute an assignment

Verifies the template repo exists, then for each student in the roster creates a private repo from
the template and adds the student as an outside collaborator with `write` access.

Safe to re-run: if a student's repo already exists, creation is skipped, but the collaborator is
still (re-)added, so re-running also repairs any invite that failed on a prior run and picks up
students who don't have a repo yet.

- `--template` defaults to the assignment short name (`ASSGN`) if omitted.
- `--classroom` defaults to `classroom.yaml` if omitted.
- `--dry-run` prints the `gh api` commands instead of running them.

```bash
$ gh assgn-dist create --template OWNER/REPO ASSGN
# Example:
$ gh assgn-dist create --template sunjae0802/cs1-a1 a1
```

### Clone or update student repos

Clones any repo that doesn't exist locally yet (into `ASSGN/CLASSROOM-ASSGN-USERNAME`); runs
`git pull` on any that do.

```bash
$ gh assgn-dist clone ASSGN
# Example:
$ gh assgn-dist clone a1
```

## Development

```bash
go build ./...
```

Project layout:

```
main.go                  # cobra root command
cmd/new.go               # `new` subcommand
cmd/create.go            # `create` subcommand
cmd/clone.go             # `clone` subcommand
internal/classroom.go    # classroom.yaml read/write
internal/roster.go       # CSV roster parsing
```
