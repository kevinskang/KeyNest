# KeyNest Refactoring Plan

> 작성일: 2026-06-05

## 1. 목적

본 리팩토링 계획은 현재 `docs/Requirements.md`와 `docs/architecture_design.md`를 기준으로, KeyNest 프로젝트의 가독성, 보안, 성능을 동시에 향상시키기 위한 전략을 제시합니다.

- 가독성: 코드 유지보수성과 개발자 생산성 향상
- 보안: 암호화/세션 처리, 오류 노출 최소화, 입력 검증 강화
- 성능: DB 액세스, 암호화 처리, 프론트엔드 렌더링 최적화

## 2. 현황 요약

현재 프로젝트는 Wails 기반 Go 백엔드와 Vue 3/TypeScript 프론트엔드를 포함한 전형적인 데스크톱 앱 구조입니다. 주요 보안 요구사항과 아키텍처는 `docs/Requirements.md`, `docs/architecture_design.md`에 잘 정리되어 있으나, 리팩토링을 통해 다음 영역에서 개선 여지가 있습니다.

- Go 코드베이스의 레이어 분리와 책임 명확화
- 암호화/로그인/세션 관리의 안전성 강화
- 프론트엔드 상태관리와 컴포넌트 재사용성 개선
- DB 쿼리와 검색 로직 성능 최적화
- 에러 처리 및 사용자 피드백 일관성 확보

## 3. 주요 개선 영역

### 3.1 백엔드 가독성 및 구조화

- `internal/service`, `internal/repository`, `internal/crypto` 계층 간 책임을 명확히 구분
- Wails 핸들러(App 구조체)에서 비즈니스 로직 제거, 서비스 호출만 수행
- 서비스 함수 반환 타입 통일과 명확한 에러 래핑
- `internal/models`와 DTO 구분: 저장 모델과 프론트엔드 전달 모델 분리

### 3.2 보안 강화

- 로그인 시 `encryptionKey` 메모리 보관 정책 검증
- 로그아웃/앱 종료 시 `encryptionKey = nil`뿐 아니라 슬라이스 재할당으로 메모리 클리어
- 입력값 검증 강화: 이메일 유효성, 비밀번호 길이, 날짜 형식, URL 유효성
- 에러 메시지 필터링: 내부 DB/암호화 오류는 로그에 남기고 사용자에게는 일반화된 설명 전달
- DB 연결 시 SQL 인젝션 방지와 PRAGMA 설정 검증

### 3.3 성능 개선

- 검색 로직 최적화: `Keyname LIKE` + 날짜 범위 조건을 적절히 인덱스 활용 가능하게 개선
- 복호화 호출 최소화: 목록 조회 시 화면 표시용으로 필요한 값만 복호화
- 프론트엔드 자동 검색 디바운스와 상태관리 로직 정리
- Wails 바인딩 호출 최소화: 반복 호출 대신 필요 시 집계/배치 API 제공

### 3.4 프론트엔드 가독성 및 유지보수

- `src/stores` Pinia 상태 구조화: auth, key list, filter, loading/error 상태 분리
- 컴포넌트 재사용성 개선: 모달/테이블/버튼 패턴 통일
- 타입 정의 정리: `src/types.ts` 및 Wails 바인딩 타입 일관성 유지
- 뷰 로직 단순화: 화면별 핵심 기능만 유지, 비즈니스 로직은 서비스 메서드 호출

### 3.5 테스트 커버리지

- Go 단위 테스트 확대: 암호화 유틸, auth 서비스, 키 CRUD, 검색 필터
- 프론트엔드 유닛 테스트 및 컴포넌트 테스트 추가 (가능한 범위)
- `docs/progress.md`와 연계하여 리팩토링 상태 추적

## 4. 단계별 리팩토링 계획

### 4.1 1단계: 분석 및 설계

