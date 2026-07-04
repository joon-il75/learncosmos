# LearnCosmos 문서 관리 설명서

상태: active
최종 업데이트: 2026-07-04
기준 프로젝트: LearnCosmos
기반 프로젝트: LearnWeaver
권장 서버 동기화 경로: docs/DOCUMENT_MANAGEMENT.md

---

## 1. 목적

이 문서는 LearnCosmos 프로젝트에서 기획 문서, 운영 문서, 구현 전 계획, 세션 기록, 버전 스냅샷을 어떻게 생성·수정·보관할지 정의한다.

목적은 다음과 같다.

```text
1. Codex가 문서 기준을 빠르게 이해하게 한다.
2. 현재 기준과 과거 이력을 분리한다.
3. Google Drive 문서와 서버 docs/ 문서의 불일치를 줄인다.
4. LearnWeaver 레거시 문서와 LearnCosmos 신규 문서의 관계를 명확히 한다.
5. 신규 기획을 구현 가능한 작업 단위로 전환한다.
6. 작업 루틴과 완료 게이트를 문서 기준으로 고정한다.
```

---

## 2. 기본 문서 구조

LearnCosmos의 기본 문서 구조는 다음과 같다.

```text
LearnCosmos/
└─ docs/
   ├─ README.md
   ├─ DOCUMENT_MANAGEMENT.md
   ├─ changelog.md
   ├─ ai/
   │  └─ QUICK_REF.md
   ├─ tasks/
   │  └─ current-task.md
   ├─ master/
   ├─ wiki/
   │  ├─ 00_INDEX.md
   │  ├─ architecture/
   │  ├─ operations/
   │  └─ references/
   ├─ plans/
   ├─ decisions/
   ├─ sessions/
   └─ versions/
```

---

## 3. 각 폴더의 역할

### 3.1 ai

`docs/ai/`는 Codex 세션 시작 시 빠르게 읽을 핵심 운영 기준을 둔다.

대표 문서:

```text
QUICK_REF.md
```

포함 항목:

```text
- 현재 저장소 / 원격 / 기본 브랜치
- 세션 시작 로딩 순서
- 작업 루틴과 테스트 게이트
- 배포 / 마감작업 조건
- 고위험 영역과 금지 사항
```

### 3.2 tasks

`docs/tasks/`는 현재 진행 중인 작업 범위와 목표를 둔다.

대표 문서:

```text
current-task.md
```

`current-task.md`는 항상 최신 상태를 유지한다. 완료된 작업이 누적되어 길어지면 별도 `docs/tasks/YYYY-MM-DD_<name>.md`로 분리한다.

### 3.3 master

`docs/master/`는 현재 제품 기준의 단일 원본이다.

포함 항목:

```text
- 확정된 제품 방향
- 확정된 용어 체계
- 확정된 UX 구조
- 확정된 학습 흐름
- 확정된 MVP 범위
- 공식 기준으로 삼을 정책
```

현재 대표 문서:

```text
LearnCosmos_기획확정안_v1.md
```

### 3.4 wiki

`docs/wiki/`는 반복 참고용 지식 허브다.

Wiki는 Master를 대체하지 않는다. Codex가 구조와 코드 위치, 운영 기준을 빠르게 찾도록 돕는 탐색·요약 레이어다.

권장 하위 구조:

```text
wiki/
├─ 00_INDEX.md
├─ architecture/
├─ operations/
└─ references/
```

### 3.5 wiki/architecture

시스템 구조, 프론트/백엔드 라이브러리, 데이터 모델, 인프라 기준을 정리한다.

### 3.6 wiki/operations

운영 규칙, 문서화 규칙, Codex 작업 방식 등을 정리한다.

### 3.7 wiki/references

문서맵, 코드맵, 참조 색인을 둔다.

### 3.8 plans

`docs/plans/`는 구현 전 계획 문서를 둔다.

### 3.9 decisions

`docs/decisions/`는 큰 의사결정과 판단 근거를 둔다.

### 3.10 sessions

`docs/sessions/`는 특정 날짜의 논의, 판단, 액션 아이템, 중간 결정 기록을 둔다.

