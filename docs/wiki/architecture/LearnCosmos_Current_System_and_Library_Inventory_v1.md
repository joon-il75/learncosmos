# LearnCosmos 현재 시스템 구성 및 라이브러리 인벤토리 v1

상태: active
최종 업데이트: 2026-06-22
기준 프로젝트: LearnCosmos
기반 프로젝트: LearnWeaver
저장 위치: LearnCosmos/docs/wiki
권장 서버 동기화 경로: docs/wiki/architecture/LearnCosmos_Current_System_and_Library_Inventory_v1.md

---

## 1. 문서 목적

이 문서는 LearnCosmos / LearnWeaver의 현재 시스템 구성, 주요 라이브러리, 런타임 기준, 관련 구조 문서 위치를 한곳에서 확인하기 위한 인벤토리 문서다.

기존 기준은 여러 문서에 나뉘어 있다.

```text
- Frontend Architecture
- Backend Architecture
- Code Map
- Knowledge Status
- 2026-06-18 Session 01
- LearnWeaver Master
- Project Tree
```

이 문서는 위 문서들을 대체하지 않는다. 최신 원본은 실제 코드, lockfile, go.mod, migration, Master 문서다. 이 문서는 빠른 확인을 위한 LearnCosmos용 요약 인벤토리다.

---

## 2. 기준 원본

현재 시스템 구성 판단 시 우선 확인할 원본은 다음과 같다.

```text
프론트엔드:
- frontend/
- frontend/package.json
- frontend/package-lock.json
- docs/wiki/architecture/frontend.md

백엔드:
- backend/
- backend/go.mod
- backend/go.sum
- docs/wiki/architecture/backend.md

DB / Migration:
- backend/migrations/
- docs/wiki/architecture/data-model.md
- docs/master/LearnWeaver_Master.md

코드 위치 지도:
- docs/wiki/references/code-map.md

문서 구조:
- docs/wiki/references/document-map.md
- docs/wiki/knowledge-map.md
- docs/wiki/knowledge-status.md

최신 운영 확인 세션:
- docs/sessions/2026-06-18_session-01.md
```

---

## 3. 현재 운영 라이브러리 기준

2026-06-18 운영 기준 문서에서는 현재 쓰는 프론트/백엔드/DB 라이브러리를 운영 기준으로 문서화하고, `frontend/package.json`, `cd frontend && npm ls --depth=0`, `backend/go.mod`, DB migration / Quick Ref 기준을 대조한 것으로 정리되어 있다.

### 3.1 프론트엔드 운영 기준

```text
Next.js: 16.2.9
React: 19.2.4
React DOM: 19.2.4
TypeScript: 5.9.3
Tailwind CSS: 4.2.2
Tiptap: 3.22.5
Zustand: 5.0.12
qrcode.react: 4.2.0
Playwright: 1.61.0
Vitest: 4.1.2
ESLint: 9.39.4
```

판단 기준:

```text
운영 확인 기준은 cd frontend && npm ls --depth=0 설치값이다.
package.json은 semver 범위를 담고 있으므로, 실제 운영 확인 버전은 lockfile / install 기준으로 판단한다.
```

### 3.2 백엔드 운영 기준

```text
Go: 1.25.0
Gin: 1.10.0
pgx: 5.9.1
go-redis: 9.18.0
JWT v5: 5.3.1
google uuid: 1.6.0
godotenv: 1.5.1
pquerna/otp: 1.5.0
golang.org/x/crypto: 0.49.0
```

판단 기준:

```text
백엔드 운영 기준은 backend/go.mod, backend/go.sum, 실제 빌드 결과를 기준으로 판단한다.
```

---

## 4. 프론트엔드 시스템 구조

프론트 구조 기준 문서는 `docs/wiki/architecture/frontend.md`다.

이 문서는 Next.js frontend를 처음 볼 때 route, shell, 주요 화면 컴포넌트, 공통 lib 위치를 빠르게 찾기 위한 1차 지도다.

### 4.1 기준 원본

```text
- frontend/
- frontend/package.json
- frontend/package-lock.json
- docs/master/LearnWeaver_Master.md
```

### 4.2 주요 구조

```text
frontend/app/layout.tsx
  -> public routes: landing / platform / policy / auth
  -> dashboard/layout.tsx gate: refresh -> me -> language / consent -> alpha access
```

### 4.3 주요 영역

```text
public landing
- frontend/app/page.tsx
- frontend/app/en/page.tsx
- landing components

platform shell
- frontend/app/platform/*/page.tsx
- frontend/app/platform/PlatformShellPageClient.tsx

auth / language / policy
- frontend/app/(auth)/login/page.tsx
- frontend/app/language-setup/*
- frontend/app/agreements/page.tsx
- frontend/app/en/*

learner shell / settings
- frontend/app/dashboard/layout.tsx
- frontend/app/dashboard/settings/*
- frontend/components/dashboard/settings/*

dashboard / community / creator
- frontend/app/dashboard/page.tsx
- frontend/app/dashboard/community/*
- frontend/app/dashboard/creator/*

course draft / planet / point
- frontend/app/dashboard/course-drafts/[id]/*
- frontend/app/dashboard/planets/_shared/*
```

---

## 5. 백엔드 시스템 구조

백엔드 구조 기준 문서는 `docs/wiki/architecture/backend.md`다.

### 5.1 기준 원본

```text
- backend/
- backend/go.mod
- backend/go.sum
- docs/master/LearnWeaver_Master.md
```

