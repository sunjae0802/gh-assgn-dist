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

GitHub API calls go through `gh`'s own REST client (see the
[`gh api` reference](https://cli.github.com/manual/gh_api)); cloning shells out to `gh repo clone`
and `git pull`. Either way, no GitHub token handling or HTTP client setup is needed beyond what `gh`
already provides.

## Install

```bash
gh extension install sunjae0802/gh-assgn-dist
```

## Concepts

- **Roster file** — a CSV whose header names a `github` column (`GitHub` and `Github` work too);
  `name` and `email` are used when present, and columns may appear in any order. A roster with no
  GitHub column is an error. Blank rows are skipped, and a row that names someone but leaves the
  GitHub cell empty is skipped with a warning naming the line. If a student is added later, just add
  a new row; if a student withdraws, just delete the row
- **Classroom name** — set once with `--name` and stored in the classroom file; it is used as the
  prefix on every student repo this classroom creates. An example would be `witcomp1000-fall26`
- **Classroom file** (`classroom.yaml` by default; override with `--classroom`) — created by this
  extension, and contains the classroom name, GitHub org, and path to roster file. One file per
  classroom, so a second class is a second file (for example `--classroom cs2.yaml`).
- **Student repos** are named `github.com/ORG/CLASSROOM-ASSGN-USERNAME`. For example, if the ORG is
  `witcomp1000`, and classroom name is `witcomp1000-fall26`, assignment name is `a1`, and username
  is `alice`, then the created repo is `github.com/witcomp1000/witcomp1000-fall26-a1-alice`.

### Why the classroom name is a prefix

A GitHub org is flat: every repo in it shares one namespace, and this tool creates one repo per
student per assignment. Without a prefix, `a1` in your fall CS1 section and `a1` in your spring CS2
section would both want `ORG/a1-alice` — the second `distribute` would quietly find the repo already
there, skip creating it, and hand the student last term's starter code.

The classroom name is what keeps those apart, so one org can hold every section you teach:

```
witcomp1000/witcomp1000-fall26-a1-alice     # CS1, fall
witcomp1000/witcomp2000-spring27-a1-alice   # CS2, spring — different repo
```

It also makes the org browsable: repos for one section sort together, and archiving a finished term
is a matter of matching one prefix. Pick something stable that names the section and the term, and
keep it short — it is the leading third of every repo name your students will see.

## Usage

### Create a classroom

Reads the roster, verifies the org exists and that you have admin access, then writes the classroom
file and reports how many students it found.

    $ gh assgn-dist create --name CLASSROOM --org ORG

Arguments:

- `--name` and `--org` are required. `--name` is the classroom name that prefixes every student
  repo; it is stored in the classroom file and never passed again.
- `--roster` is the path to the roster CSV; defaults to `roster.csv`. It must exist and list at
  least one student — a missing or empty roster stops the command before anything is written.
  Whatever it resolves to is written into the classroom file, so `distribute` and `clone` read the
  same roster afterwards.
- `--classroom` sets the classroom file to write; defaults to `classroom.yaml`. Pass it to keep
  more than one classroom side by side (for example `--classroom cs2.yaml`).
- `--dry-run` still verifies the org and your admin access, then prints what would be written
  instead of writing the classroom file.

### Distribute an assignment

Verifies the template repo exists, then for each student in the roster (or just those named by
`--only`) creates a private repo from the template and adds the student as an outside collaborator
with `write` access.

**Safe to re-run**: if a student's repo already exists, creation is skipped, but the collaborator is
still (re-)added, so re-running also repairs any invite that failed on a prior run and picks up
students who don't have a repo yet.

    $ gh assgn-dist distribute --template OWNER/REPO ASSGN

Arguments:

- `--template` defaults to `ORG/ASSGN` (the classroom org and the assignment short name) if omitted.
- `--classroom` defaults to `classroom.yaml` if omitted.
- `--only` limits distribution to a comma-separated list of GitHub usernames (for example
  `--only alice,bob`); everyone else in the roster is skipped. Matching is case-insensitive, and a
  username that isn't in the roster is reported as a warning — if none of them match, the command
  stops without making changes. Omit it to distribute to everyone.
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
Name,Email,GitHub
Sunjae Park,sunjae@email.com,sunjaeatwit
John Doe,jdoe@email.com,jdoe11atwit
Jane Austen,jausten@email.com,jaustenatwit

# Create a new classroom under github.com/witcomp1000
$ gh assgn-dist create --name witcomp1000-fall26 --org witcomp1000
Created classroom.yaml (3 students in roster.csv)

# Distribute a1 using github.com/sunjae0802/cs1-a1 as template repo
$ gh assgn-dist distribute --template sunjae0802/cs1-a1 a1

# Distribute a1 to just one student (for example, a late add)
$ gh assgn-dist distribute --template sunjae0802/cs1-a1 --only jdoe11atwit a1

# Clone a1
$ gh assgn-dist clone a1
```

## Development

```bash
go build ./...
go test ./...
```

Project layout:

```
main.go                  # cobra root command
cmd/create.go            # `create` subcommand
cmd/distribute.go        # `distribute` subcommand
cmd/clone.go             # `clone` subcommand
internal/classroom.go    # classroom.yaml read/write
internal/roster.go       # CSV roster parsing
```
