# LearnCosmos — Codex Agent Context

프로젝트: 취미 자기주도 학습 플랫폼
기반 프로젝트: LearnWeaver
작업 루트: `/home/cosmos/LearnCosmos`
Git: `git@github.com:joon-il75/learncosmos.git` (`main`)

---

## 세션 시작 시 반드시 읽어라

작업 시작 전 아래 파일을 순서대로 읽어라.

1. `docs/ai/QUICK_REF.md` — 핵심 운영 기준, 작업 루틴, 현재 저장소 기준
2. `docs/tasks/current-task.md` — 현재 진행 중인 작업 범위와 목표
3. `docs/README.md` — 문서 구조와 AI 컨텍스트 로딩 규칙
4. `docs/DOCUMENT_MANAGEMENT.md` — 문서 수정 원칙, 검증/배포/문서화 순서
5. `docs/master/LearnCosmos_기획확정안_v1.md` — 최신 제품 기준
6. `docs/wiki/00_INDEX.md` — 지식 허브 시작점

추가 규칙:

- `docs/sessions/`는 기본 AI 컨텍스트로 사용하지 않는다.
- 과거 판단 근거가 필요하면 `docs/decisions/`를 우선 본다.
- 작업 도메인이나 범위가 정해지면 관련 `docs/wiki/` 또는 `docs/plans/` 문서를 추가로 읽는다.
- 실제 코드와 문서가 충돌하면 코드, migration, lockfile, `go.mod`를 먼저 확인한 뒤 문서를 갱신한다.

---

## 작업 순서 (필수)

1. 단계별 계획 수립
   - 작업 범위와 영향 범위를 파악한다.
   - 구현 순서를 단계로 나누어 사용자에게 먼저 제시한다.
   - 고위험 영역, 배포/재시작, 데이터 변경, Git 원격/키 변경, 백업 삭제/이동은 사용자 확인 후 구현을 시작한다.
2. 구현
3. 테스트
   - 백엔드: `go build ./...` 또는 관련 테스트
   - 프론트: `npx tsc --noEmit`
   - 문서/공통: `git diff --check`
   - 실패 시 수정 후 같은 검증을 다시 실행한다.
   - 통과 전 다음 단계로 넘어가지 않는다.
4. 빌드/재시작
   - 프론트 변경 시 `deploy-frontend.sh`
   - 백엔드 변경 시 `deploy-backend.sh`
   - 둘 다 바뀌면 둘 다 실행한다.
   - 문서만 바뀌면 배포하지 않는다.
   - 배포 스크립트가 아직 없거나 LearnCosmos 경로와 맞지 않으면 임의 명령으로 대체하지 않고, 배포를 보류한 뒤 스크립트 정비 필요성을 결과에 명시한다.
5. 문서화
   - 구현 상태나 정책 기준이 바뀌면 기본 Master 파일 `docs/master/LearnCosmos_기획확정안_v1.md`를 갱신한다.
   - 변경 요약은 `docs/changelog.md`에 추가한다.
   - 판단 근거와 이슈는 해당 날짜 `docs/sessions/YYYY-MM-DD_session-XX.md` 문서에 기록한다.
   - 같은 날짜에 여러 세션 문서가 있으면 `session-01`, `session-02`처럼 다음 번호를 사용한다.
   - 주제별 지식으로 축적되어야 하면 관련 `docs/wiki/` 문서와 `docs/wiki/00_INDEX.md`를 함께 갱신한다.
6. 마감작업 (명시된 경우만)
   - `git commit`, `git push`
7. 마감작업 (명시된 경우만, 문서 변경 시)
   - Google Drive 문서 동기화
   - `rclone` 설정이 없으면 Google Drive 커넥터 기반 export/import 또는 수동 동기화 필요 상태로 보고하고, 임의 경로로 동기화하지 않는다.

추가 규칙:

- 작업 루틴의 본체인 `계획 -> 구현 -> 테스트 -> 빌드/재시작 -> 문서화`는 항상 지켜야 한다.
- `빌드/재시작`의 `필요하면`은 선택 배포를 뜻한다.
- `git commit`, `git push`, Google Drive 문서 동기화는 기본 마감작업으로 기억하되 자동 실행하지 않는다.
- 사용자가 `마감처리`, `마감작업`, `마감작업 실행`, `커밋`, `푸시`, `문서 동기화`를 명시했을 때만 해당 단계를 수행한다.

---

## GitHub 연결 확인

사용자가 GitHub/Git 연결 상태 확인을 요청하면 아래를 우선 확인한다.

1. 저장소 원격 확인: `git remote -v`
   - 기준 origin: `git@github.com:joon-il75/learncosmos.git`
2. 현재 브랜치/추적 상태 확인: `git status --short --branch`
3. SSH 인증 확인: `ssh -T git@github.com`
   - 정상 응답 예: `Hi joon-il75/learncosmos! You've successfully authenticated, but GitHub does not provide shell access.`
   - 위 메시지는 GitHub SSH 인증 성공을 뜻하며 shell 접속 실패가 아니다.

---

## 완료 보고 형식

작업 완료 시 다음 형식으로 보고한다.

1. 변경 요약
2. 수정 파일
3. 테스트 결과
4. 배포/재시작 여부
5. 문서 반영 여부
6. 남은 위험 또는 후속 작업

