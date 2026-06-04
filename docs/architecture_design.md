# KeyNest - 아키텍처 설계

## 1. 개요

KeyNest는 Wails v2 기반의 standalone 데스크톱 애플리케이션으로, Go 백엔드와 Vue3/TypeScript 프론트엔드가 단일 실행 파일로 패키징됩니다. 인터넷 연결 없이 로컬에서 완전히 동작하며, 사용자 API 키를 AES-256-GCM 암호화로 안전하게 저장합니다.

## 2. 기술 스택

| 구분 | 기술 | 버전 |
|------|------|------|
| 데스크톱 프레임워크 | Wails | v2.12.0 |
| 백엔드 언어 | Go | 1.23 |
| 프론트엔드 | Vue 3 + TypeScript | Latest |
| 프론트엔드 빌드 | Vite | Latest |
| 상태 관리 | Pinia | Latest |
| 라우터 | Vue Router 4 | Latest |
| 데이터베이스 | SQLite3 | Latest |
| SQLite 드라이버 | mattn/go-sqlite3 | Latest |
| 암호화 | AES-256-GCM + PBKDF2-SHA256 | Go 표준 라이브러리 |
| 비밀번호 해시 | bcrypt | golang.org/x/crypto |

## 3. 전체 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│                   KeyNest Desktop App                    │
│                                                         │
│  ┌──────────────────────┐   ┌────────────────────────┐  │
│  │   Frontend (Vue3/TS) │   │    Backend (Go)         │  │
│  │                      │   │                         │  │
│  │  ┌────────────────┐  │   │  ┌───────────────────┐  │  │
│  │  │   Vue Router   │  │   │  │  Wails App Struct  │  │  │
│  │  │  /login        │  │   │  │  (Handler Layer)   │  │  │
│  │  │  /register     │  │──Wails──│                 │  │  │
│  │  │  /main         │  │  Bridge │  AuthService     │  │  │
│  │  └────────────────┘  │   │  │  KeyService       │  │  │
│  │                      │   │  └────────┬──────────┘  │  │
│  │  ┌────────────────┐  │   │           │             │  │
│  │  │  Pinia Stores  │  │   │  ┌────────▼──────────┐  │  │
│  │  │  - authStore   │  │   │  │  Repository Layer  │  │  │
│  │  │  - keyStore    │  │   │  │  UserRepository    │  │  │
│  │  └────────────────┘  │   │  │  KeyRepository     │  │  │
│  │                      │   │  └────────┬──────────┘  │  │
│  │  ┌────────────────┐  │   │           │             │  │
│  │  │   Components   │  │   │  ┌────────▼──────────┐  │  │
│  │  │  KeyTable      │  │   │  │   Crypto Layer     │  │  │
│  │  │  KeyFormModal  │  │   │  │  AES-256-GCM       │  │  │
│  │  │  ExpiryModal   │  │   │  │  PBKDF2-SHA256     │  │  │
│  │  └────────────────┘  │   │  └────────┬──────────┘  │  │
│  └──────────────────────┘   │           │             │  │
│                             │  ┌────────▼──────────┐  │  │
│                             │  │   SQLite3 DB       │  │  │
│                             │  │   (로컬 파일)       │  │  │
│                             │  └───────────────────┘  │  │
│                             └────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## 4. 디렉토리 구조

