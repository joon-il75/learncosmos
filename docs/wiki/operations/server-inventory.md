# LearnCosmos Server Inventory

상태: active
최종 업데이트: 2026-07-11
기준 서버: 현재 LearnCosmos 운영 서버
저장 위치: `docs/wiki/operations/server-inventory.md`

---

## 1. 문서 목적

이 문서는 LearnCosmos 운영 서버의 사양, 런타임, 네트워크, 서비스 포트, 배포 기준을 한곳에 모아 백업, 복구, 장애 대응 때 빠르게 확인하기 위한 운영 인벤토리다.

민감 정보는 이 문서에 기록하지 않는다. DB 접속 문자열, API key, OAuth secret, JWT secret, Redis password 같은 값은 런타임 env 파일과 별도 비밀 관리 기준을 따른다.

---

## 2. 서버 기본 정보

| 항목 | 값 |
|---|---|
| 확인일 | 2026-07-11 |
| hostname | `learnweaver` |
| OS | Ubuntu 24.04.4 LTS (Noble Numbat) |
| Kernel | Linux 6.8.0-84-generic |
| Architecture | x86_64 |
| Virtualization | KVM full virtualization |
| Public IP | `175.45.200.17` |
| 기준 도메인 | `learncosmos.co.kr`, `www.learncosmos.co.kr` |
| 작업 루트 | `/home/cosmos/LearnCosmos` |
| Git remote | `git@github.com:joon-il75/learncosmos.git` |

주의: hostname은 아직 레거시 이름인 `learnweaver`다. 사용자-facing 브랜드와 저장소 기준은 LearnCosmos이며, hostname 변경은 서비스/SSH/운영 스크립트 영향이 있으므로 별도 계획 없이 변경하지 않는다.

---

## 3. 하드웨어 / 가상 자원

| 항목 | 값 |
|---|---|
| CPU | 2 vCPU |
| CPU 모델 | AMD EPYC 9454P 48-Core Processor |
| Socket | 1 |
| Core per socket | 2 |
| Thread per core | 1 |
| Memory | 7.8 GiB |
| Swap | 없음 |
| Root disk | 50G |
| Root filesystem | ext4 |
| Root mount | `/` |
| Root disk 사용량 | 35G / 50G, 74% 사용 |

디스크 구성:

```text
vda     50G disk
├─vda1   1M part
└─vda2  50G part ext4 /
```

운영 주의:

- swap이 없으므로 빌드 또는 이미지/영상 처리 작업에서 메모리 압박이 생길 수 있다.
- 루트 디스크 사용량이 70%를 넘은 상태라 백업 이미지, 대용량 로그, 업로드 파일을 만들기 전 여유 공간을 먼저 확인한다.
- 서버 전체 이미지 백업은 운영 중인 디스크를 직접 `dd`로 뜨기보다 VPS/클라우드 스냅샷 또는 서비스 중지 후 파일/DB 백업을 우선한다.

---

## 4. 주요 런타임 버전

| 항목 | 서버 확인값 |
|---|---|
| Node.js | v20.20.1 |
| npm | 10.8.2 |
| Go | go1.23.8 linux/amd64 |
| nginx | nginx/1.24.0 (Ubuntu) |
| PostgreSQL client | 16.11 |
| Redis server | 7.0.15 |

주의: 라이브러리 버전 기준은 실제 코드, lockfile, `go.mod`, migration을 우선한다. 이 문서의 런타임 버전과 `docs/wiki/architecture/LearnCosmos_Current_System_and_Library_Inventory_v1.md`의 요약값이 다르면 현재 서버와 코드 기준을 다시 대조한다.

---

## 5. 네트워크 / 포트

현재 확인된 listen 포트:

