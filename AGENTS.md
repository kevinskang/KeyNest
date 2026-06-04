# KeyNest — Agent & Developer Onboarding Guide

> 이 파일은 AI 에이전트 및 새로 합류하는 개발자가 프로젝트를 빠르게 파악하고 작업을 이어갈 수 있도록 작성된 가이드입니다.  
> 작업 시작 전 반드시 이 파일과 `docs/progress.md`를 먼저 읽으십시오.

---

## 1. 프로젝트 개요

**KeyNest**는 GitHub Personal Token, OpenAI API Key 등 개인 키·토큰을 로컬에 암호화하여 관리하는 **Windows standalone 데스크톱 앱**입니다.

- 인터넷 연결 불필요 — 완전 로컬 동작
- 데이터는 `keynest.db` (SQLite3)에 AES-256-GCM 암호화 후 저장
- 상세 요구사항 → `docs/Requirements.md`

---

## 2. 기술 스택

| 역할 | 기술 |
|------|------|
| 데스크톱 프레임워크 | **Wails v2.12.0** (Go + Web UI → 단일 exe) |
| 백엔드 | **Go 1.23** |
| 프론트엔드 | **Vue 3 + TypeScript** (Vite 빌드) |
| 상태 관리 | Pinia |
| 라우터 | Vue Router 4 |
| 데이터베이스 | SQLite3 (`mattn/go-sqlite3`) |
| 암호화 | AES-256-GCM + PBKDF2-SHA256 (Go 표준 라이브러리) |
| 비밀번호 해시 | bcrypt cost=12 (`golang.org/x/crypto`) |

---

## 3. 개발 환경 설정 및 실행

### 사전 요구사항

```powershell
# Go 1.23+
go version

# Node.js 18+ (프론트엔드 빌드용)
node --version

# Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails version

# CGO 필요 (mattn/go-sqlite3) — GCC 설치 필요
# Windows: https://www.msys2.org/ 또는 tdm-gcc
gcc --version
```

### 개발 모드 실행 (Hot Reload)

```powershell
# 프로젝트 루트에서
wails dev
```

### 프로덕션 빌드

```powershell
wails build
# 출력: build/bin/KeyNest.exe
```

### 프론트엔드만 단독 개발

```powershell
cd frontend
npm install
npm run dev
# http://localhost:5173 — Wails mock 없이 UI 확인
```

---

## 4. 디렉토리 구조

```
KeyNest/
├── AGENTS.md                   ← 이 파일
├── main.go                     # Wails 앱 초기화 진입점
├── app.go                      # Wails 바인딩 메서드 (프론트엔드 호출 엔드포인트)
├── go.mod / go.sum
│
├── internal/                   # Go 내부 패키지 (외부 임포트 금지)
│   ├── database/
│   │   └── db.go               # SQLite 연결, PRAGMA 설정, 마이그레이션
│   ├── models/
│   │   ├── user.go             # User 구조체
│   │   └── apikey.go           # APIKey 구조체 + DTO
│   ├── repository/
│   │   ├── user_repo.go        # users 테이블 CRUD
│   │   └── apikey_repo.go      # api_keys 테이블 CRUD
│   ├── service/
│   │   ├── auth_service.go     # 회원가입, 로그인, 로그아웃
│   │   └── apikey_service.go   # 키 CRUD + 암복호화 위임
│   └── crypto/
│       └── crypto.go           # AES-256-GCM 암복호화, PBKDF2 키 파생
│
├── frontend/
│   └── src/
│       ├── views/
│       │   ├── LoginView.vue       # 로그인 화면
│       │   ├── RegisterView.vue    # 회원가입 화면
│       │   └── MainView.vue        # 메인(키 목록) 화면
│       ├── components/
│       │   ├── KeyTable.vue        # 키 목록 테이블 (Copy/수정/삭제 버튼)
│       │   ├── KeyFormModal.vue    # 키 등록/수정 공용 모달
│       │   └── ExpiryModal.vue     # 만료 임박 알림 팝업
│       ├── stores/
│       │   ├── auth.ts             # Pinia: 로그인 상태
│       │   └── keys.ts             # Pinia: 키 목록 + 검색 필터
│       ├── router/
│       │   └── index.ts            # 라우트 정의 + 인증 가드
│       └── wailsjs/                # wails dev/build 시 자동 생성 — 수정 금지
│
├── build/                      # 빌드 리소스 (아이콘, 매니페스트)
└── docs/
    ├── Requirements.md         # 요구사항 명세 (확정본)
    ├── architecture_design.md  # 아키텍처 설계
    ├── database_design.md      # DB 스키마 및 쿼리
    └── progress.md             # 작업 진행 현황 ← 작업 후 반드시 업데이트
```

