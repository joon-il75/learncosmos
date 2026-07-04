#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="${PROJECT_ROOT:-/home/cosmos/LearnCosmos}"
DOMAIN="${DOMAIN:-learncosmos.co.kr}"
WWW_DOMAIN="${WWW_DOMAIN:-www.learncosmos.co.kr}"
EXPECTED_IP="${EXPECTED_IP:-175.45.200.17}"
NGINX_SITE="${NGINX_SITE:-/etc/nginx/sites-available/learncosmos}"
NGINX_LINK="${NGINX_LINK:-/etc/nginx/sites-enabled/learncosmos}"
HTTP_TEMPLATE="${HTTP_TEMPLATE:-${PROJECT_ROOT}/ops/nginx/learncosmos-http.conf}"
CERT_EMAIL="${CERT_EMAIL:-}"
BACKUP_DIR="${BACKUP_DIR:-/etc/nginx/sites-available/backup-learncosmos-$(date +%Y%m%d_%H%M%S)}"
AUTHORITATIVE_NS="${AUTHORITATIVE_NS:-ns1.whoisdomain.kr ns2.whoisdomain.kr ns3.whoisdomain.kr ns4.whoisdomain.kr}"

if [ "$(id -u)" != "0" ]; then
  echo "This script must run as root because it writes /etc/nginx and requests certificates." >&2
  exit 1
fi

if [ -z "${CERT_EMAIL}" ]; then
  echo "CERT_EMAIL is required, e.g. CERT_EMAIL=you@example.com $0" >&2
  exit 1
fi

resolve_ip() {
  if command -v dig >/dev/null 2>&1; then
    for ns in ${AUTHORITATIVE_NS}; do
      ip="$(dig @"${ns}" +short A "$1" | awk 'NF {print $1; exit}')"
      if [ -n "${ip}" ]; then
        echo "${ip}"
        return 0
      fi
    done
  fi

  getent hosts "$1" | awk '{print $1}' | head -n 1
}

root_ip="$(resolve_ip "${DOMAIN}" || true)"
www_ip="$(resolve_ip "${WWW_DOMAIN}" || true)"

if [ "${root_ip}" != "${EXPECTED_IP}" ]; then
  echo "DNS not ready for ${DOMAIN}: got '${root_ip}', expected '${EXPECTED_IP}'" >&2
  exit 1
fi
if [ "${www_ip}" != "${EXPECTED_IP}" ]; then
  echo "DNS not ready for ${WWW_DOMAIN}: got '${www_ip}', expected '${EXPECTED_IP}'" >&2
  exit 1
fi

if ! ss -ltn | grep -q ':3001 '; then
  echo "LearnCosmos frontend :3001 is not listening" >&2
  exit 1
fi
if ! ss -ltn | grep -q ':8081 '; then
  echo "LearnCosmos backend :8081 is not listening" >&2
  exit 1
fi

mkdir -p "${BACKUP_DIR}"
if [ -e "${NGINX_SITE}" ]; then
  cp -a "${NGINX_SITE}" "${BACKUP_DIR}/learncosmos"
fi
if [ -e /etc/nginx/sites-available/learnweavr ]; then
  cp -a /etc/nginx/sites-available/learnweavr "${BACKUP_DIR}/learnweavr"
fi

install -m 0644 "${HTTP_TEMPLATE}" "${NGINX_SITE}"
ln -sfn "${NGINX_SITE}" "${NGINX_LINK}"

nginx -t
systemctl reload nginx

certbot --nginx   -d "${DOMAIN}"   -d "${WWW_DOMAIN}"   --non-interactive   --agree-tos   --email "${CERT_EMAIL}"   --redirect

nginx -t
systemctl reload nginx

curl -fsSI --max-time 10 --resolve "${DOMAIN}:443:${EXPECTED_IP}" "https://${DOMAIN}/" >/dev/null
curl -fsS --max-time 10 --resolve "${DOMAIN}:443:${EXPECTED_IP}" "https://${DOMAIN}/api/health" >/dev/null || curl -fsS --max-time 10 --resolve "${DOMAIN}:443:${EXPECTED_IP}" "https://${DOMAIN}/api/v1/health" >/dev/null || true

echo "LearnCosmos domain enabled: https://${DOMAIN}"
echo "Nginx backups: ${BACKUP_DIR}"