| 포트 | 바인딩 | 용도 |
|---|---|---|
| 22 | `0.0.0.0:22` | SSH |
| 80 | `0.0.0.0:80` | HTTP / HTTPS redirect / 인증서 발급 |
| 443 | `0.0.0.0:443` | HTTPS |
| 3001 | `0.0.0.0:3001` | LearnCosmos frontend |
| 8081 | `0.0.0.0:8081` | LearnCosmos backend |
| 5432 | `127.0.0.1:5432` | PostgreSQL |
| 6379 | `127.0.0.1:6379` | Redis |
| 9000 | `127.0.0.1:9000` | local internal service, 용도 확인 필요 |
| 53 | `127.0.0.53`, `127.0.0.54` | system resolver |
| 111 | `0.0.0.0:111` | rpcbind 계열, 필요성 확인 필요 |

서비스 상태 확인 결과:

```text
nginx: active
postgresql: active
redis-server: active
```

운영 주의:

- frontend `3001`, backend `8081`은 LearnCosmos 기준 포트다.
- 기존 LearnWeaver 포트 `3000`, `8080`은 명시 override 없이 사용하지 않는다.
- PostgreSQL과 Redis는 loopback에만 바인딩되어 있어 외부 직접 접속 대상이 아니다.
- `9000`, `111`은 현재 listen 중이나 LearnCosmos 공식 문서상 용도가 명확하지 않다. 정리 전에는 중지하지 말고 별도 영향 확인이 필요하다.

---

## 6. nginx / HTTPS 기준

기준 파일:

```text
ops/nginx/learncosmos-http.conf
ops/nginx/learncosmos-ssl.conf.template
scripts/ops/enable-learncosmos-domain.sh
```

운영 경로:

```text
/etc/nginx/sites-available/learncosmos
/etc/nginx/sites-enabled/learncosmos
/etc/letsencrypt/live/learncosmos.co.kr/
```

프록시 기준:

```text
/      -> http://127.0.0.1:3001
/api/  -> http://127.0.0.1:8081
```

`client_max_body_size` 기준은 600m이다.

---

## 7. 배포 / 실행 기준

프론트엔드:

```text
스크립트: ./deploy-frontend.sh
기본 포트: 3001
빌드: frontend에서 npm run build
실행: .next/standalone/server.js
로그: /tmp/learncosmos-frontend.log
빌드 로그: /tmp/learncosmos-frontend-build.log
```

백엔드:

```text
스크립트: ./deploy-backend.sh
기본 포트: 8081
빌드: backend에서 go build -o learncosmos-backend ./cmd/server
실행 env: /run/learncosmos/backend.env
bootstrap env: /etc/learncosmos/backend.env
로그: /tmp/learncosmos-backend.log
```

문서만 변경한 경우에는 배포 또는 재시작하지 않는다.

---

## 8. 백업 기준

권장 우선순위:

1. VPS/클라우드 제공자 스냅샷
2. 앱 파일 + DB dump + 업로드 파일 + nginx/systemd/env 경로 목록 백업
3. 서비스 중지 또는 스냅샷 기반 디스크 이미지

로컬 PC 백업 전 확인할 것:

```text
- 루트 디스크 여유 공간
- DB dump 저장 위치
- 업로드/첨부파일 저장 위치
- /etc/nginx/sites-available/learncosmos
- /etc/learncosmos/backend.env 존재 여부
- /run/learncosmos/backend.env는 런타임 파일이므로 영구 원본인지 확인 필요
- /etc/letsencrypt/live/learncosmos.co.kr 인증서 백업 필요 여부
```

운영 중인 DB가 있는 서버에서 디스크를 바로 이미지로 뜨면 파일 시스템과 DB 상태가 어긋날 수 있다. 서버 전체 이미지는 제공자 스냅샷을 우선하고, 파일 단위 백업은 DB dump를 먼저 만든 뒤 진행한다.

---

## 9. 재확인 명령

서버 사양을 갱신할 때 사용하는 확인 명령:

```bash
cat /etc/os-release
uname -a
lscpu
free -h
df -h
lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINTS,MODEL
ss -ltn
systemctl is-active nginx
systemctl is-active postgresql
systemctl is-active redis-server
node -v
npm -v
go version
nginx -v
psql --version
redis-server --version
getent hosts learncosmos.co.kr
getent hosts www.learncosmos.co.kr
```
