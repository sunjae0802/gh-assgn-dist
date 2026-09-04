# gh-assgn-dist

A [GitHub CLI](https://cli.github.com/) extension that's a minimal reproduction of GitHub Classroom:
create a classroom, distribute an assignment as private per-student repos, and keep local clones up
to date. It does NOT include an autograder; I use Gradescope for this.

This extension assumes that you have a *private GitHub organization*, and you want to create private
repositories for each student. It uses a template repository as the base, makes N private copies of
the repository, and adds each student as an "external collaborator" with *write* permissions. This
ensures students will not be able to view each other's repo, but you as the instructor and the
student can.

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

    $ gh assgn-dist create --org ORG --roster roster.csv CLASSROOM

Arguments:

- `--classroom` sets the output path; defaults to `classroom.yaml`.
- `--dry-run` still verifies the org and your admin access, then prints what would be written
  instead of writing the classroom file.

### Distribute an assignment

Verifies the template repo exists, then for each student in the roster creates a private repo from
the template and adds the student as an outside collaborator with `write` access.

**Safe to re-run**: if a student's repo already exists, creation is skipped, but the collaborator is
still (re-)added, so re-running also repairs any invite that failed on a prior run and picks up
students who don't have a repo yet.

    $ gh assgn-dist distribute --template OWNER/REPO ASSGN

Arguments:

- `--template` defaults to `ORG/ASSGN` (the classroom org and the assignment short name) if omitted.
- `--classroom` defaults to `classroom.yaml` if omitted.
- `--only` limits distribution to a comma-separated list of GitHub usernames (for example
  `--only alice,bob`); everyone else in the roster is skipped. Matching is case-insensitive, and a
  username that isn't in the roster is reported as a warning. Omit it to distribute to everyone.
- `--dry-run` prints the `gh api` commands instead of running them.

### Clone or update student repos

Clones any repo that doesn't exist locally yet (into `ASSGN/CLASSROOM-ASSGN-USERNAME`); runs
`git pull` on any that do.

    $ gh assgn-dist clone ASSGN

Arguments:

- `--classroom` defaults to `classroom.yaml` if omitted.
- `--dry-run` prints the `gh repo clone` / `git pull` commands instead of running them.

## Example

```bash
# Example roster contents
$ cat roster.csv
name,email,github
Sunjae Park,sunjae@email.com,sunjaeatwit
John Doe,jdoe@email.com,jdoe11atwit
Jane Austen,jausten@email.com,jaustenatwit

# Create a new classroom under github.com/witcomp1000
$ gh assgn-dist create --org witcomp1000 --roster roster.csv witcomp1000-fall26

# Distribute a1 using github.com/sunjae0802/cs1-a1 as template repo
$ gh assgn-dist distribute --template sunjae0802/cs1-a1 a1

# Distribute a1 to just one student (for example, a late add)
$ gh assgn-dist distribute --template sunjae0802/cs1-a1 --only leopardatwit a1

# Clone a1
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
