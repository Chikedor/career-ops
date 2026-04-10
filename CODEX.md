# Career-Ops for Codex

Career-Ops is local-first. Codex should reuse the checked-in scripts, templates, tracker flow, and `modes/` files instead of inventing a parallel workflow.

## Routing

- Job URL or raw JD: read `modes/_shared.md` + `modes/auto-pipeline.md`
- Single evaluation only: read `modes/_shared.md` + `modes/oferta.md`
- Multiple offers comparison: read `modes/_shared.md` + `modes/ofertas.md`
- Portal scan: read `modes/_shared.md` + `modes/scan.md`
- PDF generation: read `modes/_shared.md` + `modes/pdf.md`
- Pipeline inbox processing: read `modes/_shared.md` + `modes/pipeline.md`
- Tracker questions: read `modes/tracker.md`
- Batch orchestration: read `modes/batch.md` + `batch/batch-prompt.md`

## Rules

- Keep user-specific customization in `config/profile.yml`, `modes/_profile.md`, `article-digest.md`, or `portals.yml`.
- Do not write directly into `data/applications.md` from an agent flow when a TSV merge path exists. Use `batch/tracker-additions/*.tsv` plus `merge-tracker.mjs`.
- Treat `templates/states.yml` as the source of truth for canonical statuses.
- Never submit an application on the user's behalf.
- Prefer Playwright for live offer verification when available. If the environment does not support it, say so explicitly.

## Batch

- `batch/batch-runner.sh` uses `codex exec` as the documented headless Codex example.
- Set `CAREER_OPS_BATCH_WORKER` or pass `--worker-cmd` if you need a different non-interactive invocation.
- The worker wrapper feeds prompts over `stdin` by default; switch to `CAREER_OPS_BATCH_WORKER_INPUT=arg` if your local CLI expects the prompt as a final argument.
- The merge step assigns the final tracker row number to avoid parallel worker collisions.

## Validation

- `npm run doctor`
- `npm run verify`
- `cd dashboard && go build ./...`
