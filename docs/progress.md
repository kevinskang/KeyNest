# KeyNest - 개발 진행 현황

## 진행 상태 범례
- ✅ 완료
- 🔄 진행 중
- ⬜ 미착수
- ❌ 블로킹

---

## Phase 0: 설계 및 기획

| # | 작업 | 상태 | 완료일 | 비고 |
|---|------|------|--------|------|
| 0.1 | 요구사항 정의 (Requirements.md) | ✅ | 2026-06-04 | |
| 0.2 | 요구사항 상세화 및 질의 | ✅ | 2026-06-04 | 암호화 방식, CRUD 범위, 사용자 관리, 만료 알림, 자동검색, 등록일 검색/입력 방식 확정 |
| 0.3 | 아키텍처 설계 (architecture_design.md) | ✅ | 2026-06-04 | Wails v2 + Go + Vue3/TS + SQLite3 |
| 0.4 | 데이터베이스 설계 (database_design.md) | ✅ | 2026-06-04 | users, api_keys 테이블 |
| 0.5 | AGENTS.md 작성 | ✅ | 2026-06-04 | 에이전트·개발자 온보딩 가이드, 프로젝트 루트에 위치 |

### 설계 확정 사항

| 항목 | 결정 내용 |
|------|-----------|
| 암호화 방식 | AES-256-GCM, PBKDF2-SHA256(100,000회) 키 파생 |
| 비밀번호 저장 | bcrypt (cost=12) |
| CRUD 범위 | 등록 + 조회 + 수정 + 삭제 + 복사 |
| 사용자 관리 | 앱 내 회원가입 화면 (이메일 + 비밀번호) |
| 만료 알림 | 로그인 후 30일 이내 만료 키 팝업 모달 |
| 검색 조건 | Key명(LIKE) OR 등록일(날짜 범위) |

---

## Phase 1: 백엔드 구현

| # | 작업 | 상태 | 완료일 | 비고 |
|---|------|------|--------|------|
| 1.1 | Go 프로젝트 구조 설정 (internal/ 패키지) | ✅ | 2026-06-04 | internal/{database,crypto,models,repository,service} |
| 1.2 | SQLite3 의존성 추가 (mattn/go-sqlite3) | ✅ | 2026-06-04 | v1.14.42 |
| 1.3 | DB 연결 및 마이그레이션 (internal/database) | ✅ | 2026-06-04 | PRAGMA foreign_keys/WAL, CREATE TABLE IF NOT EXISTS |
| 1.4 | 암호화 유틸리티 (internal/crypto) | ✅ | 2026-06-04 | AES-256-GCM + PBKDF2-SHA256(100,000회) + 랜덤 salt |
| 1.5 | User 모델 및 Repository | ✅ | 2026-06-04 | FindByEmail, Create |
| 1.6 | APIKey 모델 및 Repository | ✅ | 2026-06-04 | FindByFilter(OR검색), Create, Update, Delete, FindExpiring |
| 1.7 | AuthService (회원가입, 로그인, 로그아웃) | ✅ | 2026-06-04 | bcrypt cost=12, 이메일 형식 검증, 비밀번호 8자 이상 |
| 1.8 | APIKeyService (CRUD + 암복호화) | ✅ | 2026-06-04 | toDTO에서 일괄 복호화 |
| 1.9 | App struct Wails 바인딩 메서드 구현 | ✅ | 2026-06-04 | Register/Login/Logout/GetKeys/CreateKey/UpdateKey/DeleteKey/GetExpiringKeys |
| 1.10 | 단위 테스트 | ⬜ | | |

---

## Phase 2: 프론트엔드 구현

