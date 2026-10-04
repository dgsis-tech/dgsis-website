#!/usr/bin/env bash
# Blue/green deploy for dgsis-website behind the VPS-global Traefik.
# Pattern: require local image → up inactive (no router) → health → switch →
# stop previous.
#
# No database, backup, or migration. CD builds the image on this host first;
# this script never docker-pulls the app.
#
# Compose + .env live under /srv/docker/compose/dgsis-website/; this script may
# run from a git checkout (set ENV_FILE / COMPOSE_FILE to absolute VPS paths).
#
# Do not `source` the Compose .env in bash: values may contain `$` that bash
# would corrupt. Compose reads ENV_FILE via --env-file; caller exports
# IMAGE_TAG (and optional WEBSITE_IMAGE) which override the file.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

ENV_FILE="${ENV_FILE:-$ROOT/.env}"
COMPOSE_FILE="${COMPOSE_FILE:-$ROOT/docker-compose.yml}"
STATE_DIR="${STATE_DIR:-/var/lib/dgsis-website}"
STATE_FILE="${STATE_FILE:-$STATE_DIR/active-color}"
HEALTH_TIMEOUT_SEC="${HEALTH_TIMEOUT_SEC:-120}"
WEBSITE_IMAGE="${WEBSITE_IMAGE:-ghcr.io/dgsis-tech/dgsis-website}"

: "${IMAGE_TAG:?IMAGE_TAG is required (e.g. v0.1.0)}"

if [[ ! -f "$ENV_FILE" ]]; then
	echo "error: missing env file $ENV_FILE (copy from deploy/.env.example)" >&2
	exit 1
fi
if [[ ! -f "$COMPOSE_FILE" ]]; then
	echo "error: missing compose file $COMPOSE_FILE" >&2
	exit 1
fi
if ! grep -qE '^TRAEFIK_NETWORK=.+' "$ENV_FILE"; then
	echo "error: TRAEFIK_NETWORK must be set in $ENV_FILE" >&2
	exit 1
fi
if ! grep -qE '^WEBSITE_HOST=.+' "$ENV_FILE"; then
	echo "error: WEBSITE_HOST must be set in $ENV_FILE" >&2
	exit 1
fi

TRAEFIK_NET_NAME="$(grep -E '^TRAEFIK_NETWORK=' "$ENV_FILE" | head -n1 | cut -d= -f2- | tr -d '"' | tr -d "'")"
if [[ -z "$TRAEFIK_NET_NAME" ]]; then
	echo "error: TRAEFIK_NETWORK value empty in $ENV_FILE" >&2
	exit 1
fi
if ! docker network inspect "$TRAEFIK_NET_NAME" >/dev/null 2>&1; then
	echo "error: Docker network '$TRAEFIK_NET_NAME' not found (global Traefik on this VPS)" >&2
	exit 1
fi

mkdir -p "$STATE_DIR"

ACTIVE=""
if [[ -f "$STATE_FILE" ]]; then
	ACTIVE="$(tr -d '[:space:]' <"$STATE_FILE")"
fi

case "$ACTIVE" in
blue) INACTIVE=green ;;
green) INACTIVE=blue ;;
"") INACTIVE=blue ;;
*)
	echo "error: unknown ACTIVE_COLOR='$ACTIVE' in $STATE_FILE (expected blue|green)" >&2
	exit 1
	;;
esac

echo "deploy: active=${ACTIVE:-none} target=$INACTIVE image=${WEBSITE_IMAGE}:${IMAGE_TAG}"

export WEBSITE_IMAGE IMAGE_TAG

compose() {
	docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" "$@"
}

set_traefik_enable() {
	# $1 = color that should advertise the Traefik router (or empty for none)
	local on="${1:-}"
	if [[ "$on" == "blue" ]]; then
		export BLUE_TRAEFIK_ENABLE=true
		export GREEN_TRAEFIK_ENABLE=false
	elif [[ "$on" == "green" ]]; then
		export BLUE_TRAEFIK_ENABLE=false
		export GREEN_TRAEFIK_ENABLE=true
	else
		export BLUE_TRAEFIK_ENABLE=false
		export GREEN_TRAEFIK_ENABLE=false
	fi
}

wait_healthy() {
	local service="$1"
	local deadline=$((SECONDS + HEALTH_TIMEOUT_SEC))
	echo "waiting for healthy: $service (timeout ${HEALTH_TIMEOUT_SEC}s)"
	while ((SECONDS < deadline)); do
		local cid status
		cid="$(compose ps -q "$service" 2>/dev/null || true)"
		if [[ -n "$cid" ]]; then
			status="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$cid" 2>/dev/null || true)"
			if [[ "$status" == "healthy" ]]; then
				echo "healthy: $service"
				return 0
			fi
			echo "  status=$status"
		fi
		sleep 2
	done
	echo "error: $service not healthy within ${HEALTH_TIMEOUT_SEC}s" >&2
	compose ps "$service" || true
	docker logs "$(compose ps -q "$service")" 2>/dev/null | tail -n 50 || true
	exit 1
}

# Fixed container_name + `compose up --force-recreate` races create-before-remove
# ("name already in use"). Remove first, then up without --force-recreate.
up_color() {
	local color="$1"
	local service="website-${color}"
	compose stop "$service" 2>/dev/null || true
	compose rm -f "$service" 2>/dev/null || true
	docker rm -f "dgsis-website-${color}" 2>/dev/null || true
	compose up -d --pull never "$service"
}

if ! docker image inspect "${WEBSITE_IMAGE}:${IMAGE_TAG}" >/dev/null 2>&1; then
	echo "error: image ${WEBSITE_IMAGE}:${IMAGE_TAG} not found locally" >&2
	echo "error: build the image on this host first; CD does not pull" >&2
	exit 1
fi

# 1) Bring up the inactive color with Traefik disabled; keep current active serving.
set_traefik_enable "$ACTIVE"
if [[ "$INACTIVE" == "blue" ]]; then
	export BLUE_TRAEFIK_ENABLE=false
else
	export GREEN_TRAEFIK_ENABLE=false
fi

up_color "$INACTIVE"
wait_healthy "website-${INACTIVE}"

# 2) Flip Traefik onto the new color only.
set_traefik_enable "$INACTIVE"
up_color "$INACTIVE"
if [[ -n "$ACTIVE" && "$ACTIVE" != "$INACTIVE" ]]; then
	compose stop "website-${ACTIVE}"
	compose rm -f "website-${ACTIVE}"
fi

printf '%s\n' "$INACTIVE" >"$STATE_FILE"
echo "deploy complete: ACTIVE_COLOR=$INACTIVE (state $STATE_FILE)"
