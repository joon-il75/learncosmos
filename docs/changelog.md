# LearnCosmos Changelog

상태: active
최종 업데이트: 2026-07-11

---

## 2026-07-11 — 서버 사양 및 운영 인벤토리 문서화

- 현재 운영 서버의 OS, 커널, CPU, 메모리, 디스크, 런타임 버전, 포트, 주요 서비스 상태를 확인해 `docs/wiki/operations/server-inventory.md`에 정리했다.
- `learncosmos.co.kr` / `www.learncosmos.co.kr` 기준 IP, nginx 프록시 기준, 프론트 `3001`, 백엔드 `8081`, PostgreSQL/Redis loopback 바인딩 기준을 문서화했다.
- 루트 디스크 사용량 74%, swap 없음, `9000`/`111` 포트 용도 확인 필요 같은 운영상 주의점을 남겼다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-09 — WebGL 내부 CTA 외 문구 제거

- 공개 홈 WebGL 은하 맵 안에 있던 제목/설명 오버레이를 제거했다.
- 은하 맵 내부에는 새 학습별 생성 CTA만 남기고, CTA 위치를 상단 10px로 당겨 빈 공간을 줄였다.
- 직전 변경에서 `경로 만들기`로 바뀐 CTA 버튼 문구를 세계관 기준에 맞춰 `새 별`로 되돌렸다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — LearnWeaver 히어로 카피 적극 반영

- LearnWeaver 히어로의 `Lumi가 배우고 싶은 주제를 받아 목표를 정리하고 첫 학습탐험 경로를 만든다`는 메시지를 공개 홈 상단 카피에 적극 반영했다.
- 공개 홈의 Lumi Route 문구, 은하 맵 상단 제목/설명, 안내영상 설명, CTA 버튼 문구를 학습탐험 경로 중심으로 조정했다.
- CTA 버튼은 `새 별`에서 `경로 만들기`로 바꿔 사용자가 입력 후 무엇이 만들어지는지 더 분명하게 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 공개 홈 미니 푸터 추가

- LearnWeaver 푸터의 지원 링크 구성을 참고하되, 공개 홈 모바일 흐름에 맞춰 가벼운 미니 푸터를 추가했다.
- 푸터에는 LearnCosmos 브랜드명, 핵심 문장, 이용약관, 개인정보처리방침, 오픈소스, YouTube API 사용, 문의 정보를 배치했다.
- 기존 다중 컬럼 푸터는 사용하지 않고, 안내영상과 베타 테스트 섹션 뒤에 조용히 마감되는 구조로 정리했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — LearnWeaver 참고 베타 테스트 안내 섹션 개선

- LearnWeaver 백업의 `BetaNotice`, landing i18n, alpha 모집 페이지 문구를 참고해 공개 홈 베타 테스트 안내 섹션을 재구성했다.
- 초기 베타 단계에서는 화면, 추천 결과, AI 동작이 피드백에 따라 바뀔 수 있다는 안내를 명확히 했다.
- 베타에서 체험할 핵심 흐름을 체크 리스트로 정리하고, 섹션을 어두운 우주 카드와 골드/민트 포인트 톤으로 보강했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 베타 테스트 참가 안내와 설문 링크 섹션 추가

- 공개 홈 하단에 베타 테스트 참가 안내 섹션을 추가했다.
- `betaSurveyUrl` 상수를 추가해 향후 설문조사 URL을 넣으면 새 탭으로 연결되도록 준비했다.
- 설문 URL이 비어 있을 때는 비활성 상태 안내 문구를 보여주도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — WebGL 아래 YouTube 안내영상 섹션 추가

- 공개 홈 WebGL Galaxy Map 아래에 YouTube 안내영상 연결을 전제로 한 안내영상 섹션을 추가했다.
- 현재는 `guideVideoUrl` 상수를 비워두고, 향후 YouTube 영상 URL을 입력하면 새 탭으로 연결되는 구조로 준비했다.
- 영상 섹션은 WebGL 100vw 프레임 밖의 본문 폭에 배치해 모바일에서도 안정적으로 이어지도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 초소형 화면 CTA compact 적용

