#!/usr/bin/env bash
# LearnCosmos backend deploy script
# Derived from the backed-up deploy-learncosmos-backend.sh and adapted for /home/cosmos/LearnCosmos.
# Defaults to port 8081 to avoid touching the legacy LearnWeaver backend on 8080.

set -euo pipefail

PROJECT_ROOT="${PROJECT_ROOT:-/home/cosmos/LearnCosmos}"
BACKEND_DIR="${PROJECT_ROOT}/backend"
BACKEND_ENV_FILE="${BACKEND_ENV_FILE:-/run/learncosmos/backend.env}"
BOOTSTRAP_ENV_FILE="${BOOTSTRAP_ENV_FILE:-/etc/learncosmos/backend.env}"
BACKEND_BINARY="${BACKEND_BINARY:-learncosmos-backend}"
BACKEND_PORT="${BACKEND_PORT:-8081}"
BACKEND_LOG="${BACKEND_LOG:-/tmp/learncosmos-backend.log}"
GOFLAGS="${GOFLAGS:--buildvcs=false}"
LEARNWEAVER_PROJECT_ROOT="${LEARNWEAVER_PROJECT_ROOT:-${PROJECT_ROOT}}"
ALLOW_LEARNWEAVER_ENV="${ALLOW_LEARNWEAVER_ENV:-0}"
ALLOW_LEARNCOSMOS_PROD_PORTS="${ALLOW_LEARNCOSMOS_PROD_PORTS:-0}"

if [ "${BACKEND_PORT}" = "8080" ] && [ "${ALLOW_LEARNCOSMOS_PROD_PORTS}" != "1" ]; then
  echo "Refusing to use legacy LearnWeaver backend port 8080. Use BACKEND_PORT=8081 for LearnCosmos verification." >&2
  exit 1
fi

if [ ! -d "${BACKEND_DIR}" ]; then
  echo "Backend directory not found: ${BACKEND_DIR}" >&2
  echo "Restore or create the backend source before running this deploy script." >&2
  exit 1
fi

if [ ! -f "${BACKEND_DIR}/go.mod" ]; then
  echo "backend/go.mod not found: ${BACKEND_DIR}/go.mod" >&2
  exit 1
fi

if [ ! -r "${BACKEND_ENV_FILE}" ]; then
  echo "Runtime env file is not readable: ${BACKEND_ENV_FILE}" >&2
  if [ -x "${PROJECT_ROOT}/scripts/secrets/materialize-runtime-backend-env.sh" ]; then
    echo "Materializing runtime env from encrypted secrets..."
    PROJECT_ROOT="${PROJECT_ROOT}" \
      BOOTSTRAP_ENV_FILE="${BOOTSTRAP_ENV_FILE}" \
      RUNTIME_DIR="$(dirname "${BACKEND_ENV_FILE}")" \
      RUNTIME_ENV_FILE="${BACKEND_ENV_FILE}" \
      "${PROJECT_ROOT}/scripts/secrets/materialize-runtime-backend-env.sh"
  else
    echo "Env materializer not found. Create ${BACKEND_ENV_FILE} or restore scripts/secrets before deploying." >&2
    exit 1
  fi
fi

if [ ! -r "${BACKEND_ENV_FILE}" ]; then
  echo "Runtime env file is not readable after materialize: ${BACKEND_ENV_FILE}" >&2
  exit 1
fi

if [ "${ALLOW_LEARNWEAVER_ENV}" != "1" ] && grep -Eq '^DATABASE_URL=.*learnweaver' "${BACKEND_ENV_FILE}"; then
  echo "Refusing to start: DATABASE_URL appears to reference learnweaver instead of learncosmos." >&2
  exit 1
fi

if [ "${ALLOW_LEARNWEAVER_ENV}" != "1" ] && grep -Eq '^REDIS_URL=.*:6379/0($|[[:space:]]|\?)' "${BACKEND_ENV_FILE}"; then
  echo "Refusing to start: REDIS_URL appears to use the existing Redis DB 0." >&2
  exit 1
fi

cd "${BACKEND_DIR}"
echo "Building LearnCosmos backend..."
GOFLAGS="${GOFLAGS}" go build -o "${BACKEND_BINARY}" ./cmd/server

echo "Restarting LearnCosmos backend on port ${BACKEND_PORT}..."
PIDS=$(lsof -t -iTCP:"${BACKEND_PORT}" -sTCP:LISTEN 2>/dev/null || true)
if [ -n "${PIDS:-}" ]; then
  kill ${PIDS} || true
fi

export BACKEND_ENV_FILE BACKEND_DIR BACKEND_BINARY BACKEND_PORT PROJECT_ROOT LEARNWEAVER_PROJECT_ROOT
setsid -f bash -lc 'set -euo pipefail; set -a; source "${BACKEND_ENV_FILE}"; set +a; export PORT="${BACKEND_PORT}"; cd "${BACKEND_DIR}"; exec "./${BACKEND_BINARY}"' >"${BACKEND_LOG}" 2>&1

echo "LearnCosmos backend started. Checking port ${BACKEND_PORT}..."
for _ in 1 2 3 4 5; do
  if ss -ltn | grep -q ":${BACKEND_PORT} "; then
    tail -n 20 "${BACKEND_LOG}" || true
    echo "LearnCosmos backend listening on :${BACKEND_PORT}"
    exit 0
  fi
  sleep 1
done

echo "LearnCosmos backend did not start on :${BACKEND_PORT}" >&2
tail -n 80 "${BACKEND_LOG}" || true
exit 1
