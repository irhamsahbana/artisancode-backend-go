#!/usr/bin/env bash

set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://127.0.0.1:3939}"
UPLOAD_PATH="${UPLOAD_PATH:-/storage/upload-url}"
AUTH_TOKEN="${AUTH_TOKEN:-}"
FOLDER="${FOLDER:-attendance-face}"
CONTENT_TYPE="${CONTENT_TYPE:-image/jpeg}"
IS_PUBLIC="${IS_PUBLIC:-false}"
FILENAME="${FILENAME:-test-upload}"
FILE_PATH="${1:-${FILE_PATH:-}}"

usage() {
  cat <<'EOF'
Usage:
  AUTH_TOKEN=<bearer-token> ./scripts/test_storage_upload.sh [file_path]

Optional env:
  API_BASE_URL   Default: http://127.0.0.1:3939
  UPLOAD_PATH    Default: /storage/upload-url
  FOLDER         Default: attendance-face
  CONTENT_TYPE   Default: image/jpeg
  IS_PUBLIC      Default: false
  FILENAME       Default: test-upload
  FILE_PATH      Alternative to first positional argument

Examples:
  AUTH_TOKEN=xxx ./scripts/test_storage_upload.sh ./tmp/selfie.jpg
  AUTH_TOKEN=xxx FOLDER=attendance-face CONTENT_TYPE=image/png ./scripts/test_storage_upload.sh
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

cleanup() {
  if [[ -n "${TMP_FILE_PATH:-}" && -f "${TMP_FILE_PATH:-}" ]]; then
    rm -f "${TMP_FILE_PATH}"
  fi
  if [[ -n "${TMP_RESP_PATH:-}" && -f "${TMP_RESP_PATH:-}" ]]; then
    rm -f "${TMP_RESP_PATH}"
  fi
}

trap cleanup EXIT

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

require_cmd curl
require_cmd jq

if [[ -z "${AUTH_TOKEN}" ]]; then
  echo "AUTH_TOKEN is required" >&2
  usage >&2
  exit 1
fi

if [[ -z "${FILE_PATH}" ]]; then
  TMP_FILE_PATH="$(mktemp /tmp/storage-upload-test-XXXXXX.jpg)"
  printf 'storage upload smoke test\n%s\n' "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" >"${TMP_FILE_PATH}"
  FILE_PATH="${TMP_FILE_PATH}"
fi

if [[ ! -f "${FILE_PATH}" ]]; then
  echo "File not found: ${FILE_PATH}" >&2
  exit 1
fi

ORIGINAL_FILENAME="$(basename "${FILE_PATH}")"
REQUEST_BODY="$(jq -n \
  --arg filename "${FILENAME}" \
  --arg original_filename "${ORIGINAL_FILENAME}" \
  --arg content_type "${CONTENT_TYPE}" \
  --arg folder "${FOLDER}" \
  --argjson is_public "${IS_PUBLIC}" \
  '{
    filename: $filename,
    original_filename: $original_filename,
    content_type: $content_type,
    folder: $folder,
    is_public: $is_public
  }')"

echo "1. Requesting upload URL from ${API_BASE_URL}${UPLOAD_PATH}"

TMP_RESP_PATH="$(mktemp /tmp/storage-upload-response-XXXXXX.json)"
HTTP_STATUS="$(
  curl -sS \
    -o "${TMP_RESP_PATH}" \
    -w '%{http_code}' \
    -X POST "${API_BASE_URL}${UPLOAD_PATH}" \
    -H "Authorization: Bearer ${AUTH_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "${REQUEST_BODY}"
)"

if [[ "${HTTP_STATUS}" != "200" ]]; then
  echo "Failed to create upload URL. HTTP ${HTTP_STATUS}" >&2
  cat "${TMP_RESP_PATH}" >&2
  exit 1
fi

SUCCESS="$(jq -r '.success' "${TMP_RESP_PATH}")"
if [[ "${SUCCESS}" != "true" ]]; then
  echo "Backend returned unsuccessful response" >&2
  cat "${TMP_RESP_PATH}" >&2
  exit 1
fi

FILE_ID="$(jq -r '.data.file_id' "${TMP_RESP_PATH}")"
OBJECT_KEY="$(jq -r '.data.object_key' "${TMP_RESP_PATH}")"
UPLOAD_URL="$(jq -r '.data.upload_url' "${TMP_RESP_PATH}")"
METHOD="$(jq -r '.data.method // "PUT"' "${TMP_RESP_PATH}")"

echo "2. Uploading file to presigned URL"

UPLOAD_CMD=(
  curl -sS
  -o /dev/null
  -w '%{http_code}'
  -X "${METHOD}"
  --upload-file "${FILE_PATH}"
)

while IFS= read -r header; do
  if [[ -n "${header}" ]]; then
    UPLOAD_CMD+=(-H "${header}")
  fi
done < <(jq -r '.data.headers // {} | to_entries[] | "\(.key): \(.value)"' "${TMP_RESP_PATH}")

UPLOAD_CMD+=("${UPLOAD_URL}")
UPLOAD_STATUS="$("${UPLOAD_CMD[@]}")"

if [[ "${UPLOAD_STATUS}" != "200" && "${UPLOAD_STATUS}" != "201" && "${UPLOAD_STATUS}" != "204" ]]; then
  echo "Upload failed. HTTP ${UPLOAD_STATUS}" >&2
  exit 1
fi

cat <<EOF
Upload success.

file_id: ${FILE_ID}
object_key: ${OBJECT_KEY}
upload_method: ${METHOD}
uploaded_file: ${FILE_PATH}

You can now reuse file_id=${FILE_ID} in attendance requests that need selfie_file_id.
EOF