- 공개 홈 WebGL 상단 CTA가 344px 이하 화면에서 더 작고 촘촘하게 보이도록 compact breakpoint를 추가했다.
- 전체 화면 글씨는 유지하고, CTA 오버레이 내부의 입력창, 새 별 버튼, 추천 주제 칩만 높이, 패딩, 글자 크기를 한 단계 줄였다.
- viewport 비례 폰트 대신 `window.innerWidth <= 344` 기준의 구간형 처리로 가독성이 과도하게 작아지지 않도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — WebGL 은하 맵 100vw 확장

- 공개 홈 WebGL Galaxy Map을 모바일 화면에서 카드 폭이 아니라 viewport 전체 폭으로 보이도록 `100vw` 프레임으로 확장했다.
- Lumi 안내 카드는 최소 좌우 여백을 유지하고, 은하 맵만 화면 양끝까지 펼쳐 몰입감을 높였다.
- `100vw` 적용으로 생길 수 있는 가로 스크롤을 막기 위해 루트 영역에 `overflowX: hidden`을 적용했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 메시지와 CTA 상단 오버레이 재배치

- 공개 홈 WebGL Galaxy Map의 안내 메시지와 새 학습별 생성 CTA를 맵 상단 오버레이로 재배치했다.
- CTA 입력창, 새 별 버튼, 추천 주제 칩은 메시지 아래에 배치해 첫 시선에서 바로 학습별 생성을 시작할 수 있게 했다.
- 중심 펄스 레이어의 독립 발광 애니메이션을 고정해 은하 전체 루트 움직임과 한 덩어리로 보이도록 조정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — CTA WebGL 오버레이 배치 적용

- 공개 홈의 새 학습별 생성 CTA가 WebGL 은하 맵 아래 별도 영역으로 분리되어 보이던 구조를 맵 내부 하단 오버레이로 변경했다.
- 입력창, 새 별 버튼, 추천 주제 칩을 은하 위에 떠 있는 반투명 패널로 구성해 학습 지도와 생성 행동이 하나의 화면 안에서 이어지도록 했다.
- CTA 오버레이와 안내 문구가 겹치지 않도록 맵 높이와 안내 문구 위치를 조정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 은하 형태 유지형 동시 회전 적용

- 공개 홈 Galaxy Map에서 은하 레이어들이 서로 다른 속도로 회전해 시간이 지나면 은하 형태가 흐트러지는 문제를 줄였다.
- 은하 원반, 먼지띠, 외곽 팔, 중심 펄스, 하단 별무리의 개별 회전을 고정하고, `galaxyRoot` 전체만 아주 천천히 회전하도록 변경했다.
- 코스 별, 라벨, 추천 경로도 같은 `galaxyRoot` 안에서 함께 회전하므로 학습별이 은하와 분리되어 움직이지 않는다.
- 선택된 별로 줌인한 상태에서는 카메라 타깃이 회전 중인 별 위치를 따라가도록 보정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 공개 홈 은하 중앙 배치 보정

- 공개 홈 Galaxy Map의 은하가 화면 하단으로 치우쳐 보이는 문제를 줄이기 위해 Three.js `galaxyRoot`의 기본 y 오프셋을 상향 조정했다.
- 클릭한 코스 별로 줌인할 때도 같은 세로 오프셋을 반영해 선택 전환 위치가 어긋나지 않도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 중심 펄스 각도 90도 조정

- 공개 홈 Galaxy Map의 중심 펄스가 은하 원반 평면에서 90도 방향으로 나오도록 `upperPulse` 입자 좌표축을 조정했다.
- 기존에는 y축으로 길게 퍼졌던 펄스를 원반 법선 방향인 z축 중심으로 변경했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 3D 픽셀아트 은하 텍스처 적용

