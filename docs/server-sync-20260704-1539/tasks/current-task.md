# Current Task

상태: active
최종 업데이트: 2026-07-04
작업 루트: `/home/cosmos/LearnCosmos`

---

## 1. 현재 목표

새 GitHub 저장소 `joon-il75/learncosmos`의 초기 문서 체계를 정비했다. 현재는 기본 문서, Codex 지침, 작업 루틴 보강이 완료된 상태이며, 다음 작업은 소스 구조 설계/복원 또는 백업 소스 선별 이식 대기 상태다.

---

## 2. 완료된 작업

- 기존 소스 백업 완료
  - 백업 디렉터리: `/home/cosmos/LearnCosmos.old-20260704-135329`
  - tar 백업: `/home/cosmos/archive/source-backups/20260704-135329/LearnCosmos-source-with-git.tar.gz`
- 새 저장소 초기화 완료
  - 원격: `git@github.com:joon-il75/learncosmos.git`
  - 브랜치: `main`
  - 초기 empty commit push 완료
- Google Drive `LearnCosmos/docs` 기본 문서 복사 완료
- 백업 문서의 문서 규칙과 작업 루틴 비교 완료
- LearnCosmos 기준 작업 루틴 이식 완료
- `AGENTS.md`에 Codex 응답 원칙 추가 완료
- 백업의 LearnCosmos 전용 배포 스크립트를 공식 `deploy-frontend.sh` / `deploy-backend.sh`로 보강 완료
- `learncosmos.co.kr` 안전 전환용 nginx 템플릿과 root 실행 스크립트 준비 완료
- `learncosmos.co.kr` 권한 네임서버 기준 DNS 검사 방식 보강 완료
- `learncosmos.co.kr` / `www.learncosmos.co.kr` HTTPS 인증서 발급 및 nginx 적용 완료
- 백업 `frontend/` / `backend/` 소스를 새 레포로 이식하고 새 레포 기준으로 재배포 완료
- 공개 도메인 메타데이터/canonical/sitemap/robots/llms를 `learncosmos.co.kr` 기준으로 정리 완료

---

## 3. 이번 작업 범위

- `docs/README.md`에 세션 시작 로딩 순서와 작업 루틴 요약 반영
- `docs/DOCUMENT_MANAGEMENT.md`에 백업 문서의 작업 게이트 이식
- `docs/ai/QUICK_REF.md` 신규 작성
- `docs/tasks/current-task.md` 신규 작성
- `docs/wiki/00_INDEX.md` 신규 작성
- `docs/changelog.md` 신규 작성

---

## 4. 작업 루틴

```text
계획 -> 구현 -> 테스트 -> 빌드/재시작 -> 문서화
```

- 문서만 변경한 작업은 배포하지 않는다.
- `git commit`, `git push`, Google Drive 문서 동기화는 사용자가 명시했을 때만 실행한다.

---

## 5. 검증 체크리스트

- [x] 문서 파일 생성/수정 확인
- [x] `rg`로 기존 LearnWeaver 경로가 새 규칙에 잘못 남지 않았는지 확인
- [x] `git diff --check`
- [x] `git status --short --branch`

---

## 6. 다음 작업 후보

- 새 문서 규칙을 기준으로 실제 소스 구조를 LearnCosmos로 설계/복원
- 기존 백업 소스에서 필요한 코드만 선별 이식
- 패키지명/내부 모듈명 등 남은 레거시 LearnWeaver 명칭의 단계적 정리 범위 결정
- Google Drive 동기화 방식을 rclone 또는 커넥터 기준으로 확정
- 도메인별 AI 컨텍스트 문서를 필요 시 생성