---

## 도메인별 추가 컨텍스트

현재 LearnCosmos 새 문서 체계에서는 도메인 문서가 아직 정리되지 않았다. 도메인 문서가 생기면 작업 도메인에 맞는 파일을 추가로 읽는다.

권장 위치:

- dashboard/Galaxy/Planet/Lumi → `docs/wiki/` 또는 `docs/ai/domain/dashboard.md`
- 인증/사용자/약관 → `docs/wiki/` 또는 `docs/ai/domain/auth.md`
- 콘텐츠/임베딩/검색 → `docs/wiki/` 또는 `docs/ai/domain/content.md`
- AI 포인트/BYOK/LLM → `docs/wiki/` 또는 `docs/ai/domain/ai-features.md`
- 슈퍼관리자/운영 → `docs/wiki/` 또는 `docs/ai/domain/admin.md`

도메인 문서가 아직 없으면 `docs/wiki/00_INDEX.md`, 관련 `docs/plans/`, `docs/master/LearnCosmos_기획확정안_v1.md`를 우선 확인한다.

---

## 응답 원칙

아래 원칙은 모든 답변에 기본 적용한다. 다만 시스템/개발자 지침, 안전 정책, 도구 사용 규칙, 사용자가 해당 턴에 명시한 형식 요청이 있으면 그 지침을 우선한다.

### SELFREFINE

- 최종 답변을 내기 전에 논리, 정확성, 일관성을 스스로 검토하고 개선한다.
- 더 좋은 표현이나 더 정확한 내용이 있으면 수정한 뒤 답한다.
- 객관적인 근거가 있는 내용은 가능한 한 근거를 함께 제시한다. 자료를 인용할 때는 핵심만 요약한다.

### ELI10

- 특별히 요청하지 않는 한 어려운 개념도 초등학생도 이해할 수 있을 정도로 쉽고 직관적으로 설명한다.
- 필요한 경우 비유와 예시를 사용한다.
- 단, 코드 변경/검증 결과처럼 정확성이 중요한 부분은 쉬운 말로 풀되 기술 사실을 흐리지 않는다.

### REDTEAM

- 사용자 의견에 무조건 동의하지 않는다.
- 틀린 부분, 위험한 부분, 비용이 커질 수 있는 선택은 명확하게 지적한다.
- 객관적인 사실과 근거를 우선한다.

### /AUTOPROMPT

- 사용자가 대충 설명하거나 키워드만 적어도 의도를 추론해 완성도 높은 결과물을 만든다.
- 정보가 조금 부족해도 합리적으로 보완해 진행한다.
- 정말 필요한 정보만 추가 질문한다.

### ALT3

- 선택지가 있는 질문이면 서로 다른 접근법이나 전략 3가지를 제시한다.
- 각각의 장단점을 함께 설명한다.
- 단, 사용자가 바로 실행을 요청했거나 선택지가 명확한 작업은 불필요하게 선택지를 늘리지 않는다.

---

## 현재 기준 요약

- LearnCosmos는 사용자-facing 브랜드와 신규 기획 기준이다.
- LearnWeaver는 기반 프로젝트와 레거시 코드/문서 기준으로 본다.
- 새 저장소는 `joon-il75/learncosmos`이며 기존 `learnweavr` 원격과 분리한다.
- 현재 서버 문서 기준은 Google Drive `LearnCosmos/docs`에서 복사한 기본 문서와 백업 LearnWeaver 작업 루틴을 이식한 기준이다.
- 배포 스크립트는 백업의 `deploy-learncosmos-frontend.sh` / `deploy-learncosmos-backend.sh`를 새 레포 공식 이름 `deploy-frontend.sh` / `deploy-backend.sh`로 보강한 기준이다.
- 기본 검증 포트는 frontend `3001`, backend `8081`이며, 기존 LearnWeaver 포트 `3000`/`8080`은 명시 override 없이는 사용하지 않는다.
- 최신 제품 기준은 `docs/master/LearnCosmos_기획확정안_v1.md`를 우선한다.
- 작업 기준은 `docs/ai/QUICK_REF.md`, `docs/tasks/current-task.md`, `docs/DOCUMENT_MANAGEMENT.md`를 우선한다.

---

## 고위험 영역

아래 영역은 변경 전 영향 범위를 명시해야 한다.

- OAuth / social_accounts / users.status
- BYOK / user_api_keys / ai_usage_events
- policy_documents / user_policy_consents
- ai_point_wallets / ai_point_transactions
- Object Storage / presigned URL / private attachments
- super-admin 권한 API
- explorer_regions / explorer_subregions / explorer_nodes
- course_point_* 학습 기록 테이블
- 배포 스크립트 / 서비스 포트 / runtime env
- Git 원격 / 브랜치 / 배포 키 / SSH 설정

---

## 절대 금지

- 커널 업그레이드
- `docs/` 외부에 기획문서 생성
- Master 문서에 긴 변경 이력 누적
- 아키텍처를 명시적 지시 없이 재설계
- 기존 컴포넌트 대신 새로 작성 (확장 우선)
- 기존 백업 디렉터리 `/home/cosmos/LearnCosmos.old-20260704-135329` 직접 수정
- 사용자가 명시하지 않은 `git commit`, `git push`, Google Drive 문서 동기화
