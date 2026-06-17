# EmBox — CLAUDE.md

> Full project context: see [docs/llms.txt](docs/llms.txt)

## Docker is not available in this environment

Claude cannot execute `docker` or `docker compose` commands. When database operations are needed, output the commands for the user to run locally:

```sh
docker compose up -d   # start database
docker compose down    # stop database
```

## Verification (no tests yet)

Until a test suite exists, verify changes with:

| Area     | Command                                           |
|----------|---------------------------------------------------|
| Backend  | `cd api && go build ./...` and `go vet ./...`     |
| Frontend | `cd app && npm run lint` and `npm run build`      |

Run the relevant check after every non-trivial change.

## Deployment

Deployment is handled manually by the user — do not build, push, or deploy unless explicitly asked.

## LuckyCloud Integration

Original media files are stored in LuckyCloud via REST API.
API docs: https://docs.luckycloud.de/en/cloud-storage/rest-api

The local `api/media/` directory contains only WebP thumbnails. Original files live in LuckyCloud at path `yyyy/mm/dd_ID.FileExt`.

## Known Gotchas

**GORM cascade delete (album_media):** GORM does not correctly apply `ON DELETE CASCADE` for many-to-many relations. After migrations, the FK must be manually fixed in the DB:

```sql
ALTER TABLE album_media DROP FOREIGN KEY fk_albums_album_media;
ALTER TABLE album_media
ADD CONSTRAINT fk_albums_album_media
FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE;
```

**Port:** API runs on port `2705` by default.

**Thumbnails vs. originals:** Thumbnails are stored locally as WebP in `api/media/yyyy/mm/`. Original files (any format) are streamed directly from LuckyCloud. Never confuse local paths (`Path()`) with remote paths (`RemotePath()`) in the Media model.
