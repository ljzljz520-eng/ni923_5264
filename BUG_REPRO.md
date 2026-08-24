# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/toolnav/cmd/toolnav	[no test files]
?   	example.com/toolnav/model	[no test files]
?   	example.com/toolnav/report	[no test files]
?   	example.com/toolnav/store	[no test files]
?   	example.com/toolnav/workflow	[no test files]
--- FAIL: TestToolBackupSeesConsistentOrder (0.00s)
    workflow_test.go:117: backup during drag failed: duplicate tool id "policykit"
FAIL
FAIL	example.com/toolnav	0.032s
ok  	example.com/toolnav/app	0.007s
ok  	example.com/toolnav/audit	0.002s
ok  	example.com/toolnav/backup	0.001s
ok  	example.com/toolnav/catalog	0.006s
ok  	example.com/toolnav/exporter	0.001s
ok  	example.com/toolnav/governance	0.001s
ok  	example.com/toolnav/importer	0.001s
ok  	example.com/toolnav/metrics	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/toolnav): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/toolnav): exit `0`
- Frontend build (web): exit `0`
