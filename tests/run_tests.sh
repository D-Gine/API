#!/usr/bin/env bash
# Local CI. Lives in <project>/tests/run_tests.sh

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
API_DIR="${API_DIR:-$ROOT_DIR/api}"
OUT_DIR="$ROOT_DIR/tests/test_outputs"
rm -rf $OUT_DIR
mkdir -p "$OUT_DIR"

# Load the root .env
set -a
source "$ROOT_DIR/.env"
set +a

G=$'\e[0;32m'
R=$'\e[0;31m'
NC=$'\e[0m'

export DB_HOST=127.0.0.1

export GOTOOLCHAIN="${GOTOOLCHAIN:-go1.27.1}"
export PATH="$(go env GOPATH)/bin:$PATH"

# '--reset' flag to reset installations
GOBIN_DIR="$(go env GOPATH)/bin"
if [[ $* == '--re-install' ]]; then
    echo "Remove old versions"
    rm -f "$GOBIN_DIR/govulncheck" "$GOBIN_DIR/golangci-lint"
    echo "Re-installations"
    go install golang.org/x/vuln/cmd/govulncheck@latest
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
    echo "${G}DONE${NC}"
    echo ""
fi

cleanup() {
  [ -n "$API_PID" ] && kill "$API_PID" 2>/dev/null
  docker rm -f "$DB_CONTAINER" >/dev/null 2>&1
}
trap cleanup EXIT

run() {
  local title="$1" log="$2"
  shift 2
  echo "RUN $title"
  if "$@" > "$OUT_DIR/$log" 2>&1; then
    echo "${G}DONE${NC}"
  else
    echo "${R}FAILED${NC} (see tests/test_outputs/$log)"
  fi
  echo ""
}

run "Govulncheck" govulncheck.log bash -c "cd '$API_DIR' && '$GOBIN_DIR/govulncheck' ./..."

run "Golangci" golangci.log bash -c "cd '$API_DIR' && '$GOBIN_DIR/golangci-lint' run"

run "Build" build.log bash -c "cd '$API_DIR' && go mod download && go build -v ./..."
run "Unit tests" unit_tests.log bash -c "cd '$API_DIR' && go test -v ./..."

echo "RUN Start database"
docker rm -f "$DB_CONTAINER" >/dev/null 2>&1
docker run -d --name "$DB_CONTAINER" \
  -e PGPORT="$DB_PORT" \
  -e POSTGRES_DB="$DB_NAME" \
  -e POSTGRES_USER="$DB_USER" \
  -e POSTGRES_PASSWORD="$DB_PWD" \
  -p "$DB_PORT:$DB_PORT" \
  --health-cmd="pg_isready -U $DB_USER -d $DB_NAME" \
  --health-interval=2s --health-timeout=5s --health-retries=30 \
  ghcr.io/d-gine/database-test:latest >/dev/null || { echo "FAILED"; exit 1; }

for _ in $(seq 1 60); do
  [ "$(docker inspect -f '{{.State.Health.Status}}' "$DB_CONTAINER" 2>/dev/null)" = "healthy" ] && break
  sleep 2
done
echo "${G}DONE${NC}"
echo ""

run "Build API" build_api.log bash -c "cd '$API_DIR' && make build"

echo "RUN Start API"
(cd "$API_DIR" && ./api_bin > "$OUT_DIR/api.log" 2>&1) &
API_PID=$!
for _ in $(seq 1 30); do
  curl -sf "http://127.0.0.1:$API_PORT/api/health" >/dev/null && break
  sleep 2
done
echo "${G}DONE${NC}"
echo ""

run "Hoppscotch" hoppscotch.log npx -y @hoppscotch/cli@latest test \
  "$ROOT_DIR/tests/Hoppscotch/D-Gine API.json" \
  -e "$ROOT_DIR/tests/Hoppscotch/dengine api.json"

echo "ALL TESTS PASSED"
