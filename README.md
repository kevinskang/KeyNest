# KeyNest

> 개인 API 키와 토큰을 로컬에 안전하게 보관하는 Windows 데스크톱 앱

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2.12-red)](https://wails.io)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js)](https://vuejs.org)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?logo=sqlite)](https://sqlite.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 소개

KeyNest는 GitHub Personal Access Token, OpenAI API Key 등 다양한 키·토큰을 **오프라인 환경에서 안전하게 저장하고 관리**하는 standalone 데스크톱 애플리케이션입니다.

- 단일 실행 파일(`.exe`) — 별도 설치 불필요
- 인터넷 연결 없이 완전 로컬 동작
- AES-256-GCM 암호화로 키값을 안전하게 보관
- 암호화 키는 메모리에만 존재하며 디스크에 저장되지 않음

---

## 주요 기능

| 기능 | 설명 |
|------|------|
| 계정 관리 | 앱 내 회원가입 / 로그인 / 로그아웃 |
| 키 CRUD | 키 등록, 조회, 수정, 삭제 |
| 클립보드 복사 | Copy 버튼 한 번으로 키값을 클립보드에 복사 |
| 만료 알림 | 로그인 시 30일 이내 만료 예정 키를 팝업으로 알림 |
| 실시간 검색 | Key Name 부분 일치 검색 및 등록일 범위 검색 (debounce 300ms) |
| 만료 상태 표시 | 만료됨(빨간색) / 30일 이내 임박(주황색)으로 시각화 |

---

## 보안 구조

```
사용자 비밀번호
      │
      ▼ PBKDF2-SHA256 (100,000 iterations)
암호화 키 (메모리에만 보관)
      │
      ├─ 키 저장 시 ──▶ AES-256-GCM 암호화 ──▶ keynest.db
      └─ 키 조회 시 ◀── AES-256-GCM 복호화 ◀── keynest.db

로그인 비밀번호: bcrypt (cost=12) 해시 저장
```

- 암호화 키는 로그인 세션 동안만 메모리에 존재
- 로그아웃 또는 앱 종료 시 메모리에서 즉시 제거(zeroing)
- 데이터 파일(`keynest.db`)은 실행 파일과 같은 디렉터리에 생성

---

## 기술 스택

| 구분 | 기술 |
|------|------|
| 프레임워크 | [Wails v2.12](https://wails.io) |
| 백엔드 | Go 1.23 |
| 프론트엔드 | Vue 3 + TypeScript + Pinia + Vue Router 4 |
| 빌드 도구 | Vite |
| 데이터베이스 | SQLite3 (`mattn/go-sqlite3`) |
| 암호화 | AES-256-GCM, PBKDF2-SHA256, bcrypt |

---

## 사전 요구사항 (빌드 시)

| 도구 | 버전 |
|------|------|
| Go | 1.21 이상 |
| Node.js | 18 이상 |
| npm | 9 이상 |
| Wails CLI | v2.x |
| GCC (CGO) | TDM-GCC 또는 MinGW-w64 (go-sqlite3 컴파일에 필요) |

### Wails CLI 설치

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

---

## 빌드

```bash
# 저장소 클론
git clone https://github.com/<your-username>/KeyNest.git
cd KeyNest

# 의존성 설치
go mod tidy
cd frontend && npm install && cd ..

# 프로덕션 빌드 (build/bin/KeyNest.exe 생성)
wails build

# 개발 모드 실행 (핫리로드)
wails dev
```

빌드 결과물은 `build/bin/KeyNest.exe`에 생성됩니다.  
`KeyNest.exe`와 동일한 디렉터리에 `keynest.db`가 자동으로 생성됩니다.

---

## 사용 방법

### 1. 첫 실행 — 계정 생성

1. `KeyNest.exe` 실행
2. **회원가입** 버튼을 클릭하여 이메일과 비밀번호(8자 이상)로 계정 생성
3. 로그인 화면에서 이메일과 비밀번호를 입력하여 로그인

> 비밀번호를 잊어버리면 저장된 키를 복구할 수 없습니다. 비밀번호는 반드시 별도로 기억해 두세요.

### 2. 키 등록

1. 메인 화면에서 **`+ 키 등록`** 버튼 클릭
2. 모달에서 항목 입력 후 저장

| 항목 | 필수 여부 |
|------|-----------|
| Key Name | 필수 |
| Key Value | 필수 |
| URL | 선택 |
| 만료예정일 | 선택 |
| 등록일 | 선택 |
| 메모 | 선택 |

### 3. 키 조회 / 검색

- Key Name 검색창에 입력하면 0.3초 후 자동 검색
- 등록일 범위(From / To)로 날짜 기반 필터링 가능
- 두 조건은 OR 조건으로 동작

### 4. 키 수정 / 삭제

- 목록의 각 행에서 **수정** 또는 **삭제** 버튼 클릭
- 삭제된 데이터는 복구할 수 없음

### 5. 클립보드 복사

- 목록의 **Copy** 버튼을 클릭하면 키값이 클립보드에 복사됨
- 복사 완료 시 화면에 알림 표시

### 6. 만료 알림

- 로그인 성공 시 오늘 기준 **30일 이내 만료 예정 키**가 있으면 팝업으로 알림
- 이미 만료된 키도 포함하여 표시

---

## 프로젝트 구조

```
KeyNest/
├── main.go                     # Wails 진입점
├── app.go                      # Frontend ↔ Backend 바인딩
├── wails.json                  # Wails 프로젝트 설정
├── internal/
│   ├── crypto/                 # AES-256-GCM, PBKDF2 암호화
│   ├── database/               # SQLite 연결 및 마이그레이션
│   ├── models/                 # 도메인 모델 (User, APIKey, DTO)
│   ├── repository/             # DB 접근 계층
│   └── service/                # 비즈니스 로직 (Auth, APIKey)
├── frontend/
│   └── src/
│       ├── views/              # LoginView, RegisterView, MainView
│       ├── components/         # KeyTable, KeyFormModal, ExpiryModal 등
│       └── stores/             # Pinia 상태 관리
└── docs/                       # 요구사항, 설계 문서
```





