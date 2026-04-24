#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"
WAILS_EXE="${WAILS_EXE:-}"

if [[ -z "${WAILS_EXE}" ]]; then
  if command -v wails.exe >/dev/null 2>&1; then
    WAILS_EXE="$(command -v wails.exe)"
  else
    for candidate in /mnt/c/Users/*/go/bin/wails.exe; do
      if [[ -x "${candidate}" ]]; then
        WAILS_EXE="${candidate}"
        break
      fi
    done
  fi
fi

if [[ ! -f "${ENV_FILE}" ]]; then
  echo ".env not found at ${ENV_FILE}" >&2
  exit 1
fi

if [[ ! -x "${WAILS_EXE}" ]]; then
  echo "wails.exe not found at ${WAILS_EXE}" >&2
  exit 1
fi

set -a
source "${ENV_FILE}"
set +a

required_vars=(
  DB_HOST
  DB_PORT
  DB_USER
  DB_PASSWORD
  DB_NAME
)

for var_name in "${required_vars[@]}"; do
  if [[ -z "${!var_name:-}" ]]; then
    echo "Required variable ${var_name} is empty in ${ENV_FILE}" >&2
    exit 1
  fi
done

BOOTSTRAP_ADMIN_USERNAME="${BOOTSTRAP_ADMIN_USERNAME:-}"
BOOTSTRAP_ADMIN_PASSWORD="${BOOTSTRAP_ADMIN_PASSWORD:-}"
BOOTSTRAP_ADMIN_FULL_NAME="${BOOTSTRAP_ADMIN_FULL_NAME:-}"

ldflags=(
  "-X" "inventory-app/backend/internal/db.BuildDBHost=${DB_HOST}"
  "-X" "inventory-app/backend/internal/db.BuildDBPort=${DB_PORT}"
  "-X" "inventory-app/backend/internal/db.BuildDBUser=${DB_USER}"
  "-X" "inventory-app/backend/internal/db.BuildDBPassword=${DB_PASSWORD}"
  "-X" "inventory-app/backend/internal/db.BuildDBName=${DB_NAME}"
)

if [[ -n "${BOOTSTRAP_ADMIN_USERNAME}" ]]; then
  ldflags+=("-X" "inventory-app/backend/internal/db.BuildBootstrapAdminUsername=${BOOTSTRAP_ADMIN_USERNAME}")
fi

if [[ -n "${BOOTSTRAP_ADMIN_PASSWORD}" ]]; then
  ldflags+=("-X" "inventory-app/backend/internal/db.BuildBootstrapAdminPassword=${BOOTSTRAP_ADMIN_PASSWORD}")
fi

if [[ -n "${BOOTSTRAP_ADMIN_FULL_NAME}" ]]; then
  ldflags+=("-X" "inventory-app/backend/internal/db.BuildBootstrapAdminFullName=${BOOTSTRAP_ADMIN_FULL_NAME}")
fi

cd "${ROOT_DIR}"
"${WAILS_EXE}" build -ldflags "${ldflags[*]}"