- 현재 코드베이스 구조와 주요 함수/모듈 목록화
- `docs/architecture_design.md`와 `docs/Requirements.md` 기준으로 개선 우선순위 선정
- 리팩토링 대상 목록 작성
  - `app.go` 핸들러 정리
  - 서비스 계층 단위 함수 검토
  - 암호화 유틸 인터페이스 점검
  - DB 쿼리/인덱스 설계 확인
  - 프론트엔드 상태 및 컴포넌트 구조 검토

### 4.2 2단계: 코드 정리 및 책임 분리

- `app.go`에 비즈니스 로직이 남아있다면 서비스로 이동
- 서비스와 저장소 호출 규약 통일
- DTO/Model 경계 명확화
- 반복 코드 제거 및 공통 함수 추출
- 함수/변수 명칭을 한글 요구사항과 일치하도록 개선

### 4.3 3단계: 보안 패치

- `internal/crypto/crypto.go` 테스트로 PBKDF2/AES-GCM 동작 검증
- `auth_service.go`와 `apikey_service.go`에 입력 검증 추가
- `Logout()` 및 앱 종료 시 키 클리어 로직 강화
- 에러 메시지 처리 일관성 확보

### 4.4 4단계: 성능 최적화

- 검색 쿼리와 인덱스 검토, 필요 시 `database_design.md`에 인덱스 추가 제안
- 불필요한 복호화 지연 또는 반복 방지
- 프론트엔드 검색/리스트 리렌더링 최소화
- Wails 호출 횟수 최적화, 가능하면 페이지 초기화 시 일괄 데이터 로딩

### 4.5 5단계: 검증 및 문서화

- 주요 리팩토링 결과를 단위 테스트와 수동 검증으로 확인
- `docs/progress.md`에 리팩토링 완료 상태 및 변경사항 반영
- `docs/architecture_design.md`와 `docs/database_design.md`에 주요 변경점 추가
- 코드 주석/문서 보강: WHY 중심 주석, 보안/성능 결정 근거 명시

## 5. 구체적 개선 항목 예시

### 5.1 Go 코드

- `internal/crypto/crypto.go`:
  - `DeriveKey`/`Encrypt`/`Decrypt` 인터페이스 명확화
  - 복호화 실패 시 상세 에러를 노출하지 않고 `ErrDecryptionFailed` 변환
- `internal/service/auth_service.go`:
  - `Register`와 `Login`에서 입력 검증 함수로 분리
  - `Login` 성공 시 `encryptionKey`를 안전하게 반환하지 않고 App 레벨에서 처리
- `internal/repository/apikey_repo.go`:
  - 검색 조건별 쿼리 빌더 정리
  - `GetKeys` 호출 시 필요한 컬럼만 선택

### 5.2 프론트엔드

- `frontend/src/stores/auth.ts`:
  - 인증 상태, 사용자 프로필, 오류 상태 분리
  - `login`, `logout`, `register` 액션명 일관성 유지
- `frontend/src/stores/keys.ts`:
  - 검색 필터 상태와 키 목록 상태 분리
  - `fetchKeys`, `refreshKeys`, `fetchExpiringKeys` 액션 명시적 분리
- `frontend/src/components/KeyTable.vue`:
  - 복사/수정/삭제 버튼 로직 단순화
  - 렌더링 조건 최소화

## 6. 기대 효과

- 유지보수 용이성 증가: 계층별 책임 명확화로 코드 변경이 쉬워집니다.
- 보안 신뢰도 향상: 암호화 키 처리와 입력 검증이 강화되어 데이터 노출 위험이 감소합니다.
- 성능 개선: 검색/조회 로직 최적화로 사용자 체감 속도가 빨라집니다.
- 문서 일치성: 리팩토링 결과가 `docs/`의 아키텍처·요구사항과 일치하게 됩니다.

## 7. 추적 및 완료 기준

- `docs/progress.md`에 리팩토링 단계별 진행 상태 갱신
- 주요 변경 사항을 커밋 단위로 분리하여 리뷰 가능하도록 정리
- 테스트 통과 및 기능 수동 점검 완료
- `docs/architecture_design.md`, `docs/database_design.md`에 리팩토링 반영 여부 확인
