# hook-deadline contribution and verification

## Source and scope

Normal Fork of [adnanh/webhook](https://github.com/adnanh/webhook), based on
`2bbdeb9f90b5f98c7b7b8976881efd6da9f69571`. MIT LICENSE retains the original
2015 Adnan Hajdarevic notice. Vendor licenses remain unchanged. Existing HTTP
routing, authentication/signature rules, command argument mapping, config reload
and output responses belong to upstream. This change was prepared with Codex
assistance; it does not demonstrate independent mastery or production experience.

Reviewed recent changes #769, #765 and #764, and related proposals #768/#759
(request read limits) and #703 (socket listeners). This work does not copy those
patches or claim their features. No upstream issue, PR or message was created.

## Reproduced need

The baseline ignores unknown JSON fields, including `command-timeout`. A synthetic
command configured with `100ms` ran to completion in about 2 seconds. The regression
failed before implementation. This is an optional capability added to upstream,
not a claim that upstream promised timeout behavior and violated it.

Example configuration fragment (keep your existing authenticated trigger rule):

```json
{"id":"bounded-job","execute-command":"/absolute/path/to/your-command",
 "command-timeout":"30s","include-command-output-in-response":true}
```

Only the new field and lifecycle logic are added. A positive Go duration creates
a separate context for the invocation, cancels the direct process on expiry, and
limits inherited output-pipe waiting to one second after exit/cancellation.
The timer starts after executable lookup, before argument/file preparation. It
is not a request-body deadline or hard real-time bound on all handler operations.
Omitted, empty, `0` and `0s` keep unlimited legacy behavior and pipe waiting.
Invalid/negative/overflowing durations return an error before execution when the
hook is invoked; they are not rejected by configuration loading itself.

Timeout errors wrap `context.DeadlineExceeded`. Ordinary nonzero exits retain
`*exec.ExitError`. In captured-output mode, the existing HTTP 500 behavior applies;
partial output is returned only if the existing error-output option permits it.
Asynchronous hooks still return their configured acknowledgement immediately;
they share the execution timeout but cannot change an already returned response.
Successfully created payload temporary files follow the existing post-command
cleanup path, now also after a timeout. No unrelated tempfile error path is fixed.

## Verification

See PORTFOLIO_PROGRESS.md for exact results and remote CI evidence. Use Go 1.21+
(the go.mod toolchain declaration may select 1.22 automatically):

```sh
go test ./...
go test -race -run '^TestCommandTimeout' -count=1
go build
```

Tests execute only short synthetic helper processes, use fake payloads and cover
expiry/partial output, successful defaults, nonzero exits, invalid durations,
JSON/YAML loading, payload-file cleanup, HTTP response handling and inherited
pipes. The pipe test creates a finite child, records its PID and terminates it
after checking WaitDelay. It does not pretend the implementation kills descendants.

The local Linux sandbox does not allow Unix sockets: two pre-existing socket
subtests fail both before changes and independently of timeout logic. The local
full supported run uses `go test ./... -skip 'TestWebhook/.*:socket@' -count=1`.
Remote CI runs the **unfiltered** upstream suite on Linux, Windows and macOS, plus
a focused Linux race check. Do not call full baseline green without that caveat.

## Limitations

This is not process-tree termination, graceful shutdown, output-size limiting,
sandboxing, retry, persistence or a delivery queue. Descendants may continue
running; inherited pipes are closed after the one-second grace period and output
can be truncated. A normal parent exit with an inherited open pipe can return
`exec.ErrWaitDelay` even before the configured execution deadline. Commands still
run with the server user's privileges. Configure trusted commands/arguments and
existing signature rules; a timeout is not an authorization boundary. No real
secrets, external jobs, deployment, online service or paid resource was used.
