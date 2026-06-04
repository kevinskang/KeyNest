# KeyNest - 데이터베이스 설계

## 1. 개요

- DBMS: SQLite3
- 파일 위치: 앱 실행 파일과 동일 디렉토리 `keynest.db`
- 문자 인코딩: UTF-8
- 외래키 제약: `PRAGMA foreign_keys = ON`

## 2. 테이블 목록

| 테이블명 | 설명 |
|----------|------|
| `users` | 사용자 계정 정보 |
| `api_keys` | 사용자별 API 키/토큰 정보 |

## 3. 테이블 상세

### 3.1 users

사용자 계정 및 암호화 키 파생에 필요한 salt를 저장합니다.

```sql
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT    NOT NULL UNIQUE,
    password_hash TEXT    NOT NULL,    -- bcrypt 해시 (cost=12)
    enc_salt      TEXT    NOT NULL,    -- PBKDF2 salt (hex, 32바이트)
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

#### 컬럼 설명

| 컬럼 | 타입 | NULL | 설명 |
|------|------|------|------|
| `id` | INTEGER | NOT NULL | 기본키, 자동증가 |
| `email` | TEXT | NOT NULL | 로그인 ID (이메일 형식), 유니크 |
| `password_hash` | TEXT | NOT NULL | bcrypt 해시 (인증용) |
| `enc_salt` | TEXT | NOT NULL | PBKDF2 salt (암호화 키 파생용, hex 인코딩) |
| `created_at` | DATETIME | NOT NULL | 계정 생성 일시 |
| `updated_at` | DATETIME | NOT NULL | 최종 수정 일시 |

> **보안 노트:** `enc_salt`는 암호화 키 파생(PBKDF2) 전용이며 `password_hash`(bcrypt)와 독립적으로 관리합니다. 로그인 시 비밀번호로 bcrypt 검증과 동시에 PBKDF2로 AES 키를 파생합니다.

---

### 3.2 api_keys

사용자가 등록한 API 키 및 토큰 정보를 저장합니다. `key_value`는 AES-256-GCM으로 암호화됩니다.

```sql
CREATE TABLE IF NOT EXISTS api_keys (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_name        TEXT     NOT NULL,
    key_value       TEXT     NOT NULL,    -- AES-256-GCM 암호화, base64(nonce||ciphertext)
    url             TEXT,
    expiry_date     DATE,                 -- 만료예정일 (YYYY-MM-DD)
    registered_date DATE,                 -- 키 최초 발급일 (YYYY-MM-DD)
    memo            TEXT,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id    ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_name   ON api_keys(user_id, key_name);
CREATE INDEX IF NOT EXISTS idx_api_keys_expiry     ON api_keys(user_id, expiry_date);
CREATE INDEX IF NOT EXISTS idx_api_keys_created_at ON api_keys(user_id, created_at);
```

#### 컬럼 설명

| 컬럼 | 타입 | NULL | 설명 |
|------|------|------|------|
| `id` | INTEGER | NOT NULL | 기본키, 자동증가 |
| `user_id` | INTEGER | NOT NULL | 소유 사용자 FK (users.id, CASCADE 삭제) |
| `key_name` | TEXT | NOT NULL | 키 식별 이름 (예: "GitHub Personal Token") |
| `key_value` | TEXT | NOT NULL | 암호화된 키값 `base64(12B nonce \|\| ciphertext)` |
| `url` | TEXT | NULL | 관련 서비스 URL (선택) |
| `expiry_date` | DATE | NULL | 만료예정일 (선택, YYYY-MM-DD) |
| `registered_date` | DATE | NULL | 키 최초 발급일 (선택, YYYY-MM-DD) |
| `memo` | TEXT | NULL | 메모 (선택) |
| `created_at` | DATETIME | NOT NULL | DB 등록 일시 |
| `updated_at` | DATETIME | NOT NULL | 최종 수정 일시 |

---

## 4. 암호화 상세

### 4.1 key_value 암호화 형식

```
평문 키값 (예: "ghp_xxxxxxxxxxxx")
         │
         ▼ AES-256-GCM
         │  - Key: PBKDF2(비밀번호, enc_salt, 100000, 32) → 32바이트
         │  - Nonce: crypto/rand 12바이트 (매 암호화마다 새로 생성)
         │
         ▼
base64(nonce[12B] || ciphertext[N+16B])  ← DB 저장값
```

### 4.2 복호화 흐름

```
DB에서 읽은 base64 문자열
    │
    ▼ base64 decode
nonce[12B] || ciphertext[...]
    │
    ▼ AES-256-GCM Decrypt (Key = PBKDF2(비밀번호, enc_salt))
평문 키값
```

---

## 5. 주요 쿼리

### 5.1 키 검색 (Key명 OR 등록일)

```sql
-- key_name 또는 created_at(날짜)로 검색 (OR 조건)
-- 조건 없으면 전체 조회
SELECT id, key_name, key_value, url, expiry_date, registered_date, memo, created_at, updated_at
FROM   api_keys
WHERE  user_id = ?
  AND  (
         (:key_name = '' OR key_name LIKE '%' || :key_name || '%')
         OR
         (:date_from = '' OR DATE(created_at) >= :date_from)
         OR
         (:date_to   = '' OR DATE(created_at) <= :date_to)
       )
ORDER BY created_at DESC;
```

> 조건이 모두 비어있으면 전체 행이 반환됩니다 (OR 조건 특성상 모두 true).

### 5.2 만료 임박 키 조회

```sql
-- 오늘부터 N일 이내 만료 예정인 키 + 이미 만료된 키
SELECT id, key_name, expiry_date
FROM   api_keys
WHERE  user_id = ?
  AND  expiry_date IS NOT NULL
  AND  DATE(expiry_date) <= DATE('now', '+' || ? || ' days')
ORDER BY expiry_date ASC;
```

### 5.3 만료예정일 기준 정렬 포함 전체 조회

```sql
SELECT id, key_name, key_value, url, expiry_date, registered_date, memo,
       created_at, updated_at,
       CASE
           WHEN expiry_date IS NULL                              THEN 0
           WHEN DATE(expiry_date) < DATE('now')                  THEN 3  -- 만료됨
           WHEN DATE(expiry_date) <= DATE('now', '+30 days')     THEN 2  -- 만료 임박
           ELSE                                                       1  -- 정상
       END AS expiry_status
FROM   api_keys
WHERE  user_id = ?
ORDER BY created_at DESC;
```

---

## 6. 마이그레이션 전략

- 초기 실행 시 테이블 자동 생성 (`CREATE TABLE IF NOT EXISTS`)
- 스키마 버전 관리를 위한 `schema_version` 테이블 (향후 확장용)

```sql
CREATE TABLE IF NOT EXISTS schema_version (
    version    INTEGER PRIMARY KEY,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

## 7. ER 다이어그램

```
┌──────────────────────┐         ┌───────────────────────────────┐
│        users         │         │           api_keys             │
├──────────────────────┤         ├───────────────────────────────┤
│ id            PK     │◄──┐     │ id              PK            │
│ email         UNIQUE │   └─FK─ │ user_id         FK NOT NULL   │
│ password_hash        │         │ key_name        NOT NULL       │
│ enc_salt             │         │ key_value       NOT NULL(암호화)│
│ created_at           │         │ url                            │
│ updated_at           │         │ expiry_date                    │
└──────────────────────┘         │ registered_date                │
                                 │ memo                           │
                                 │ created_at                     │
                                 │ updated_at                     │
                                 └───────────────────────────────┘
```