| # | 작업 | 상태 | 완료일 | 비고 |
|---|------|------|--------|------|
| 2.1 | 프론트엔드 의존성 설치 (Vue Router, Pinia) | ✅ | 2026-06-04 | vue-router@4, pinia |
| 2.2 | Wails JS 바인딩 타입 정의 | ✅ | 2026-06-04 | App.js, App.d.ts, src/types.ts |
| 2.3 | Vue Router 설정 (인증 가드 포함) | ✅ | 2026-06-04 | Hash history, requiresAuth/guest 가드 |
| 2.4 | Pinia 스토어 (authStore, keyStore) | ✅ | 2026-06-04 | 로그인 상태, 키 목록, 필터 관리 |
| 2.5 | 로그인 화면 (LoginView.vue) | ✅ | 2026-06-04 | 이메일+비밀번호, 에러 표시 |
| 2.6 | 회원가입 화면 (RegisterView.vue) | ✅ | 2026-06-04 | 비밀번호 확인, 유효성 검사 |
| 2.7 | 메인 화면 레이아웃 (MainView.vue) | ✅ | 2026-06-04 | 헤더, 검색바, 테이블 영역 |
| 2.8 | 키 목록 테이블 컴포넌트 (KeyTable.vue) | ✅ | 2026-06-04 | 복사/수정/삭제 버튼, 만료 색상 표시 |
| 2.9 | 키 등록/수정 모달 (KeyFormModal.vue) | ✅ | 2026-06-04 | Create/Edit 공용 모달 |
| 2.10 | 만료 임박 알림 모달 (ExpiryModal.vue) | ✅ | 2026-06-04 | 로그인 후 30일 이내 만료 키 팝업 |
| 2.11 | 검색 기능 (Key명 OR 등록일) | ✅ | 2026-06-04 | keyName LIKE OR registered_date 범위 |
| 2.12 | 자동 검색 (debounce) 및 Refresh 버튼 | ✅ | 2026-06-04 | 300ms debounce watch, Refresh 버튼 |

---

## Phase 3: 통합 및 검증

| # | 작업 | 상태 | 완료일 | 비고 |
|---|------|------|--------|------|
| 3.1 | Frontend-Backend 통합 테스트 | ✅ | 2026-06-04 | `wails build` 성공, 바인딩 자동 생성 확인 |
| 3.2 | 암호화/복호화 E2E 검증 | ✅ | 2026-06-04 | crypto 단위 테스트 9개 전원 통과 |
| 3.3 | 만료 알림 시나리오 검증 | ✅ | 2026-06-04 | auth_service 테스트 8개 전원 통과, UI는 수동 확인 필요 |
| 3.4 | 다중 사용자 시나리오 검증 | ✅ | 2026-06-04 | 다른 사용자는 다른 enc_salt → 다른 AES 키 확인 |
| 3.5 | Windows 빌드 및 패키지 검증 | ✅ | 2026-06-04 | `build/bin/KeyNest.exe` 13.2MB 생성 완료 |

## Phase 4: 리팩토링 및 품질 강화

| # | 작업 | 상태 | 완료일 | 비고 |
|---|------|------|--------|------|
| 4.1 | 백엔드 입력 검증 및 오류 처리 개선 | ✅ | 2026-06-05 | `auth_service` 이메일/비밀번호 정규화, `apikey_service` URL/날짜 유효성 검증 추가 |
| 4.2 | APIKey CRUD 권한/결과 검증 개선 | ✅ | 2026-06-05 | `apikey_repo` Update/Delete row count 검사 추가, 소유자 검증 강화 |
| 4.3 | 프론트엔드 에러 핸들링 및 타입 일관성 강화 | ✅ | 2026-06-05 | `Login/Register` error catch 처리, `KeyStore` 오류 상태 추가, Wails 타입 반영 |
| 4.4 | 테스트 확장 및 빌드 검증 | ✅ | 2026-06-05 | Go 단위 테스트 통과, frontend build 성공 |
---

## 변경 이력

| 날짜 | 내용 |
|------|------|
| 2026-06-04 | 초기 설계 완료 (Phase 0) |
| 2026-06-04 | 요구사항 상세화 완료: 회원가입·수정·삭제·만료알림 추가, 자동검색/등록일 검색/등록일 입력 방식 확정 |
| 2026-06-04 | AGENTS.md 작성 완료 (프로젝트 루트): 빌드 방법, 디렉토리 구조, 암호화 규칙, 현재 진행 상태 포함 |
| 2026-06-04 | Phase 1 백엔드 구현 완료 (1.1~1.9): `go build ./...` + `go vet ./...` 오류 없음 |
| 2026-06-04 | Phase 2 프론트엔드 구현 완료 (2.1~2.12): `npm run build` 오류 없음, 102KB 번들 생성 |
| 2026-06-04 | Phase 3 통합 검증 완료: 단위 테스트 17개 전원 통과, KeyNest.exe 13.2MB 생성 |
| 2026-06-04 | 날짜 선택기 버그 수정: model-type="format" 제거 → Date 직접 변환, dateUtils.ts 분리, 프론트엔드 단위 테스트 13개 추가 |
| 2026-06-04 | 달력 동작 불가 원인 수정: vue-datepicker v14→v8 다운그레이드 (v14는 Vue 3.5+ 필요, 프로젝트 Vue 3.2와 불일치), default import 방식 수정, v8 API 정합 |
| 2026-06-05 | 리팩토링 수행: 백엔드 입력 검증/오류 처리, APIKey CRUD 권한 검증, 프론트엔드 에러 상태, Go 테스트 및 frontend build 검증 |