```
KeyNest/
├── main.go                     # 앱 진입점 (Wails 초기화)
├── app.go                      # Wails App struct & 바인딩 메서드
├── go.mod
├── go.sum
├── internal/
│   ├── database/
│   │   └── db.go               # SQLite 연결, 마이그레이션
│   ├── models/
│   │   ├── user.go             # User 구조체
│   │   └── apikey.go           # APIKey 구조체
│   ├── repository/
│   │   ├── user_repo.go        # Users CRUD
│   │   └── apikey_repo.go      # API Keys CRUD
│   ├── service/
│   │   ├── auth_service.go     # 인증 로직 (회원가입, 로그인)
│   │   └── apikey_service.go   # 키 관리 로직 (암복호화 포함)
│   └── crypto/
│       └── crypto.go           # 암호화 유틸리티
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── LoginView.vue   # 로그인 화면
│   │   │   ├── RegisterView.vue# 회원가입 화면
│   │   │   └── MainView.vue    # 메인(키 목록) 화면
│   │   ├── components/
│   │   │   ├── KeyTable.vue    # 키 목록 테이블 컴포넌트
│   │   │   ├── KeyFormModal.vue# 키 등록/수정 모달
│   │   │   └── ExpiryModal.vue # 만료 임박 알림 모달
│   │   ├── stores/
│   │   │   ├── auth.ts         # 인증 상태 (Pinia)
│   │   │   └── keys.ts         # 키 목록 상태 (Pinia)
│   │   ├── router/
│   │   │   └── index.ts        # 라우트 정의 (인증 가드 포함)
│   │   ├── wailsjs/            # Wails 자동 생성 바인딩
│   │   │   └── go/main/        # Go 메서드 TypeScript 바인딩
│   │   ├── App.vue
│   │   └── main.ts
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── build/                      # 빌드 리소스 (아이콘 등)
└── docs/
    ├── Requirements.md
    ├── architecture_design.md  # (본 파일)
    ├── database_design.md
    └── progress.md
```

## 5. Wails 통신 구조

Wails는 Go 메서드를 JavaScript/TypeScript에서 직접 호출 가능하게 바인딩합니다.

```
Frontend (TypeScript)          Backend (Go)
─────────────────────          ────────────────────
window.go.main.App             app.go (App struct)
  .Login(email, pw)     ──▶    func (a *App) Login(email, password string) LoginResult
  .Register(email, pw)  ──▶    func (a *App) Register(email, password string) error
  .GetKeys(filter)      ──▶    func (a *App) GetKeys(filter KeyFilter) []APIKeyDTO
  .CreateKey(req)       ──▶    func (a *App) CreateKey(req CreateKeyRequest) error
  .UpdateKey(req)       ──▶    func (a *App) UpdateKey(req UpdateKeyRequest) error
  .DeleteKey(id)        ──▶    func (a *App) DeleteKey(id int) error
  .GetExpiringKeys(days)──▶    func (a *App) GetExpiringKeys(days int) []APIKeyDTO
  .Logout()             ──▶    func (a *App) Logout()
```

## 6. 보안 설계

### 6.1 암호화 흐름

```
로그인 시:
  사용자 비밀번호 + enc_salt(DB 저장)
      │
      ▼ PBKDF2-SHA256 (100,000 iterations)
  32바이트 AES 암호화 키 (메모리에만 보관)

키 저장 시:
  평문 키값 + AES키
      │
      ▼ AES-256-GCM
  base64(12바이트 nonce + 암호문) → DB 저장

키 조회 시:
  DB 암호문 → AES-256-GCM 복호화 → 평문 → 프론트엔드 전달
```

### 6.2 세션 관리

- Go 백엔드 `App` 구조체에 `currentUserID`와 `encryptionKey`를 메모리에 보관
- 앱 종료 시 메모리에서 자동 소멸 (디스크에 저장하지 않음)
- 로그아웃 시 메모리의 키 정보 즉시 초기화

### 6.3 비밀번호 저장

- bcrypt (cost factor 12)로 해시 후 DB 저장
- 평문 비밀번호는 절대 DB에 저장하지 않음

## 7. 화면 흐름

```
앱 시작
   │
   ▼
로그인 화면 ──[회원가입 클릭]──▶ 회원가입 화면
   │ [로그인 성공]                       │ [등록 완료]
   ▼                                    ▼
만료 임박 키 체크                    로그인 화면
   │ [만료 임박 키 있음]
   ▼
만료 알림 모달 (팝업)
   │ [확인]
   ▼
메인 화면 (키 목록)
   ├─ 검색 (Key명 OR 등록일)
   ├─ 키 등록 모달
   ├─ 키 수정 모달
   ├─ 키 삭제 확인
   └─ 복사 버튼 (클립보드)
```

## 8. 에러 처리 원칙

- 모든 Go 메서드는 에러를 포함한 결과 구조체를 반환
- 프론트엔드는 에러 메시지를 토스트 알림으로 표시
- DB 오류, 복호화 실패는 로그에 기록하고 사용자에게 안전한 메시지만 전달
