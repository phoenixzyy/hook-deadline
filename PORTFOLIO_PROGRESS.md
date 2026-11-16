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

## Publication checkpoint

Normal Fork verified: https://github.com/phoenixzyy/hook-deadline, forked from
adnanh/webhook; owner phoenixzyy ID 7067843. Eight upstream branches retained.
Feature commit 8172d3987876257bc3f44bc68f54262c83e09bb1 has tree
eb3dea589174254ce177656b6726b165bc527858, identical to locally verified commit
69a206e43720ca648e71ef05e4da5d6fe7ffb90f. Parent is the pinned baseline.
No force update. Actions enabled; this checkpoint triggers unfiltered remote CI.
CI results must be checked separately; no remote success is claimed here.

## Remote verification checkpoint

- Final code commit: 80ebc23dc9044bc7c27ed3ee86f1624dc7591019. A small follow-up
  logs invalid duration errors for asynchronous hooks; its local rejection tests passed.
- [CI run 43923772086](https://github.com/phoenixzyy/hook-deadline/actions/runs/43923772086):
  Ubuntu Go 1.21/1.22/1.26 and Windows Go 1.21/1.22 all succeeded. Each ran the
  full unfiltered build/test suite. Logs confirm both socket cases passed on Linux.
  Ubuntu Go 1.26 also passed the focused command-lifecycle race check.
- macOS Go 1.26 is still queued awaiting a runner as of this checkpoint; no macOS
  success is claimed. The live CI link is authoritative for later completion.
- Public fetch verifies remote history/tree. MIT and vendor licenses are unchanged.
- This documentation-only checkpoint uses [skip ci] to avoid repeating already
  verified code tests. It does not cancel or bypass the pending macOS code run.
- No deployment, service, package release, purchase, upstream contact or PR.
  Process-tree termination, durable delivery and output-size limits remain out of scope.
