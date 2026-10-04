#!/usr/bin/env bash
# cd-prune-host.sh — bounded host prune after CD/CI (images and buildx).
#
# Idempotent. Never runs `docker system prune -af` (would kill Traefik and
# other products, including academy). Does not use sudo. Image `rmi` failures
# (in use) are logged; exit stays 0.
#
# Usage:
#   ./scripts/cd-prune-host.sh [--dry-run] [--help]
#
# Environment:
#   WEBSITE_IMAGE      default ghcr.io/dgsis-tech/dgsis-website
#   BACKUP_DIR         default empty (this site has no dumps; set only if added later)
#   BUILDX_CACHE       default $HOME/.cache/dgsis-website-buildx
#   KEEP_RELEASES      default 3 (extra v* tags beyond images used by running containers)
#   KEEP_BACKUPS       default 7 (unused while BACKUP_DIR is empty)
#   BUILDX_MAX_BYTES   default 4294967296 (4 GiB)
#   EXTRA_PRUNE_IMAGES optional space-separated refs (e.g. dgsis-website:ci)
#
set -euo pipefail

DRY_RUN=0

usage() {
  cat <<'EOF'
Usage:
  ./scripts/cd-prune-host.sh [--dry-run] [--help]

Bounded prune for the dgsis-website VPS runner after CD/CI:
  1. docker image prune -f (dangling only)
  2. drop unused tags of WEBSITE_IMAGE (keep running + last KEEP_RELEASES v* tags)
  3. if BACKUP_DIR is set, retain KEEP_BACKUPS newest dgsis-website-backup-*.sql.gz
  4. if BUILDX_CACHE is >= BUILDX_MAX_BYTES, remove that directory
  5. optionally remove EXTRA_PRUNE_IMAGES

Never runs docker system prune -af. Never uses sudo. Never deletes $HOME
or other products' paths (/srv/docker/compose/academy, /var/lib/academy, …).

Options:
  --dry-run   Print actions without deleting
  --help      Show this help

Environment: WEBSITE_IMAGE, BACKUP_DIR, BUILDX_CACHE, KEEP_RELEASES,
KEEP_BACKUPS, BUILDX_MAX_BYTES, EXTRA_PRUNE_IMAGES (see script header).
EOF
}

log() { printf 'cd-prune: %s\n' "$*"; }
log_warn() { printf 'cd-prune: WARN: %s\n' "$*" >&2; }

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      log_warn "unknown argument: $1"
      usage
      exit 2
      ;;
  esac
done

WEBSITE_IMAGE="${WEBSITE_IMAGE:-ghcr.io/dgsis-tech/dgsis-website}"
BACKUP_DIR="${BACKUP_DIR:-}"
BUILDX_CACHE="${BUILDX_CACHE:-${HOME}/.cache/dgsis-website-buildx}"
KEEP_RELEASES="${KEEP_RELEASES:-3}"
KEEP_BACKUPS="${KEEP_BACKUPS:-7}"
BUILDX_MAX_BYTES="${BUILDX_MAX_BYTES:-4294967296}"
EXTRA_PRUNE_IMAGES="${EXTRA_PRUNE_IMAGES:-}"

run_or_echo() {
  if [[ "$DRY_RUN" -eq 1 ]]; then
    log "dry-run: $*"
  else
    "$@"
  fi
}

# --- 1. dangling images ---
if command -v docker >/dev/null 2>&1; then
  if [[ "$DRY_RUN" -eq 1 ]]; then
    log "dry-run: docker image prune -f"
  else
    docker image prune -f || log_warn "docker image prune failed (continuing)"
  fi
else
  log_warn "docker not found; skipping image prune"
fi