`docs/sessions/`는 기본 AI 컨텍스트로 사용하지 않는다. 과거 판단 근거가 필요하면 먼저 `docs/decisions/`를 확인하고, 세부 경위가 필요할 때만 sessions를 읽는다.

### 3.11 versions

`docs/versions/`는 과거 기준 스냅샷을 보관한다.

---

## 4. Codex 문서 확인 순서

Codex가 작업을 시작할 때는 아래 순서로 문서를 읽는다.

```text
1. docs/ai/QUICK_REF.md
2. docs/tasks/current-task.md
3. docs/README.md
4. docs/DOCUMENT_MANAGEMENT.md
5. docs/master/LearnCosmos_기획확정안_v1.md
6. docs/wiki/00_INDEX.md
7. 작업 도메인에 맞는 wiki 문서
8. 작업 범위에 맞는 plans 문서
9. 필요 시 decisions / versions / sessions 확인
```

기술 구조 작업이면 다음을 추가로 확인한다.

```text
- docs/wiki/architecture/LearnCosmos_Current_System_and_Library_Inventory_v1.md
- docs/wiki/references/code-map.md
- docs/wiki/references/document-map.md
```

---

## 5. 문서 생성 위치 규칙

새 문서는 반드시 `docs/` 아래에 둔다.

```text
AI 핵심 기준 → docs/ai/
현재 작업 범위 → docs/tasks/
제품 공식 기준 → docs/master/
구현 전 계획 → docs/plans/
큰 의사결정 → docs/decisions/
작업/논의 기록 → docs/sessions/
지식 허브 → docs/wiki/
과거 기준 → docs/versions/
```

루트 폴더에 새 기획문서나 세션문서를 만들지 않는다.

---

## 6. Master와 Wiki의 역할 구분

### Master

Master는 최신 제품 기준이다.

```text
정책 판단, 용어 판단, 확정 UX 판단은 master를 우선한다.
```

### Wiki

Wiki는 반복 참고용 지식 허브다.

```text
구조 이해, 코드 위치 탐색, 운영 방식 확인은 wiki를 사용한다.
```

Wiki가 Master와 충돌하면 Master를 우선하고, Wiki를 갱신한다.

---

## 7. LearnWeaver와 LearnCosmos 문서 관계

LearnCosmos는 LearnWeaver 기반 프로젝트다.

```text
LearnWeaver = 기반 프로젝트 / 코드·레거시 문서 기준
LearnCosmos = 사용자-facing 브랜드 / 최신 제품 기획 기준
```

현재 코드, 경로, 일부 문서는 LearnWeaver 이름을 유지할 수 있다. 그러나 신규 기획, 화면 문구, 사용자-facing 문서 기준은 LearnCosmos로 통일한다.

판단 우선순위:

```text
1. 실제 코드 / migration / lockfile / go.mod
2. LearnCosmos master 문서
3. LearnCosmos wiki / plans
4. LearnWeaver 백업 문서
5. LearnWeaver 과거 통합 기획문서 / versions
```

---

## 8. Google Drive와 서버 저장소 동기화 원칙

Google Drive 문서를 생성하거나 수정하면, 서버 저장소의 대응 경로에도 동기화한다.

Google Drive 문서는 편집·공유용이다. Codex와 배포 기준은 서버 저장소의 `docs/`와 맞춰야 한다.

기본 절차:

```text
1. Google Drive 문서 생성 또는 수정
2. 문서 제목과 서버 저장 경로 결정
3. Google Drive 내용을 Markdown 형식으로 정리
4. 서버의 대응 경로에 복사
5. git diff로 변경 확인
6. 필요한 경우 changelog 또는 session 기록
7. 사용자에게 Google Drive 링크와 서버 반영 경로 보고
```

Google Drive 동기화는 문서 변경이 있고 사용자가 마감작업 또는 문서 동기화를 명시했을 때만 수행한다.

`rclone` 설정이 없으면 임의 경로로 동기화하지 않는다. 이 경우 Google Drive 커넥터 기반 export/import 또는 수동 동기화 필요 상태로 보고한다.

---

## 9. 작업 루틴

작업 명령을 받으면 아래 순서를 기본 규칙으로 따른다. 예외는 `git commit`, `git push`, Google Drive 동기화 같은 마감작업 명시 여부에만 적용된다.

