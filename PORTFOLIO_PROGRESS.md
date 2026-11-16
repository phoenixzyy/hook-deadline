# hook-deadline progress

## Local verification — 2026-11-17 (Asia/Shanghai)

- Account reverified: phoenixzyy, ID 7067843. Source adnanh/webhook fixed at
  2bbdeb9f90b5f98c7b7b8976881efd6da9f69571; full history cloned.
- MIT LICENSE and vendor notices unchanged. Read CONTRIBUTING.md; no applicable
  repository AGENTS.md/SKILL.md found. No upstream contact or copied PR patch.
- Reviewed recent changes and related PRs #769/#765/#764/#768/#759/#703.
- Unmodified full baseline failed only two TestWebhook socket subtests in this
  container. Original supported suite with the exact socket exclusion passed:
  79 Go test pass events (includes parent tests/subtests, not 79 independent cases).
- Red reproduction: JSON command-timeout=100ms was ignored and a synthetic child
  completed after about 2 seconds. One new test failed before production changes.
- Implemented optional duration, direct-child cancellation, one-second inherited
  output-pipe wait, wrapped deadline errors. Legacy default behavior retained.
- Final supported suite: 96 pass events, no failures. New focused race check
  passed; Go build and built binary -version passed. Go 1.26.3 Linux.
- Tests cover default/zero/valid durations, invalid/negative/overflow values,
  nonzero exit codes, partial output, temporary-file cleanup, inherited pipes,
  HTTP error-output policy and JSON/YAML loading. All subprocesses synthetic.
- English contribution guide and Chinese learning guide explain architecture,
  evidence, design, interviews/exercises and limits. No fabricated experience.
- CI retains upstream unfiltered Linux/Windows Go 1.21/1.22 and macOS Go 1.26,
  adds Ubuntu Go 1.26 focused race check, and has read-only contents permissions.
- Logs: webhook-baseline.log, webhook-baseline-tcp.jsonl, webhook-red.log,
  webhook-focused.log, webhook-final.jsonl and webhook-race.log in workspace.
- Limits: does not kill process trees, enforce output size, gracefully shut down,
  persist/retry jobs or provide a delivery queue. Invalid durations fail when
  invoked, not at config load. Descendants can survive. No deployment/purchase.
- Normal Fork creation requested as phoenixzyy/hook-deadline after validation;
  remote publication and CI still require verification.