- 생성된 참고 이미지 방향에 맞춰 공개 홈 Galaxy Map의 은하 입자와 중심핵을 3D 픽셀아트 질감으로 조정했다.
- `createParticleTexture`를 저해상도 사각 픽셀 블록 기반 텍스처로 교체하고 `NearestFilter`를 적용했다.
- `createGalaxyCoreTexture`도 저해상도 픽셀 블록 중심핵으로 바꿔 따뜻한 중심 발광이 픽셀처럼 보이도록 했다.
- 픽셀 질감이 보이도록 은하 레이어별 Points 크기를 소폭 키웠다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-09 — 3D 은하 유지와 픽셀아트 코스 별 적용

- 공개 홈 Galaxy Map의 3D 은하 배경과 드래그 회전은 유지하고, 클릭 가능한 코스 별 텍스처만 픽셀아트 스프라이트풍으로 변경했다.
- `createStarTexture`를 16x16 그리드 기반 캔버스 드로잉으로 바꾸고, CanvasTexture에 `NearestFilter`를 적용해 픽셀 느낌이 유지되도록 했다.
- 픽셀 별이 과하게 튀지 않도록 후광, 궤도 링, 별 스프라이트 크기를 조금 줄였다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 1단계 은하 원반 Points 레이어 정교화

- 단계별 은하 개선의 1단계로 공개 홈 Galaxy Map의 은하 원반을 Points 레이어만으로 더 정교하게 조정했다.
- `innerDust` 레이어를 추가해 따뜻한 중심 먼지띠를 만들고, `blueRim` 레이어를 추가해 참고 이미지의 파란 외곽 원반감을 보강했다.
- 기존 Plane 없이 입자 레이어만 유지해 바닥 조각 문제가 재발하지 않도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 참고 이미지 기반 은하 별무리 재설계

- 공개 홈 Galaxy Map의 은하 별무리 디자인을 참고 이미지 기준으로 처음부터 다시 구성했다.
- 기존 3-arm 나선 생성 함수와 중심 펄스 함수를 제거하고, `createReferenceGalaxyLayer`로 disk, outerArm, core, upperPulse, lowerSparkles 레이어를 분리했다.
- 따뜻한 중심핵, 파란 외곽 원반, 보라 상단 흐름, 하단 별무리가 한 덩어리로 보이도록 색상과 입자 분포를 다시 설계했다.
- 평면 Plane 없이 Points 레이어만 사용해 바닥 조각 문제가 재발하지 않게 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 납작한 은하 원반과 중심 펄스 적용

- 공개 홈 Galaxy Map의 은하 y축 분포를 낮춰 더 납작한 원반형 입자 은하로 조정했다.
- 중심부에서 위아래로 뻗는 `corePulse` Points 레이어를 추가해 보라/푸른 펄스가 중심에서 솟는 느낌을 만들었다.
- 평면 Plane을 사용하지 않고 모든 효과를 Points 레이어로 구성해 바닥 조각이 다시 보이지 않게 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 코스 별을 은하 내부 입자로 통합

- 공개 홈 Galaxy Map에서 코스 별이 은하와 별도로 떠 보이지 않도록 배치 반경과 z 위치를 은하 입자층 안쪽으로 조정했다.
- 코스 별, 후광, 궤도 링, 라벨 크기를 줄이고 기본 라벨 투명도를 더 낮췄다.
- 추천 경로 연결선도 은하 내부 깊이에 맞추고 투명도를 낮춰 은하와 함께 회전하는 느낌을 강화했다.
- 별 크기를 줄인 뒤 클릭성이 너무 떨어지지 않도록 Raycaster Sprite threshold를 보강했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 은하 바닥 평면 조각 제거

- 공개 홈 Galaxy Map에서 기울어진 사각형 조각처럼 보이던 `galaxyDisc` Plane 레이어를 제거했다.
- 은하 원반감은 텍스처 Plane이 아니라 Points 입자 레이어의 밀도와 퍼짐으로 표현하도록 조정했다.
- 평면 제거로 줄어든 밀도는 외곽 haze와 주 나선팔 입자 수를 늘려 보완했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 검은 우주 배경을 딥 퍼플 우주로 전환

