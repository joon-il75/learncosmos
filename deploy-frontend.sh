#!/usr/bin/env bash
# LearnCosmos frontend deploy script
# Derived from the backed-up deploy-learncosmos-frontend.sh and adapted for /home/cosmos/LearnCosmos.
# Defaults to port 3001 to avoid touching the legacy LearnWeaver frontend on 3000.

set -euo pipefail

PROJECT_ROOT="${PROJECT_ROOT:-/home/cosmos/LearnCosmos}"
FRONTEND_DIR="${PROJECT_ROOT}/frontend"
FRONTEND_PORT="${FRONTEND_PORT:-3001}"
FRONTEND_HOSTNAME="${FRONTEND_HOSTNAME:-0.0.0.0}"
INTERNAL_API_URL="${INTERNAL_API_URL:-http://127.0.0.1:8081}"
BUILD_LOG="${BUILD_LOG:-/tmp/learncosmos-frontend-build.log}"
FRONTEND_LOG="${FRONTEND_LOG:-/tmp/learncosmos-frontend.log}"
LEARNWEAVER_PROJECT_ROOT="${LEARNWEAVER_PROJECT_ROOT:-${PROJECT_ROOT}}"
PLANET_TEXTURE_MAP_DIR="${PLANET_TEXTURE_MAP_DIR:-${FRONTEND_DIR}/public/textures/planets/maps}"
ALLOW_LEARNCOSMOS_PROD_PORTS="${ALLOW_LEARNCOSMOS_PROD_PORTS:-0}"

if [ "${FRONTEND_PORT}" = "3000" ] && [ "${ALLOW_LEARNCOSMOS_PROD_PORTS}" != "1" ]; then
  echo "Refusing to use legacy LearnWeaver frontend port 3000. Use FRONTEND_PORT=3001 for LearnCosmos verification." >&2
  exit 1
fi

if [ ! -d "${FRONTEND_DIR}" ]; then
  echo "Frontend directory not found: ${FRONTEND_DIR}" >&2
  echo "Restore or create the frontend source before running this deploy script." >&2
  exit 1
fi

if [ ! -f "${FRONTEND_DIR}/package.json" ]; then
  echo "frontend/package.json not found: ${FRONTEND_DIR}/package.json" >&2
  exit 1
fi

cd "${FRONTEND_DIR}"

echo "Building LearnCosmos frontend..."
export NODE_OPTIONS="${NODE_OPTIONS:-} --max-old-space-size=4096"
export INTERNAL_API_URL PROJECT_ROOT LEARNWEAVER_PROJECT_ROOT PLANET_TEXTURE_MAP_DIR
npm run build 2>&1 | tee "${BUILD_LOG}" || {
  echo "Build failed" >&2
  tail -n 120 "${BUILD_LOG}" || true
  exit 1
}

echo "Checking standalone assets..."
test -f .next/standalone/server.js || { echo "standalone/server.js not found" >&2; exit 1; }
test -d .next/standalone/public/images || { echo "public/images was not copied" >&2; exit 1; }
test -d .next/standalone/public/videos || { echo "public/videos was not copied" >&2; exit 1; }
test -f .next/standalone/public/images/favicon.webp || { echo "favicon.webp was not copied" >&2; exit 1; }
test -d .next/standalone/.next/static || { echo ".next/static was not copied" >&2; exit 1; }

if [ -f public/videos/LearnWeaver-Seamless-Loop.mp4 ]; then
  test -f .next/standalone/public/videos/LearnWeaver-Seamless-Loop.mp4 || { echo "representative video was not copied" >&2; exit 1; }
fi

echo "Restarting LearnCosmos frontend on port ${FRONTEND_PORT}..."
PIDS="$(lsof -t -iTCP:"${FRONTEND_PORT}" -sTCP:LISTEN 2>/dev/null || true)"
if [ -z "${PIDS:-}" ] && command -v fuser >/dev/null 2>&1; then
  PIDS="$(fuser -n tcp "${FRONTEND_PORT}" 2>/dev/null || true)"
fi
if [ -z "${PIDS:-}" ]; then
  PIDS="$(ss -ltnp 2>/dev/null | awk -v port=":${FRONTEND_PORT}" '$4 ~ port { gsub(/.*pid=/, "", $0); gsub(/,.*/, "", $0); if ($0 ~ /^[0-9]+$/) print $0 }' || true)"
fi
if [ -n "${PIDS:-}" ]; then
  kill ${PIDS} || true
  sleep 1
fi

cd .next/standalone
setsid -f env PORT="${FRONTEND_PORT}" HOSTNAME="${FRONTEND_HOSTNAME}" INTERNAL_API_URL="${INTERNAL_API_URL}" PROJECT_ROOT="${PROJECT_ROOT}" LEARNWEAVER_PROJECT_ROOT="${LEARNWEAVER_PROJECT_ROOT}" PLANET_TEXTURE_MAP_DIR="${PLANET_TEXTURE_MAP_DIR}" node server.js >"${FRONTEND_LOG}" 2>&1

echo "Checking LearnCosmos frontend port ${FRONTEND_PORT}..."
for _ in 1 2 3 4 5; do
  if ss -ltn | grep -q ":${FRONTEND_PORT} "; then
    echo "LearnCosmos frontend listening on :${FRONTEND_PORT}"
    exit 0
  fi
  sleep 1
done

echo "LearnCosmos frontend did not start on :${FRONTEND_PORT}" >&2
tail -n 80 "${FRONTEND_LOG}" || true
exit 1
