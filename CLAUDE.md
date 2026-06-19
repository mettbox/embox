# EmBox — CLAUDE.md

> Full project context: see [docs/llms.txt](docs/llms.txt)

## Docker is not available in this environment

Claude cannot execute `docker` or `docker compose` commands. When database operations are needed, output the commands for the user to run locally:

```sh
docker compose up -d   # start database
docker compose down    # stop database
```

## Verification

| Area          | Command                                                        |
|---------------|----------------------------------------------------------------|
| Backend build | `cd api && go build ./...` and `go vet ./...`                  |
| Integration   | `cd api && go test ./internal/tests/... -v -count=1`           |
| Frontend      | `cd app && npm run lint` and `npm run build`                   |

Run the relevant check after every non-trivial change.

> Integration tests require `ffmpeg` in PATH and run against SQLite in-memory — no Docker needed.
> Tests are skipped gracefully when `ffmpeg` is not available.

## Deployment

Deployment is handled manually by the user — do not build, push, or deploy unless explicitly asked.

## LuckyCloud Integration

Original media files are stored in LuckyCloud via REST API.
API docs: https://docs.luckycloud.de/en/cloud-storage/rest-api

The local `api/media/` directory contains only WebP thumbnails. Original files live in LuckyCloud at path `yyyy/mm/dd_ID.FileExt`.

## Known Gotchas

**Port:** API runs on port `2705` by default.

**Thumbnails vs. originals:** Thumbnails are stored locally as WebP in `api/media/yyyy/mm/`. Original files (any format) are streamed directly from LuckyCloud. Never confuse local paths (`Path()`) with remote paths (`RemotePath()`) in the Media model.