- 공개 홈 Galaxy Map의 검은 우주 배경을 딥 퍼플 우주 톤으로 변경했다.
- 별도 레이어처럼 보이던 CSS 먼지 패턴의 투명도를 낮춰 배경이 덜 거슬리도록 조정했다.
- 각진 평면 성운은 다시 추가하지 않고, 넓은 보라 그라데이션과 입자 은하 중심으로 유지했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 각진 성운 제거와 은하 중심 보라 발광 조정

- 공개 홈 Galaxy Map에서 평면처럼 보이던 보라색 성운 Plane과 CSS veil 레이어를 제거했다.
- 은하 중심으로 갈수록 입자가 보라색으로 발광하도록 나선 입자 색상 보간 로직을 조정했다.
- 중심핵과 은하 원반 텍스처도 따뜻한 노란빛보다 보라빛이 강하게 보이도록 보정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 어두운 3D 입자 은하와 드래그 회전 적용

- 첨부 참고 이미지처럼 공개 홈 Galaxy Map을 밝은 카드형 성운보다 어두운 3D 입자 은하에 가깝게 조정했다.
- Three.js 은하 레이어를 `galaxyRoot` 그룹으로 묶고 마우스/터치 드래그로 은하를 회전해 볼 수 있게 했다.
- 클릭과 드래그가 충돌하지 않도록 포인터 이동 거리가 일정 이상이면 별 선택 대신 은하 회전으로 처리한다.
- 기본 학습별 라벨 투명도를 더 낮추고, 안내 문구를 은하 회전 탐색 중심으로 바꿨다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-08 — 공개 홈 은하 단독 구현 강화

