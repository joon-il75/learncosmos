# LearnCosmos Changelog

상태: active
최종 업데이트: 2026-07-04

---

## 2026-07-04 — 기본 문서와 작업 루틴 이식

- Google Drive `LearnCosmos/docs` 기본 문서를 서버 `docs/`로 복사했다.
- 백업된 LearnWeaver 문서의 세션 시작 로딩 순서, 작업 루틴, 테스트/배포/문서화 게이트, 마감작업 조건을 LearnCosmos 이름과 경로에 맞게 이식했다.
- `docs/ai/QUICK_REF.md`, `docs/tasks/current-task.md`, `docs/wiki/00_INDEX.md`, `docs/changelog.md`를 추가했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-04 — Codex 응답 원칙 추가

- `AGENTS.md`에 SELFREFINE, ELI10, REDTEAM, /AUTOPROMPT, ALT3 응답 원칙을 추가했다.
- 상위 시스템/개발자 지침, 안전 정책, 도구 사용 규칙, 사용자의 명시 형식 요청이 우선한다는 적용 경계를 함께 명시했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-04 — 문서 규칙 허점 보강

- `current-task.md`의 현재 목표를 완료된 작업 기준에서 다음 작업 대기 상태로 갱신했다.
- session 기록 생성/번호 규칙, 고위험 작업 사용자 확인, 배포 스크립트 부재 시 보류 규칙, Google Drive 동기화 방식 미확정 시 보고 규칙을 보강했다.
- `docs/sessions/2026-07-04_session-01.md`에 이번 문서 정비 판단 근거를 기록했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-04 — 배포 스크립트 보강

- 백업 폴더의 `deploy-learncosmos-frontend.sh`와 `deploy-learncosmos-backend.sh`를 기준으로 새 레포 공식 `deploy-frontend.sh`, `deploy-backend.sh`를 생성했다.
- 기본 포트는 frontend `3001`, backend `8081`로 두어 기존 LearnWeaver `3000`/`8080` 포트를 건드리지 않게 했다.
- 현재 소스가 아직 없으면 명확한 오류로 중단하도록 frontend/backend 디렉터리와 핵심 파일 존재 검사를 추가했다.
- 문서만 변경했으며 실제 배포/재시작은 하지 않았다.

## 2026-07-04 — learncosmos.co.kr 전환 준비

- `learncosmos.co.kr` 안전 전환을 위한 nginx HTTP 템플릿과 HTTPS 템플릿을 `ops/nginx/`에 추가했다.
- root 권한 실행용 `scripts/ops/enable-learncosmos-domain.sh`를 추가했다.
- 스크립트는 DNS가 `175.45.200.17`로 전파되지 않으면 중단하고, LearnCosmos `3001/8081` 서비스 상태를 확인한 뒤 nginx 설정/인증서 발급을 진행하도록 작성했다.
- 현재 `learncosmos.co.kr` 루트 도메인은 아직 이전 IP `118.67.131.217`로 확인되어 실제 nginx/certbot 전환은 보류했다.

## 2026-07-04 - learncosmos.co.kr DNS 검사 보강

- whoisdomain 권한 네임서버(`ns1`~`ns4.whoisdomain.kr`)는 `learncosmos.co.kr`과 `www.learncosmos.co.kr` 모두 `175.45.200.17`을 응답하는 것으로 확인했다.
- 서버의 일반 resolver/getent 캐시는 루트 도메인을 이전 IP `118.67.131.217`로 응답할 수 있어, `scripts/ops/enable-learncosmos-domain.sh`의 DNS 준비 검사를 권한 네임서버 직접 조회 우선 방식으로 보강했다.
- 스크립트 문법 검사는 `bash -n scripts/ops/enable-learncosmos-domain.sh`로 통과했다.
- 실제 nginx/certbot 전환은 아직 root 권한으로 실행하지 않았다.

## 2026-07-04 - learncosmos.co.kr HTTPS 전환 완료

- root 권한으로 `scripts/ops/enable-learncosmos-domain.sh`를 실행해 `learncosmos.co.kr`과 `www.learncosmos.co.kr`의 Let's Encrypt 인증서 발급 및 nginx 적용을 완료했다.
- 인증서 경로는 `/etc/letsencrypt/live/learncosmos.co.kr/fullchain.pem`이며 만료일은 2026-10-02로 확인했다.
- 새 IP `175.45.200.17` 강제 지정 기준으로 루트와 `www` 모두 HTTPS `200 OK`를 확인했다.
- 서버 로컬 resolver는 루트 도메인을 이전 IP로 캐시할 수 있어, 전환 스크립트의 최종 smoke check도 `--resolve`로 새 IP를 고정하도록 보강했다.

## 2026-07-04 - 새 레포 소스 이식과 공개 도메인 정리

- 백업 폴더의 `frontend/`와 `backend/` 소스를 새 레포로 이식했다. 단, `.env*`, `.next`, `node_modules` Git 추적, runtime storage, 빌드 바이너리는 Git 제외 대상으로 정리했다.
- 프론트 공개 메타데이터, canonical, Open Graph, sitemap, robots.txt, llms.txt를 `https://learncosmos.co.kr` 기준으로 정리했다.
- 공개 화면 카피의 `LearnWeaver` 사용자 노출 표기를 주요 영역에서 `LearnCosmos`로 정리했다. 내부 이벤트 키, 레거시 자산 파일명, 개발 도메인 noindex 기준은 호환성 때문에 유지했다.
- `deploy-frontend.sh`의 포트 종료 로직을 `lsof` 단독에서 `fuser`/`ss` fallback 방식으로 보강했다.
- 새 레포 기준으로 프론트와 백엔드를 재배포했고, 프론트는 `/home/cosmos/LearnCosmos/frontend/.next/standalone`, 백엔드는 `/home/cosmos/LearnCosmos/backend`에서 실행 중임을 확인했다.
