# Setup Guide

## Prerequisites

- Codex available in your local workflow
- Node.js 18+
- Playwright Chromium for PDF generation and live page checks
- Go 1.21+ if you want the dashboard
- Bash if you want to use `batch/batch-runner.sh`

## Install

```bash
git clone https://github.com/santifer/career-ops.git
cd career-ops
npm install
npx playwright install chromium
```

PowerShell equivalent for the local copies:

```powershell
Copy-Item config/profile.example.yml config/profile.yml
Copy-Item templates/portals.example.yml portals.yml
```

## Configure Your Local Data

```bash
cp config/profile.example.yml config/profile.yml
cp templates/portals.example.yml portals.yml
```

Then create:

- `cv.md` in the repo root
- `article-digest.md` if you want extra proof points
- `modes/_profile.md` if you want user-specific archetypes or negotiation language

## First Checks

```bash
npm run doctor
npm run verify
```

## Start Using It With Codex

Open Codex in the repo root and use prompts like:

- `Evaluate this job URL with the full career-ops pipeline.`
- `Read modes/_shared.md and modes/pdf.md, then generate the tailored PDF for this JD.`
- `Process the pending URLs in data/pipeline.md.`

## Batch

```bash
./batch/batch-runner.sh --dry-run
CAREER_OPS_BATCH_WORKER="codex exec" ./batch/batch-runner.sh --parallel 2
```

On Windows with Git Bash installed:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' -lc "cd '/c/Users/<you>/path/to/career-ops' && ./batch/batch-runner.sh --dry-run"
& 'C:\Program Files\Git\bin\bash.exe' -lc "cd '/c/Users/<you>/path/to/career-ops' && export CAREER_OPS_BATCH_WORKER='codex exec' && ./batch/batch-runner.sh --parallel 2"
```

If your CLI expects the prompt as a final argument instead of stdin:

```bash
CAREER_OPS_BATCH_WORKER="codex exec" CAREER_OPS_BATCH_WORKER_INPUT=arg ./batch/batch-runner.sh --parallel 2
```

## Dashboard

From the repo root:

```bash
go build -o dashboard/career-dashboard ./dashboard
./dashboard/career-dashboard --path .
```

You can also run it from inside `dashboard/`:

```bash
go build -o career-dashboard .
./career-dashboard --path ..
```

## Known Limits

- The documented non-interactive batch command is `codex exec`.
- If your local Codex CLI expects prompt input differently, set `CAREER_OPS_BATCH_WORKER_INPUT=arg` or override `CAREER_OPS_BATCH_WORKER`.