- 공개 홈 Three.js Galaxy Map에서 CSS성 가짜 안개보다 WebGL 내부 은하 레이어가 먼저 보이도록 조정했다.
- 밝은 성운 텍스처, 은하 원반 텍스처, 외곽 먼지층, 주 나선팔, 팔 sparkles, 중심 별무리를 분리해 은하 조감도 감각을 강화했다.
- 초기 조감도 카메라를 더 넓게 잡아 특정 학습별보다 은하 전체 구조가 먼저 읽히도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/`가 200 OK로 응답하는 것을 확인했다.

## 2026-07-07 — 공개 홈 은하 우선 구조와 코스 별 특수 효과 적용

- Three.js 장면을 학습별 중심 지도에서 은하 우선 구조로 재구성했다.
- 임의의 배경 별밭, 외곽 haze, 나선팔, 중심부 별무리, 중심핵이 먼저 은하를 구성하고, 학습별은 그 은하 안의 특별한 Course Star로 배치된다.
- 클릭 가능한 코스 별에는 일반 은하 별과 구분되는 후광, 궤도 링, 느린 맥동 효과를 추가했다.
- 선택된 코스 별은 링 회전과 광채가 강화되고, Raycaster 클릭 시 카메라가 해당 별로 줌인한다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈이 200 OK로 응답하는 것을 확인했다.

## 2026-07-07 — 공개 홈 은하 시각 강화

- Three.js Galaxy Map이 별 지도처럼 보이고 은하감이 약해, 중심핵, 3중 나선팔 입자층, 외곽 먼지층을 추가했다.
- 초기 조감도에서는 학습별 라벨 투명도를 낮춰 은하 구조가 먼저 보이도록 조정했다.
- 공개 홈 맵 높이와 성운 배경 대비를 보강해 320px 모바일에서도 은하의 중심과 팔이 더 분명하게 보이도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈이 200 OK로 응답하는 것을 확인했다.

## 2026-07-07 — 공개 홈 초기 은하 조감도 상태 적용

- 공개 홈 첫 진입 시 특정 학습별이 선택된 상태로 시작하지 않도록 `selectedLabel` 초기값을 제거했다.
- 초기 카메라는 특정 별 줌인이 아니라 은하 전체 조감도 위치에 머문다.
- 별을 클릭하면 Raycaster로 선택하고 해당 별로 줌인하며, 배경을 클릭하면 다시 조감도 상태로 돌아간다.
- 예시 주제 칩은 별 선택이 아니라 새 별 입력 보조 동작으로 분리했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈 HTML에서 `학습 은하를 먼저 둘러보세요` 문구가 포함된 것을 확인했다.

## 2026-07-07 — 공개 홈 LearnCosmos 로고 헤더 유지

- 공개 홈에서 브랜드 신호가 사라지지 않도록 얇은 LearnCosmos 로고 헤더를 복원했다.
- 헤더는 로고와 햄버거 버튼만 포함하며, 기존 탑 메뉴/좌측 메뉴/섹션 메뉴는 계속 제거 상태로 유지한다.
- Galaxy Map 본문은 로고 헤더 높이만큼 아래에서 시작하도록 여백을 조정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈 HTML에서 `LearnCosmos` 문구가 포함된 것을 확인했다.

## 2026-07-07 — 공개 홈 탑 메뉴 제거와 햄버거 메뉴 유지

- 공개 홈에서 `LearnerAppShell`을 제거해 고정 탑 메뉴, 좌측 메뉴, 섹션 메뉴, 상단 여백을 없앴다.
- 우주 맵 몰입감을 유지하되 다른 메뉴 진입이 막히지 않도록 플로팅 햄버거 버튼과 `MobileMainNavDrawer`를 유지했다.
- 햄버거 드로어에는 로그인 상태에 따라 Login 안내 또는 사용자 정보/Logout 블록을 표시하도록 `accountSlot`을 추가했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈 HTML에서 `메인 메뉴 열기` 문구가 포함된 것을 확인했다.

## 2026-07-07 — 공개 홈 Three.js/WebGL Galaxy Map 전환

- 공개 홈의 Canvas 2D 우주 효과를 Three.js/WebGL 기반 `ThreeCosmosMap`으로 전환했다.
- 학습별 샘플 데이터는 `THREE.Sprite` 별, `THREE.Points` 은하 입자, `THREE.LineSegments` 추천 항로로 렌더링한다.
- `Raycaster` 클릭으로 별을 선택하고, 선택한 별 방향으로 카메라가 줌인하도록 구현했다.
- 모바일 성능을 위해 renderer DPR을 제한하고, `prefers-reduced-motion` 환경에서는 반복 애니메이션을 줄인다.
- `three`와 `@types/three` 의존성을 추가했고 `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했다.

## 2026-07-07 — 공개 홈 canvas 기반 라이트 코스모스 효과 보강

- 공개 홈의 정적 CSS 별 지도에 실제 `canvas` 렌더링 레이어를 추가해 성운, 은하 먼지, 학습별, 추천 경로 연결선을 애니메이션으로 표현했다.
- `cosmos.codeghost.cloud`처럼 데이터 기반 별/선 시각화 방향을 따르되, LearnCosmos의 밝은 흰색/보라색/민트 색감을 유지하고 모바일 성능을 위해 별 입자 수를 제한했다.
- `prefers-reduced-motion` 환경에서는 반복 애니메이션을 중단하고 정적 캔버스 프레임만 그리도록 했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 공개 홈 HTML에서 `<canvas>` 요소가 포함된 것을 확인했다. Playwright 브라우저 바이너리가 서버에 없어 캔버스 픽셀 검사는 보류했다.

## 2026-07-07 — 공개 홈 기존 랜딩 제거와 넓은 모바일 기반 화면 적용