```text
계획 -> 구현 -> 테스트 -> 빌드/재시작 -> 문서화
```

### 9.1 단계별 계획 수립

- 작업 범위와 영향 범위를 파악한다.
- 구현 순서를 단계로 나누고 필요한 경우 사용자에게 먼저 제시한다.
- 고위험 영역, 배포/재시작, 데이터 변경, Git 원격/키 변경, 백업 삭제/이동은 사용자 확인 후 구현한다.
- 고위험 영역이면 영향 범위를 명시하고 진행한다.

### 9.2 단계별 구현

- 각 단계를 순서대로 구현한다.
- 기존 컴포넌트와 구조를 우선 확장하고, 명시적 지시 없이 아키텍처를 재설계하지 않는다.
- `docs/` 외부에 기획문서를 만들지 않는다.

### 9.3 단계별 테스트

- 각 구현 단계 직후 해당 범위에서 검증을 수행한다.
- 백엔드: `go build ./...` 또는 관련 단위 테스트
- 프론트: `npx tsc --noEmit`
- 프론트 UI/레이아웃 변경은 모바일 우선으로 확인한다. 가장 좁은 320px viewport에서 구조, 텍스트, CTA, 주요 비주얼, body overflow를 먼저 안정화한 뒤 344/360/375px과 데스크톱으로 확장 확인한다.
- 테스트 실패 시 다음 단계로 넘어가지 않고 먼저 수정한다.
- 수정 후 반드시 테스트를 다시 실행한다. 통과할 때까지 수정과 테스트를 반복한다.

### 9.4 빌드 및 재시작

- 프론트 코드/UI 변경: `deploy-frontend.sh`
- 백엔드 API/서버 로직 변경: `deploy-backend.sh`
- 둘 다 바뀌면 둘 다 배포한다.
- 문서만 바뀌면 배포하지 않는다.
- 배포 또는 재시작을 수행하지 못한 경우에는 완료 보고에 명시한다.

### 9.5 문서화

- 구현 상태나 정책 기준이 바뀌면 기본 Master 파일 `docs/master/LearnCosmos_기획확정안_v1.md`를 갱신한다.
- 변경 요약은 `docs/changelog.md`에 추가한다.
- 판단 근거와 이슈는 해당 날짜 `docs/sessions/YYYY-MM-DD_session-XX.md` 문서에 기록한다.
- 같은 날짜에 여러 세션 문서가 있으면 `session-01`, `session-02`처럼 다음 번호를 사용한다.
- 작업 결과가 주제별 지식으로 축적되어야 하면 관련 `docs/wiki/` 문서와 필요 시 `docs/wiki/00_INDEX.md`를 갱신한다.
- 문서만 변경한 작업은 배포하지 않지만, 기준 변경이면 `Master`, `QUICK_REF`, `changelog`, `sessions`에 필요한 범위만 반영한다.

### 9.6 마감작업

`git commit`, `git push`, Google Drive 문서 동기화는 자동 실행하지 않는다.

사용자가 아래 표현을 명시했을 때만 해당 작업을 수행한다.

```text
마감처리
마감작업
마감작업 실행
커밋
푸시
문서 동기화
```

마감작업 명시 시 기본 순서는 다음과 같다.

```text
문서화 완료 -> git status 확인 -> stage -> commit -> push -> 필요한 경우 Google Drive 동기화
```

### 9.7 완료 보고 전 체크 게이트

- 코드/API/UI 변경이 있으면 해당 범위 테스트와 필요한 빌드/재시작을 먼저 끝낸다.
- 구현 상태나 정책 기준이 바뀌면 `master -> changelog -> sessions` 순서로 문서화한다.
- 지식 구조나 주제별 설명이 바뀌면 관련 `docs/wiki/` 문서를 확인한다.
- `git commit`, `git push`, Google Drive 동기화는 사용자가 명시했을 때만 실행한다.
- 수행하지 못한 항목은 완료 보고에 누락 사유와 남은 작업을 명시한다.

---

## 10. 배포 문서 필수 반영 사항

