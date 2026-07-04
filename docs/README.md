# LearnCosmos Docs

상태: active
최종 업데이트: 2026-07-04
목적: Codex와 기획자가 LearnCosmos 문서 기준을 빠르게 이해하기 위한 문서 루트 안내

---

## 1. 문서 루트 목적

이 디렉터리는 LearnCosmos의 최신 기준, 구현 전 계획, 의사결정, 세션 기록, 지식 허브, 버전 스냅샷을 분리 관리한다.

LearnCosmos는 LearnWeaver 기반 프로젝트이지만, 사용자-facing 브랜드와 신규 기획 기준은 LearnCosmos를 사용한다.

```text
LearnWeaver = 기반 프로젝트 / 레거시 코드·문서 기준
LearnCosmos = 사용자-facing 브랜드 / 최신 제품 기획 기준
```

---

## 2. Codex가 먼저 읽을 문서

Codex가 LearnCosmos 작업을 시작할 때는 아래 순서로 확인한다.

```text
1. docs/ai/QUICK_REF.md
2. docs/tasks/current-task.md
3. docs/README.md
4. docs/DOCUMENT_MANAGEMENT.md
5. docs/master/LearnCosmos_기획확정안_v1.md
6. docs/wiki/00_INDEX.md
7. 작업 성격에 맞는 docs/wiki/* 또는 docs/plans/*
8. 필요 시 docs/decisions/*, docs/versions/*
```

`docs/sessions/`는 기본 AI 컨텍스트로 사용하지 않는다. 과거 판단 근거가 필요하면 `docs/decisions/`를 먼저 확인한다.

서버 저장소 기준 권장 파일명은 다음과 같다.

```text
docs/README.md
docs/DOCUMENT_MANAGEMENT.md
docs/ai/QUICK_REF.md
docs/tasks/current-task.md
docs/wiki/00_INDEX.md
docs/wiki/architecture/LearnCosmos_Current_System_and_Library_Inventory_v1.md
```

---

## 3. 폴더 구조

```text
LearnCosmos/
└─ docs/
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

## 4. 각 폴더 역할

| 폴더 | 역할 | Codex 사용 기준 |
|---|---|---|
| `ai/` | AI 세션용 핵심 사실 요약 | 세션 시작 시 항상 먼저 확인 |
| `tasks/` | 현재 작업 범위와 목표 | 작업 착수 전 현재 상태 확인 |
| `master/` | 현재 제품 기준의 단일 원본 | 공식 정책·용어·UX 판단 시 우선 |
| `wiki/` | 반복 참고용 지식 허브 | 구조 이해, 코드맵, 운영 기준 탐색 |
| `wiki/architecture/` | 시스템 구조·라이브러리·데이터 구조 | 구현 전 기술 구조 확인 |
| `wiki/operations/` | 운영 방식·문서화 규칙·기획 비서 지침 | 작업 방식 확인 |
| `wiki/references/` | 문서맵·코드맵·참조 색인 | 원본 위치 탐색 |
| `plans/` | 구현 전 계획 문서 | Codex 작업 범위 잡기 |
| `decisions/` | 큰 의사결정 기록 | 왜 그렇게 결정했는지 확인 |
| `sessions/` | 논의/작업 세션 기록 | 필요 시 세부 경위 확인 |
| `versions/` | 과거 기준 스냅샷 | 이전 기준 비교 |

---

## 5. 최신 기준과 레거시 문서 관계

LearnCosmos 문서가 우선이다.

다만 실제 코드와 과거 구현은 LearnWeaver 이름을 사용한 문서와 경로에 남아 있을 수 있다.

판단 우선순위는 다음과 같다.

```text
1. 실제 코드 / migration / lockfile / go.mod
2. LearnCosmos master 문서
3. LearnCosmos wiki / plans
4. LearnWeaver 백업 문서
5. LearnWeaver 과거 통합 기획문서 / versions
```

문서와 코드가 충돌하면 코드와 migration을 확인한 뒤, Master와 wiki를 갱신한다.

---

## 6. 문서 생성 위치 원칙

새 문서는 루트에 만들지 않는다.

```text
AI 핵심 기준: docs/ai/
현재 작업 범위: docs/tasks/
확정 기준: docs/master/
구현 전 계획: docs/plans/
의사결정: docs/decisions/
논의 기록: docs/sessions/
지식 허브: docs/wiki/
과거 보관: docs/versions/
```

---

## 7. 서버 동기화 원칙

Google Drive에서 생성·수정한 LearnCosmos 문서는 서버 저장소의 같은 `docs/` 경로에도 동기화해야 한다.

Google Drive 문서는 편집·공유용이고, Codex와 배포 작업 기준은 서버 저장소 문서와 맞추는 것을 원칙으로 한다.

`git commit`, `git push`, Google Drive 문서 동기화는 사용자가 마감작업을 명시했을 때만 실행한다.

---

## 8. 작업 루틴

작업 루틴의 본체는 항상 다음 순서를 따른다.

```text
계획 -> 구현 -> 테스트 -> 빌드/재시작 -> 문서화
```

- 작업 전 범위와 영향 범위를 파악하고 단계별 계획을 세운다.
- 구현 후 테스트가 실패하면 다음 단계로 넘어가지 않고 수정 후 재테스트한다.
- 프론트 변경은 `npx tsc --noEmit`, 백엔드 변경은 `go build ./...` 또는 관련 테스트를 기본 검증으로 삼는다.
- 서비스 반영이 필요한 프론트 변경은 `deploy-frontend.sh`, 백엔드 변경은 `deploy-backend.sh`를 사용한다.
- 구현 상태나 정책 기준이 바뀌면 `master -> changelog -> sessions` 순서로 문서화한다.
- `git commit`, `git push`, Google Drive 문서 동기화는 사용자가 마감작업을 명시했을 때만 실행한다.

---

## 9. 현재 주요 기준 문서

```text
ai:
- QUICK_REF.md

tasks:
- current-task.md

master:
- LearnCosmos_기획확정안_v1.md

wiki:
- 00_INDEX.md

plans:
- LearnCosmos 메뉴 기획안 v1
- LearnCosmos 지점 학습 UX 기획안 v1

wiki/architecture:
- LearnCosmos 현재 시스템 구성 및 라이브러리 인벤토리 v1
```

---

## 10. 핵심 문장

LearnCosmos의 모든 기획과 구현 판단은 아래 문장으로 되돌아간다.

```text
콘텐츠를 학습 경험으로 바꾼다.
```
