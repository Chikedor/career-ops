# Codex Workflow

This repo is now organized to work with Codex as the primary local agent.

## What Codex Reads

- Entry point: `AGENTS.md`
- Main routing rules: `CODEX.md`
- Task-specific instructions: `modes/*.md`
- User customization: `config/profile.yml`, `modes/_profile.md`, `article-digest.md`, `portals.yml`

## Recommended Prompts

- `Evaluate this job URL with the full career-ops pipeline.`
- `Read data/pipeline.md and process the pending offers one by one.`
- `Generate the ATS PDF for this offer using my current cv.md and profile.`
- `Review batch/batch-input.tsv and run the batch pipeline.`

## Batch Mode

`batch/batch-runner.sh` runs a configurable headless worker command:

```bash
CAREER_OPS_BATCH_WORKER="codex exec" ./batch/batch-runner.sh --parallel 2
CAREER_OPS_BATCH_WORKER_INPUT=stdin ./batch/batch-runner.sh --parallel 2
```

Notes:

- The validated headless Codex command is `codex exec`.
- Default input mode is `stdin`; set `CAREER_OPS_BATCH_WORKER_INPUT=arg` if your local CLI expects the prompt as a final argument instead.
- If your local Codex CLI exposes a different non-interactive command, override it with `CAREER_OPS_BATCH_WORKER` or `--worker-cmd`.
- Final tracker numbering is assigned during merge, not by each worker, to avoid collisions in parallel runs.

## Windows Validation

Check what PowerShell can resolve:

```powershell
Get-Command codex
where.exe codex
```

Then test whether the CLI is available from your shell:

```powershell
codex --help
codex --version
```

## Files That Still Matter

- `generate-pdf.mjs`: PDF generation
- `merge-tracker.mjs`: merge TSV additions into tracker
- `verify-pipeline.mjs`: pipeline integrity check
- `normalize-statuses.mjs`: backward-compatible status cleanup
- `dashboard/`: local TUI pipeline viewer

## Validation

```bash
npm run doctor
npm run verify
cd dashboard && go build ./...
```