---

## 5. 백엔드 아키텍처

### 레이어 구조

```
app.go (Wails Handler)
    └── service/auth_service.go      ← 인증 로직
    └── service/apikey_service.go    ← 키 관리 로직
            └── repository/user_repo.go
            └── repository/apikey_repo.go
                    └── internal/database/db.go  (SQLite)
            └── internal/crypto/crypto.go        (암복호화)
```

### app.go — Wails 바인딩 메서드 (프론트엔드 인터페이스)

| 메서드 | 파라미터 | 반환 | 설명 |
|--------|----------|------|------|
| `Register` | email, password string | error | 회원가입 |
| `Login` | email, password string | LoginResult | 로그인 + AES 키 파생 |
| `Logout` | - | - | 세션 초기화 |
| `GetKeys` | KeyFilter | []APIKeyDTO | 키 목록 조회 (검색 포함) |
| `CreateKey` | CreateKeyRequest | error | 키 등록 |
| `UpdateKey` | UpdateKeyRequest | error | 키 수정 |
| `DeleteKey` | id int | error | 키 삭제 |
| `GetExpiringKeys` | days int | []APIKeyDTO | 만료 임박 키 조회 |

### App 구조체 세션 상태

```go
type App struct {
    ctx           context.Context
    currentUserID int64
    encryptionKey []byte  // 32바이트, 메모리에만 보관
    authService   *service.AuthService
    keyService    *service.APIKeyService
}
```

> **규칙:** `encryptionKey`는 절대 파일/DB에 저장하지 않습니다. 로그아웃 또는 앱 종료 시 `encryptionKey = nil`로 즉시 초기화합니다.

---

## 6. 데이터베이스 스키마 요약

상세 DDL과 인덱스 → `docs/database_design.md`

### users 테이블

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | INTEGER PK | 자동증가 |
| email | TEXT UNIQUE | 로그인 ID |
| password_hash | TEXT | bcrypt 해시 (인증용) |
| enc_salt | TEXT | PBKDF2 salt — 암호화 키 파생 전용 (hex) |
| created_at / updated_at | DATETIME | 자동 기록 |

### api_keys 테이블

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | INTEGER PK | 자동증가 |
| user_id | INTEGER FK | users.id CASCADE DELETE |
| key_name | TEXT NOT NULL | 키 식별 이름 |
| key_value | TEXT NOT NULL | **AES-256-GCM 암호화값** `base64(nonce‖ciphertext)` |
| url | TEXT | 관련 URL (선택) |
| expiry_date | DATE | 만료예정일 (선택, YYYY-MM-DD) |
| registered_date | DATE | 키 발급일 — 사용자 직접 입력 (선택) |
| memo | TEXT | 메모 (선택) |
| created_at / updated_at | DATETIME | 자동 기록 |

### 시작 시 실행할 PRAGMA

```sql
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
```

---

## 7. 암호화 규칙 (변경 금지)

```
로그인:  비밀번호 + enc_salt  →  PBKDF2-SHA256(100,000회)  →  32바이트 AES 키 (메모리)
저장:    평문  +  AES 키  →  AES-256-GCM  →  base64(12B nonce ‖ ciphertext)  →  DB
복호화:  DB값 base64 decode  →  nonce + ciphertext  →  AES-256-GCM decrypt  →  평문
```

- nonce는 매 암호화마다 `crypto/rand`로 새로 생성합니다.
- `enc_salt`(PBKDF2용)와 `password_hash`(bcrypt용)는 독립적으로 관리합니다.
- 화면에 표시하는 Key Value는 **마스킹 없이** 복호화된 실제 값을 보여줍니다.

---

## 8. 프론트엔드 규칙