# --- 2. website image tags ---
if command -v docker >/dev/null 2>&1; then
  declare -A keep_tags=()
  declare -A in_use_ids=()

  while IFS= read -r cid; do
    [[ -z "$cid" ]] && continue
    img_id="$(docker inspect --format '{{.Image}}' "$cid" 2>/dev/null || true)"
    [[ -n "$img_id" ]] && in_use_ids["$img_id"]=1
  done < <(docker ps -q 2>/dev/null || true)

  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    tag="${line%%|*}"
    id="${line##*|}"
    [[ "$tag" == "<none>" ]] && continue
    if [[ -n "${in_use_ids[$id]:-}" ]]; then
      keep_tags["$tag"]=1
    fi
  done < <(docker images --format '{{.Tag}}|{{.ID}}' "$WEBSITE_IMAGE" 2>/dev/null || true)

  mapfile -t v_tags < <(
    docker images --format '{{.Tag}}' "$WEBSITE_IMAGE" 2>/dev/null \
      | grep -E '^v[0-9]' \
      | sort -V \
      | tail -n "$KEEP_RELEASES" \
      || true
  )
  for t in "${v_tags[@]+"${v_tags[@]}"}"; do
    [[ -n "$t" ]] && keep_tags["$t"]=1
  done

  while IFS= read -r tag; do
    [[ -z "$tag" || "$tag" == "<none>" ]] && continue
    if [[ -n "${keep_tags[$tag]:-}" ]]; then
      continue
    fi
    ref="${WEBSITE_IMAGE}:${tag}"
    id="$(docker images --format '{{.ID}}' "$ref" 2>/dev/null | head -n1 || true)"
    if [[ -n "$id" && -n "${in_use_ids[$id]:-}" ]]; then
      log "keep in-use $ref"
      continue
    fi
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log "dry-run: docker rmi $ref"
    else
      if docker rmi "$ref" 2>/dev/null; then
        log "removed $ref"
      else
        log_warn "could not remove $ref (in use or missing); continuing"
      fi
    fi
  done < <(docker images --format '{{.Tag}}' "$WEBSITE_IMAGE" 2>/dev/null || true)

  # shellcheck disable=SC2086
  for ref in $EXTRA_PRUNE_IMAGES; do
    [[ -z "$ref" ]] && continue
    if [[ "$DRY_RUN" -eq 1 ]]; then
      log "dry-run: docker rmi $ref"
    else
      if docker rmi "$ref" 2>/dev/null; then
        log "removed $ref"
      else
        log_warn "could not remove $ref; continuing"
      fi
    fi
  done
fi

# --- 3. backup retention (off unless BACKUP_DIR is set) ---
if [[ -n "$BACKUP_DIR" && -d "$BACKUP_DIR" ]]; then
  mapfile -t dumps < <(ls -1t "$BACKUP_DIR"/dgsis-website-backup-*.sql.gz 2>/dev/null || true)
  if ((${#dumps[@]} > KEEP_BACKUPS)); then
    for ((i = KEEP_BACKUPS; i < ${#dumps[@]}; i++)); do
      run_or_echo rm -f "${dumps[$i]}"
      log "backup prune: ${dumps[$i]}"
    done
  else
    log "backups: ${#dumps[@]} <= KEEP_BACKUPS=${KEEP_BACKUPS}; nothing to drop"
  fi
else
  log "BACKUP_DIR unset or missing; skip dump retention"
fi

# --- 4. buildx cache size cap ---
if [[ -d "$BUILDX_CACHE" ]]; then
  size="$(du -sb "$BUILDX_CACHE" 2>/dev/null | cut -f1 || echo 0)"
  if ((size >= BUILDX_MAX_BYTES)); then
    log "BUILDX_CACHE ${BUILDX_CACHE} is ${size} bytes (>= ${BUILDX_MAX_BYTES}); removing"
    run_or_echo rm -rf "$BUILDX_CACHE"
  else
    log "BUILDX_CACHE ${BUILDX_CACHE} size=${size} under cap"
  fi
else
  log "BUILDX_CACHE ${BUILDX_CACHE} absent; skip"
fi

log "done (dry_run=${DRY_RUN})"
exit 0