- 공개 홈 `www.learncosmos.co.kr/`에서 기존 긴 랜딩 여정, 소개 섹션, 알림 배너 렌더를 제거하고 `PublicCosmosMobileFirstView`를 단일 첫 화면으로 적용했다.
- 320px 모바일 기준으로 설계한 레이아웃을 최대 720px까지 넓혀 작은 모바일에서 출발하되 태블릿/데스크톱 폭에서도 같은 학습별 지도 경험이 유지되도록 조정했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 `https://www.learncosmos.co.kr/` HTML에서 `Lumi Route`, `하나의 학습별에서 시작하세요`, `새 별` 문구가 포함된 것을 확인했다.

## 2026-07-07 — 320px 모바일 Galaxy Map 첫 화면 MVP 구현

- `/dashboard` phone layout에서 기존 리스트 fallback 대신 `CosmosMobileFirstView`를 렌더링하도록 구현했다.
- 320px 모바일 기준으로 Lumi 오늘의 항로, 라이트 코스모스 별 지도, 선택 별 상세 패널, 이어서 탐험하기 CTA, 새 별 만들기 입력을 배치했다.
- 코스 하나를 별 하나로 표현하고, 코스 상태에 따라 별 색상과 광채를 다르게 표시했다.
- `npx tsc --noEmit`, `git diff --check`, `deploy-frontend.sh`를 통과했고 프론트는 3001 포트에서 재시작되었다.

## 2026-07-07 — 라이트 코스모스 색감과 밝은 보라색 성운 배경 기준 정리

- Galaxy Map의 전체 색감은 현재 `www.learncosmos.co.kr`의 밝은 LearnCosmos 톤을 유지하기로 했다.
- 첫 화면 배경은 검은 우주가 아니라 흰색과 밝은 보라색 성운 위에 은하가 얹힌 라이트 코스모스 톤으로 정리했다.
- 모바일에서는 성운과 은하 효과보다 선택한 별, 추천 항로, CTA 가독성을 우선한다는 기준을 메뉴 기획안과 Master에 반영했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-07 — 320px 모바일 햄버거 메뉴와 사용자 정보 기준 정리

- 320px 모바일에서 햄버거 메뉴는 전체 이동과 계정 확인을 위한 보조 서랍으로 유지하기로 했다.
- 학습 핵심 행동은 햄버거 안에 숨기지 않고 본문 또는 선택한 별 패널에 직접 노출하는 기준을 정리했다.
- 햄버거 메뉴 상단에는 사용자 정보 블록을 두고, Logout은 하단에 배치하는 원칙을 `docs/plans/LearnCosmos_메뉴_기획안_v1.md`와 Master에 반영했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

## 2026-07-07 — 세계관 용어와 Point/Content 구분 기준 정리

- LearnCosmos 세계관 기준을 `Course = 별 / 학습별`, `Lesson = 지역`, `Point = 학습지점`, `Content = 학습자료`로 정리했다.
- Point는 학습 행동이 일어나는 미션 단위, Content는 그 미션을 수행하기 위해 쓰는 학습자료로 구분했다.
- 완료 상태의 기본 기준은 Content가 아니라 Point에 둔다는 운영 규칙을 `docs/master/LearnCosmos_기획확정안_v1.md`에 반영했다.
- `docs/sessions/2026-07-07_session-01.md`에 논의 배경과 보류 항목을 기록했다.
- 문서만 변경했으므로 배포/재시작은 하지 않는다.

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

## 2026-07-06 - Google Drive 문서 동기화 루틴 확정

- `rclone` remote `Gdrive:` 기준으로 `Gdrive:LearnCosmos/docs`와 서버 `/home/cosmos/LearnCosmos/docs`를 동기화하는 표준 루틴을 확정했다.
- 삭제 위험을 줄이기 위해 서버 docs 백업 후 `Drive -> 서버 copy -> 서버 -> Drive sync dry-run -> 실제 sync -> 최종 dry-run` 순서로 진행한다.
- `docs/ai/QUICK_REF.md`, `docs/DOCUMENT_MANAGEMENT.md`, `docs/README.md`에 해당 운영 규칙을 반영했다.
