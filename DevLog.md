 Machine & Devbox Setup
Verified you were using WSL2 with VS Code and .venv for Python tooling.

Wanted to check system specs — used:

bash
Copy
Edit
lscpu
free -h
df -h
to verify CPU, RAM, and disk availability.

🔎 Kubernetes Cluster Runtime Surprise
Created a K3D cluster in a terminal.

Closed the terminal and were surprised the cluster was still running.

Lesson: K3D runs as Docker containers, which are managed by the Docker daemon, not your shell session. Restarting the terminal doesn’t affect the cluster.

If your machine reboots and Docker isn’t set to auto-start the containers, then yes — your cluster dies until manually restarted.

🧵 Git Tracking Issues
Tried to add the obs-k8s-demo/ folder but Git treated it as a submodule.

Resolved by:

bash
Copy
Edit
git rm --cached obs-k8s-demo
rm .gitmodules
git add obs-k8s-demo/
git commit -m "Fix: track demo folder as part of main tree"
❌ Pre-commit Hook Blockers
Pre-commit failed due to:

Missing --- in YAML files

Mis-indented PrometheusRule.yaml

Line-length violations

Syntax errors like mapping values not allowed here

✅ Recovery Steps:
Created .pre-commit-config.yaml and used:

yamllint

yamlfmt

unicode-sanitizer for invisible junk

safe-commit.sh to bypass hook during dev

fix-skipped.py to re-run hooks on skipped files later

Realized yamllint doesn’t fix structural YAML errors.

Installed ruamel.yaml to help auto-format YAML with:

python
Copy
Edit
from ruamel.yaml import YAML
yaml = YAML()
yaml.preserve_quotes = True
with open("file.yaml") as f: data = yaml.load(f)
yaml.dump(data, sys.stdout)
🧠 AI-Savvy Linting Goals
Wanted to prevent:

Invisible characters (from copy-paste)

Malicious <!-- hidden commands --> from LLM-suggested code

Setup a fixer script that logs issues and diffs before + after

🛠️ GO Installation Issues
Tried installing yamlfmt via Go:

bash
Copy
Edit
go install github.com/google/yamlfmt/cmd/yamlfmt@latest
Got errors like defined in both Go and assembly — likely due to incorrect unpacking or partial Go setup.

Recommended: use system package manager (apt, brew) or install from the official Go tarball cleanly.

Let me know if you want this synced to your GitHub in a specific format (like a CHANGELOG.md, devlog.md, or as commits inside your repo). I can also build a simple devlog-tracker.py that appends to your logs when you run into stuff like this again.