### 5.2 주요 기술 기준

```text
Go + Gin API 서버
PostgreSQL / pgx 기반 데이터 접근
Redis / go-redis 기반 캐시 또는 큐/상태 관리
JWT 기반 인증
OTP 기반 관리자 인증 보조
OpenAI / 기타 LLM provider 연동 구조
```

### 5.3 확인 원칙

```text
구현 최신 원본은 backend 코드와 migration이다.
문서가 코드와 다를 경우, 먼저 코드와 migration을 확인하고 Master / wiki를 갱신한다.
```

---

## 6. 데이터 / DB 구조 기준

데이터 모델과 DB 구조는 다음 문서와 코드 기준을 함께 확인한다.

```text
- backend/migrations/
- docs/wiki/architecture/data-model.md
- docs/master/LearnWeaver_Master.md
- docs/wiki/references/code-map.md
```

핵심 기준:

```text
- migration이 DB 구조의 최종 기술 원본이다.
- Master는 현재 정책과 구현 상태의 기준이다.
- wiki는 빠른 탐색과 요약 레이어다.
```

---

## 7. 코드맵 기준

코드 위치 지도는 `docs/wiki/references/code-map.md`를 기준으로 한다.

이 문서는 주요 코드 위치를 wiki 도메인과 연결하기 위한 색인이다. 구현의 최신 원본은 코드와 migration이며, wiki는 탐색을 빠르게 하기 위한 안내 레이어다.

### 7.1 주요 프론트 코드맵

```text
public landing
platform shell
auth / language / policy
learner shell / settings
dashboard / community / creator
course draft / planet / point
diary stage
diary tree
diary map
planning editor
```

### 7.2 Explorer Diary / Point Workspace 관련 위치

```text
diary stage:
frontend/app/dashboard/planets/_shared/diary/PlanetDiaryStage.tsx

diary tree:
frontend/app/dashboard/planets/_shared/diary/DiaryTreePanel.tsx
components/explorer-plan/tree/*Row.tsx

diary map:
frontend/app/dashboard/planets/_shared/diary/DiaryMapPanel.tsx
course-drafts/[id]/sections/planning/map/*

planning editor:
frontend/app/dashboard/course-drafts/[id]/sections/PlanningSection.tsx
components/explorer-plan/useExplorerPlan*.ts
```

---

## 8. 프로젝트 트리 기준

프로젝트 전체 구조 확인에는 `PROJECT_TREE.md`를 사용한다.

```text
/home/weaver/learnweaver
```

이 문서는 특정 시점의 실제 프로젝트 파일 트리를 기록한다. 파일 존재 여부, 주요 경로 확인, 문서/코드 구조 파악에 사용한다.

---

## 9. 과거 기술 스택 문서와의 관계

과거 통합 기획 문서에도 기술 스택이 정리되어 있다.

```text
LearnWeaver_기획문서_통합_v9_수정_v1.md 등
```

해당 문서에는 Naver Cloud Platform, Ubuntu 24.04 LTS 등 초기 인프라 및 기술 스택 기준이 포함되어 있다.

다만 현재 운영 라이브러리 기준은 2026-06-18 기준의 frontend / backend architecture wiki와 session 문서를 우선한다.

---

## 10. 현재 인프라 / 운영 기준 요약

현재까지 확인된 운영 인프라 기준은 다음과 같다.

```text
인프라:
- Naver Cloud Platform
- Ubuntu 24.04 LTS

프론트:
- Next.js 16 계열
- React 19 계열
- TypeScript
- Tailwind CSS

백엔드:
- Go
- Gin
- PostgreSQL / pgx
- Redis / go-redis

AI / 검색:
- OpenAI 중심 LLM / embedding 운영 기준
- BYOK 구조
- 내부 검색 + 외부 검색 보강 구조

문서 / 작업:
- docs/master: 최신 정책 원본
- docs/wiki: 탐색 / 요약 / 도메인 지식
- docs/plans: 구현 전 계획
- docs/sessions: 세션 기록
- docs/changelog.md: 변경 요약
```

---

## 11. 보강 필요

현재 이 문서는 기존 LearnWeaver 기준을 LearnCosmos 문서 체계 안에 모은 1차 인벤토리다.

후속으로 보강할 항목은 다음과 같다.

```text
1. 실제 서버 저장소의 frontend/package-lock.json 기준 재확인
2. backend/go.mod / go.sum 기준 재확인
3. DB migration 최신 번호와 주요 테이블 요약 추가
4. AI / BYOK / 포인트 / 검색 엔진 라이브러리 별도 섹션 확장
5. 배포 스크립트와 운영 명령 요약 추가
6. LearnCosmos 명칭 기준으로 LearnWeaver 레거시 문서와 연결 관계 정리
```

---

## 12. 최종 정리

LearnCosmos / LearnWeaver에는 현재 시스템 구성과 라이브러리 기준이 이미 여러 문서에 나뉘어 정리되어 있다.

이 문서는 그 기준을 LearnCosmos 폴더 안에서 빠르게 확인하기 위한 인벤토리다.

```text
최신 판단 우선순위:
1. 실제 코드 / lockfile / go.mod / migration
2. docs/master/LearnWeaver_Master.md
3. docs/wiki/architecture/frontend.md
4. docs/wiki/architecture/backend.md
5. docs/wiki/references/code-map.md
6. docs/sessions/2026-06-18_session-01.md
7. 과거 통합 기획 문서
```
