# assgn-dist Implementation Status

## Completed & Verified

### Infrastructure
- Added `cobra` and `gopkg.in/yaml.v3` dependencies
- Restructured project into `cmd/` and `internal/` packages

### Subcommands
- **`new`** — verifies org exists + user has admin access, writes `CLASSROOM.yaml`
- **`create`** — loads classroom + roster, verifies template repo, creates private student repos, adds students as outside collaborators with `write` role; supports `--dry-run`
- **`student-repos`** — clones missing repos or pulls existing ones

### File Structure
```
main.go                  # cobra root command
cmd/new.go               # new subcommand
cmd/create.go            # create subcommand
cmd/student_repos.go     # student-repos subcommand
internal/classroom.go    # CLASSROOM.yaml read/write
internal/roster.go       # CSV roster parsing
```

## Test Results (against org `witcomp-summer26-park`)

| Step | Result |
|------|--------|
| `new --org witcomp-summer26-park --roster roster.csv myclass` | ✅ Created `myclass.yaml` |
| `create --template witcomp-summer26-park/ch01-sample --dry-run a1` | ✅ Printed correct `gh api` commands |
| `create --template witcomp-summer26-park/ch01-sample a1` | ✅ Created private repo `myclass-a1-sunjae0802-150`, invited collaborator |
| `student-repos a1` (first run) | ✅ Cloned repo locally |
| `student-repos a1` (second run) | ✅ Ran `git pull` instead of clone |

- Test student: `Sunjae,sunjae0802@150mail.com,sunjae0802-150`
- Collaborator added as pending invitation (expected for outside collaborators)
- Pull warning on second run is expected — test template repo had no commits

## All tasks complete
