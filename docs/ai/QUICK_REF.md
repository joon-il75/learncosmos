# LearnCosmos Quick Reference

상태: active
최종 업데이트: 2026-07-04
목적: Codex 세션 시작 시 반드시 확인할 핵심 운영 기준

---

## 1. 프로젝트 기준

```text
프로젝트: LearnCosmos
기반 프로젝트: LearnWeaver
작업 루트: /home/cosmos/LearnCosmos
GitHub: git@github.com:joon-il75/learncosmos.git
기본 브랜치: main
```

LearnCosmos는 사용자-facing 브랜드와 신규 기획 기준이다. 실제 코드와 일부 레거시 문서는 LearnWeaver 이름을 유지할 수 있다.

판단 우선순위:

```text
1. 실제 코드 / migration / lockfile / go.mod
2. LearnCosmos master 문서
3. LearnCosmos wiki / plans
4. LearnWeaver 백업 문서
5. 과거 versions / sessions
```

---

## 2. 세션 시작 로딩 순서

작업 시작 전 아래 파일을 순서대로 읽는다.

```text
1. docs/ai/QUICK_REF.md
2. docs/tasks/current-task.md
3. docs/README.md
4. docs/DOCUMENT_MANAGEMENT.md
5. docs/master/LearnCosmos_기획확정안_v1.md
6. docs/wiki/00_INDEX.md
```

작업 도메인이나 범위가 정해지면 관련 `docs/wiki/` 또는 `docs/plans/` 문서를 추가로 읽는다.

`docs/sessions/`는 기본 AI 컨텍스트로 사용하지 않는다. 과거 판단 근거가 필요하면 `docs/decisions/`를 먼저 확인한다.

---

## 3. 작업 루틴

작업 루틴의 본체는 항상 아래 순서다.

```text
계획 -> 구현 -> 테스트 -> 빌드/재시작 -> 문서화
```

- 작업 범위와 영향 범위를 먼저 파악한다.
- 고위험 영역, 배포/재시작, 데이터 변경, Git 원격/키 변경, 백업 삭제/이동은 사용자 확인 후 구현한다.
- 구현 후 테스트가 실패하면 다음 단계로 넘어가지 않는다.
- 실패 원인을 수정한 뒤 같은 검증을 다시 실행한다.
- 배포 또는 재시작을 수행하지 못하면 완료 보고에 명시한다.
- 문서만 변경한 작업은 배포하지 않는다.

---

## 4. 기본 검증

```text
백엔드: go build ./...
프론트: npx tsc --noEmit
문서/공통: git diff --check
```

프론트 UI 변경은 모바일 우선으로 확인한다.

```text
320px -> 344px -> 360px -> 375px -> desktop
```

---

## 5. 배포/재시작 기준

```text
프론트 변경: deploy-frontend.sh (기본 포트 3001)
백엔드 변경: deploy-backend.sh (기본 포트 8081)
문서만 변경: 배포 없음
```

배포 스크립트는 백업의 LearnCosmos 전용 스크립트를 새 레포 공식 이름으로 보강한 것이다. 스크립트가 현재 소스 구조와 맞지 않으면 임의 명령으로 대체하지 않고 보류/수정 필요성을 결과에 명시한다.

---

## 6. 문서화 기준

구현 상태나 정책 기준이 바뀌면 아래 순서로 반영한다.

```text
1. docs/master/LearnCosmos_기획확정안_v1.md
2. docs/changelog.md
3. docs/sessions/YYYY-MM-DD_session-XX.md
```

작업 결과가 반복 참고용 지식이면 관련 `docs/wiki/`와 `docs/wiki/00_INDEX.md`도 갱신한다.

---

## 7. 마감작업 기준

`git commit`, `git push`, Google Drive 문서 동기화는 자동 실행하지 않는다.

사용자가 아래 표현을 명시했을 때만 수행한다.

```text
마감처리
마감작업
마감작업 실행
커밋
푸시
문서 동기화
```

마감작업 명시 시에는 문서화와 검증을 먼저 끝낸 뒤 `git status`로 범위를 확인한다. Google Drive 동기화는 rclone 설정이 없으면 커넥터 기반 동기화 또는 수동 동기화 필요 상태로 보고한다.

---

## 8. 현재 문서 상태

- Google Drive `LearnCosmos/docs` 기본 문서 7개를 서버 `docs/`로 복사했다.
- 백업된 기존 문서 위치: `/home/cosmos/LearnCosmos.old-20260704-135329/docs`
- 새 레포는 `joon-il75/learncosmos` 기준이며 기존 `learnweavr` 원격과 분리한다.
- 현재 문서 규칙은 백업의 LearnWeaver 작업 루틴을 LearnCosmos 이름/경로로 이식한 기준이다.

---

## 9. 금지/주의

- `docs/` 외부에 기획문서를 만들지 않는다.
- 명시적 지시 없이 아키텍처를 재설계하지 않는다.
- 기존 백업 문서를 직접 수정하지 않는다.
- 커널 업그레이드는 하지 않는다.
- 커밋/푸시/문서 동기화는 명시 요청 전 실행하지 않는다.