### 라우트 구성

| 경로 | 컴포넌트 | 인증 가드 |
|------|----------|-----------|
| `/login` | LoginView.vue | 비인증 전용 |
| `/register` | RegisterView.vue | 비인증 전용 |
| `/` | MainView.vue | **인증 필수** |

### Wails 백엔드 호출 방법

```typescript
// wailsjs/go/main/App 에서 임포트 (자동 생성)
import { Login, GetKeys, CreateKey } from '../wailsjs/go/main/App'

const result = await Login(email, password)
```

### 자동 검색 (debounce)

- 검색 입력값이 변경되면 **300ms debounce** 후 `GetKeys(filter)` 호출
- `watch`로 `searchFilter`를 감시하고 debounce 타이머를 적용합니다.

### 만료 상태 색상 규칙

| 상태 | expiry_status 값 | 색상 |
|------|-----------------|------|
| 만료됨 (오늘 이전) | 3 | 빨간색 (red) |
| 임박 (30일 이내) | 2 | 주황색 (orange) |
| 정상 | 1 | 기본색 |
| 만료예정일 없음 | 0 | 기본색 |

---

## 9. 검색 조건 동작 규칙

- `key_name`(LIKE 부분 일치) **OR** `registered_date` 범위(From~To)
- 두 조건 모두 비어있으면 해당 사용자의 **전체 키 반환**
- `registered_date`는 `api_keys.registered_date` (사용자 입력 키 발급일) 기준

---

## 10. 현재 진행 상태

> 최신 진행 현황은 항상 `docs/progress.md`를 확인하십시오.

| Phase | 내용 | 상태 |
|-------|------|------|
| Phase 0 | 설계 및 기획 | ✅ 완료 |
| Phase 1 | 백엔드 구현 (Go) | ⬜ 미착수 |
| Phase 2 | 프론트엔드 구현 (Vue3/TS) | ⬜ 미착수 |
| Phase 3 | 통합 검증 및 Windows 빌드 | ⬜ 미착수 |

### 다음 작업 (Phase 1 시작점)

1. `go get github.com/mattn/go-sqlite3` — SQLite3 드라이버 추가
2. `go get golang.org/x/crypto` — bcrypt 추가 (이미 go.mod에 있을 수 있음)
3. `internal/database/db.go` — SQLite 연결 및 테이블 마이그레이션 구현
4. `internal/crypto/crypto.go` — PBKDF2 + AES-256-GCM 유틸리티 구현
5. `internal/models/`, `internal/repository/`, `internal/service/` 순서로 구현

---

## 11. 코딩 컨벤션

### Go

- 패키지는 `internal/` 하위에 기능별로 분리 (database / models / repository / service / crypto)
- 모든 공개 함수는 `error`를 반환하거나 결과 구조체에 `Error string` 포함
- Repository는 `*sql.DB`를 직접 받음 — 서비스 레이어에서 DB 인스턴스 주입
- 컨텍스트(`context.Context`)는 Wails startup에서 받은 것을 사용

### TypeScript / Vue

- `<script setup lang="ts">` 방식 (Composition API)
- Props/Emits는 명시적 타입 정의
- Wails 바인딩 타입은 `wailsjs/go/main/` 자동 생성 파일 기준으로 맞춤
- 컴포넌트 파일명은 PascalCase

### 공통

- 주석은 **WHY가 비자명할 때**만 작성 (WHAT 설명 주석 금지)
- 에러 메시지는 사용자에게 안전한 메시지만 노출 (내부 DB/암호화 오류 상세 노출 금지)

---

## 12. 참고 문서

| 문서 | 경로 | 내용 |
|------|------|------|
| 요구사항 명세 | `docs/Requirements.md` | 기능 요구사항 전체 |
| 아키텍처 설계 | `docs/architecture_design.md` | 레이어 구조, Wails 통신, 보안 설계 |
| DB 설계 | `docs/database_design.md` | 테이블 DDL, 인덱스, 주요 쿼리 |
| 진행 현황 | `docs/progress.md` | Phase별 작업 완료 여부 — **작업 후 반드시 업데이트** |
| Wails 공식 문서 | https://wails.io/docs/introduction | Wails v2 API 참조 |
