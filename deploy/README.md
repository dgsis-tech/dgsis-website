# Deploy artifacts (VPS) — dgsis-website

Reference files for the public website on the shared Traefik VPS.

Policy mirrored from [`dgsis-tech/academy`](https://github.com/dgsis-tech/academy)
(`deploy/` + `.github/workflows/`). Adapted because this repo is a single Go
server with no database.

CI/CD workflow: [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml)
and [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).

| Concern | Choice |
| --- | --- |
| Traefik | **Global** on the VPS (`external` Docker network `frontend`). Not started here. |
| Blue/green | **Local** color switch via `traefik.enable` + `/var/lib/dgsis-website/active-color` |
| Postgres | None. No backup and no migrate step. |
| Image | Built **on the VPS** by CD from the repo-root `Dockerfile`. `WEBSITE_IMAGE` is a **local** Docker name (default `ghcr.io/dgsis-tech/dgsis-website`). No GHCR push/pull. |
| Environments | **Single production stack** (no separate staging hostname). |
| GitHub secrets | **None.** Host `.env` holds Traefik hostname and port only. |
| Runner | Repo-scoped self-hosted runner for **this** repo, label `self-hosted`. Academy's runner does not pick up these workflows. |

## What lives on the VPS vs in the checkout

CD checks out the tag on the self-hosted runner, builds the image, then runs
`deploy/deploy-blue-green.sh` **from that checkout** with:

- `ENV_FILE=/srv/docker/compose/dgsis-website/.env`
- `COMPOSE_FILE=/srv/docker/compose/dgsis-website/docker-compose.yml`
- `STATE_DIR=/var/lib/dgsis-website`

Keep on the host (update manually when Compose changes):

| Path | Purpose |
| --- | --- |
| `/srv/docker/compose/dgsis-website/docker-compose.yml` | Copy from `deploy/docker-compose.yml` when it changes |
| `/srv/docker/compose/dgsis-website/.env` | From `.env.example` — Traefik network and public host |
| `/var/lib/dgsis-website` | `active-color` |
| `$HOME/.cache/dgsis-website-buildx` | Buildx local cache (runner user; **not** `/var/cache`, **not** academy's cache) |

Do not write into academy paths (`/srv/docker/compose/academy`, `/var/lib/academy`, `/var/backups/academy`, `$HOME/.cache/academy-buildx`).

Also ensure:

- Docker network **`frontend`** already exists (global Traefik edge on this VPS).
- A self-hosted runner is registered to **`dgsis-tech/dgsis-website`** (label `self-hosted`) with Docker + Buildx. Same machine as academy is fine; it must be a registration for this repository (or an org runner — see the operator notes). Academy's existing repo runner will not run these jobs.
- Runner user can write `$HOME/.cache/dgsis-website-buildx` and `/var/lib/dgsis-website` (no sudo in the job).
- DNS for `WEBSITE_HOST` points at the Traefik VPS. The example host is `dgsis.com`; confirm it before the first tag.
- Runtime image includes `wget` so the Compose healthcheck can call `GET /health`.

Post-job CD/CI: [`scripts/cd-prune-host.sh`](../scripts/cd-prune-host.sh) (dangling images, old website tags, Buildx size cap) then empty `$GITHUB_WORKSPACE`. **Never** `docker system prune -af`.

**Never** start the stack with bare `docker compose up` without the Traefik enable flags that `deploy-blue-green.sh` sets — defaults leave both colors with `traefik.enable=false` and the host returns 404.

## One-time host prep

```bash
sudo mkdir -p /srv/docker/compose/dgsis-website /var/lib/dgsis-website
sudo chown -R github-runner:github-runner \
  /srv/docker/compose/dgsis-website /var/lib/dgsis-website
mkdir -p /home/github-runner/.cache/dgsis-website-buildx
```

Copy `deploy/docker-compose.yml` and a filled `.env` into the compose directory **before** the first tag. Re-copy the compose file only when it changes.

## Local image build (same as CD)

```bash
docker build -f Dockerfile -t ghcr.io/dgsis-tech/dgsis-website:local .
# Tag name is local only — CD does not push/pull this registry path.
```

## Release

Production only after an explicit operator request. Preferred trigger is an annotated tag `v*`. `workflow_dispatch` can redeploy an existing tag. There is no automatic rollback: a failed health check stops the script and leaves the previous color serving if it was still up.
