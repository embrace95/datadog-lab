#!/bin/bash

# Commit message passed as an argument
COMMIT_MSG="$1"

if [ -z "$COMMIT_MSG" ]; then
  echo "❌ Please provide a commit message"
  exit 1
fi

# Log skipped files
SKIPPED_LOG=".precommit-skips.log"
echo "🔕 Skipping pre-commit hooks for this commit at $(date)" >> "$SKIPPED_LOG"
git diff --cached --name-only >> "$SKIPPED_LOG"
echo "---" >> "$SKIPPED_LOG"

# Commit with --no-verify
git commit -m "$COMMIT_MSG" --no-verify