- 공식 프론트 배포 경로는 `/home/cosmos/LearnCosmos/deploy-frontend.sh`로 둔다.
- 공식 백엔드 배포 경로는 `/home/cosmos/LearnCosmos/deploy-backend.sh`로 둔다.
- 두 스크립트는 백업의 `deploy-learncosmos-frontend.sh` / `deploy-learncosmos-backend.sh`를 LearnCosmos 새 레포 기준으로 보강한 것이다.
- 기본 포트는 frontend `3001`, backend `8081`이며 기존 LearnWeaver 포트 `3000`/`8080`은 `ALLOW_LEARNCOSMOS_PROD_PORTS=1` 없이는 사용하지 않는다.
- 실제 소스 구조와 스크립트가 맞지 않는 경우 임의 배포 명령으로 대체하지 않는다. 배포를 보류하고 문서와 스크립트 기준을 함께 갱신하는 후속 작업으로 분리한다.
- 프론트가 Next.js `standalone` 배포 방식을 사용할 경우 `public/` 폴더는 자동 포함된다고 가정하지 않는다.
- 이미지/영상 깨짐 이슈가 있으면 먼저 `public/` 복사 누락 여부를 확인한다.

---

## 11. 문서 이름 규칙

### Google Drive 문서 제목

한글 제목을 사용할 수 있다.

예:

```text
LearnCosmos 메뉴 기획안 v1
LearnCosmos 지점 학습 UX 기획안 v1
LearnCosmos 현재 시스템 구성 및 라이브러리 인벤토리 v1
```

### 서버 저장소 파일명

서버 저장소에서는 영문 또는 혼합 파일명을 사용하되, 경로와 목적이 명확해야 한다.

예:

```text
docs/plans/LearnCosmos_메뉴_기획안_v1.md
docs/plans/LearnCosmos_지점_학습_UX_기획안_v1.md
docs/wiki/architecture/LearnCosmos_Current_System_and_Library_Inventory_v1.md
```

---

## 12. 세션 문서 규칙

세션 문서는 다음 형식을 권장한다.

```text
docs/sessions/YYYY-MM-DD_session-XX.md
```

내용 구조:

```text
# YYYY-MM-DD Session XX

## 1. 작업 배경
## 2. 확인한 기준 문서
## 3. 논의 및 판단
## 4. 변경/생성한 문서
## 5. 보류 항목
## 6. 다음 작업
```

---

## 13. Changelog 규칙

`docs/changelog.md`에는 길게 설명하지 않고 핵심 변경만 적는다.

```text
## YYYY-MM-DD — 변경 제목
- 무엇을 만들었는가
- 어떤 문서를 갱신했는가
- 서버 동기화 여부
```

---

## 14. Codex 작업 지시문 작성 규칙

Codex에게 넘길 작업 지시문은 다음 구조를 사용한다.

```text
# 작업 목적
# 현재 기준
# 변경 범위
# 변경하지 말아야 할 것
# 구현 단계
# 파일 / 컴포넌트 영향 범위
# API / DB 영향 범위
# 모바일 UX 고려사항
# 테스트 체크리스트
# 문서 반영 항목
```

추상적으로 쓰지 않고, 어느 화면/컴포넌트/API/DB에 영향을 주는지 명확히 쓴다.

---

## 15. 현재 주요 LearnCosmos 문서

```text
ai/
- QUICK_REF.md

tasks/
- current-task.md

master/
- LearnCosmos_기획확정안_v1.md

wiki/
- 00_INDEX.md

plans/
- LearnCosmos 메뉴 기획안 v1
- LearnCosmos 지점 학습 UX 기획안 v1

wiki/architecture/
- LearnCosmos 현재 시스템 구성 및 라이브러리 인벤토리 v1
```

---

## 16. 최종 원칙

LearnCosmos 문서 관리는 다음 흐름을 따른다.

```text
아이디어
→ 기준 문서 확인
→ 기획 정리
→ 문서 위치 결정
→ Google Drive 문서화
→ 서버 docs/ 동기화
→ Codex 작업 단위화
→ 계획
→ 구현
→ 테스트
→ 빌드/재시작
→ master / changelog / sessions 반영
```

LearnCosmos의 핵심 방향은 항상 다음 문장으로 되돌아간다.

```text
콘텐츠를 학습 경험으로 바꾼다.
```
