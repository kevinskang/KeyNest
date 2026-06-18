# KeyNest - 만료일 알람 기능 코드 검토 보고서

- **검토 일자:** 2026-06-11
- **검토 대상:** 만료일 알람(ExpiryModal) 전체 흐름
- **검토 결과:** 버그 4건 발견 → 전건 수정 완료

---

## 1. 기능 개요

로그인 성공 직후, 만료일이 **30일 이내**이거나 **이미 만료된** API 키가 존재하면 경고 모달(`ExpiryModal`)을 자동으로 표시한다.

| ExpiryStatus | 조건 | 표시 |
|:---:|---|---|
| `0` | `expiry_date IS NULL` | 표시 안 함 |
| `1` | `expiry_date > now + 30d` | 표시 안 함 (정상) |
| `2` | `expiry_date <= now + 30d` (미만료) | 주황색 — 임박 |
| `3` | `expiry_date < now` (이미 만료) | 빨간색 — 만료됨 |

---

## 2. 데이터 흐름

```
LoginView.vue
  └─ handleLogin()
       ├─ await Login(email, password)          ← 세션 수립
       ├─ await keyStore.loadExpiringKeys()      ← 만료 키 조회 (days=30)
       └─ router.push('/')

         ↓ Wails IPC

app.go: GetExpiringKeys(days int)
  └─ keyService.GetExpiringKeys(uid, encKey, days)
       └─ keyRepo.FindExpiring(uid, days)
            └─ SQLite: WHERE expiry_date <= DATE('now', '+N days')
                        CASE → expiry_status 계산
            └─ 복호화 → []APIKeyDTO 반환

MainView.vue
  └─ onMounted()
       ├─ await keyStore.loadKeys()
       └─ expiringKeys.length > 0 → showExpiry = true
            └─ <ExpiryModal :keys="expiringKeys" :days="30" />
```

---

## 3. 검토 범위 파일 목록

| 계층 | 파일 |
|------|------|
| Wails 바인딩 | `app.go` |
| 서비스 | `internal/service/apikey_service.go` |
| 레포지토리 | `internal/repository/apikey_repo.go` |
| 테스트 | `internal/service/apikey_service_test.go` |
| Pinia 스토어 | `frontend/src/stores/keys.ts` |
| 모달 컴포넌트 | `frontend/src/components/ExpiryModal.vue` |
| 메인 뷰 | `frontend/src/views/MainView.vue` |
| 로그인 뷰 | `frontend/src/views/LoginView.vue` |
| 라우터 | `frontend/src/router/index.ts` |

---

## 4. 발견된 버그 및 수정 내역

### Bug 1 — [중요] `FindExpiring` CASE 문이 `days` 파라미터를 무시

**파일:** `internal/repository/apikey_repo.go`

**원인:**  
WHERE 절은 동적 `days` 파라미터를 올바르게 사용했으나, 같은 쿼리 내 CASE 문은 `+30 days`를 하드코딩하고 있었다.  
`days > 30`으로 호출 시(예: 60일) 31~60일 구간 키가 WHERE를 통과하면서 `expiry_status = 1`(정상)로 반환되어, 모달에 배지 없이 표시되는 UX 오류 발생.

```sql
-- 수정 전
WHEN DATE(expiry_date) <= DATE('now', '+30 days') THEN 2  -- 하드코딩

-- 수정 후
WHEN DATE(expiry_date) <= DATE('now', '+'||?||' days') THEN 2  -- 동적 바인딩
```

**파라미터 순서 변경:** `(days, userID, days)` — CASE용 `days`가 첫 번째 `?`로 추가됨.

---

### Bug 2 — [중요] CRUD 후 `expiringKeys` 스토어 미갱신

**파일:** `frontend/src/views/MainView.vue`

**원인:**  
`handleSave`, `handleDelete` 실행 후 `keyStore.loadKeys()`만 호출하고 `loadExpiringKeys()`를 호출하지 않아 스토어의 `expiringKeys`가 stale 상태로 유지됐다.  
키 삭제 후에도 만료 목록에 잔존하거나, 만료 임박 키 등록 후 목록에 미반영.

```typescript
// 수정 전
await keyStore.loadKeys()

// 수정 후
await Promise.all([keyStore.loadKeys(), keyStore.loadExpiringKeys()])
```

`handleSave`와 `handleDelete` 양쪽에 동일하게 적용. 두 조회를 병렬 실행하여 성능 유지.

---

### Bug 3 — [중요] `loadExpiringKeys` 에러 묵살

**파일:** `frontend/src/stores/keys.ts`

**원인:**  
`GetExpiringKeys` 호출 실패 시 `catch {}` 블록이 에러를 완전히 삼키고 `expiringKeys = []`로 설정했다. 백엔드 오류 발생 시 모달이 뜨지 않으면서 사용자는 만료 키 부재와 조회 실패를 구별할 수 없었다.

```typescript
// 수정 전
} catch {
  expiringKeys.value = []
}

// 수정 후
} catch (e) {
  expiringKeys.value = []
  error.value = e instanceof Error ? e.message : '만료 키 조회 중 오류가 발생했습니다.'
  return false
}
```

반환 타입을 `Promise<boolean>`으로 변경하여 호출부에서 성공/실패 판별 가능.

---

### Bug 4 — [낮음] `ExpiryModal` UI 문구 "30일" 하드코딩

**파일:** `frontend/src/components/ExpiryModal.vue`

**원인:**  
모달 안내 문구에 `30`이 리터럴로 고정되어 있어, `days` 값이 변경되어도 UI가 반영되지 않았다.

```html
<!-- 수정 전 -->
<p>아래 키들이 30일 이내에 만료되거나 이미 만료되었습니다.</p>

<!-- 수정 후 -->
<p>아래 키들이 {{ props.days }}일 이내에 만료되거나 이미 만료되었습니다.</p>
```

`days` prop 추가 (`withDefaults` 기본값 `30`), Vue 3.2 호환 방식 적용.

---

## 5. 정상 동작 확인 항목

| 항목 | 상태 | 비고 |
|------|:----:|------|
| SQLite CASE 문 status 0/1/2/3 계산 | ✅ | `FindByFilter` 정확 |
| 로그인 → loadExpiringKeys → navigate 순서 보장 | ✅ | `await` 체인 정상 |
| `validateISODate` 날짜 유효성 검증 | ✅ | Go `time.Parse`가 허위 날짜도 거부 |
| 라우터 guard (`requiresAuth`) | ✅ | 미인증 시 `/login` 리다이렉트 |
| 암호화 키 메모리 보관 (DB 미저장) | ✅ | `requireSession()` 패턴 정상 |
| 만료일 NULL 처리 (`NULLIF`, `sql.NullString`) | ✅ | INSERT/UPDATE/SELECT 일관 |
| `expiry_date` 인덱스 존재 | ✅ | `idx_api_keys_expiry(user_id, expiry_date)` |

---

## 6. 참고 사항

- **과거 만료 키 누적 문제:** `FindExpiring`은 만료 하한 날짜 제한이 없어, 수년 전 만료된 키도 모달에 계속 표시된다. 현재는 운영 중인 앱이므로 영향이 없지만, 장기 사용 시 모달이 누적 키로 혼잡해질 수 있다.
- **테스트 커버리지:** `GetExpiringKeys` 및 `FindExpiring` SQL 로직에 대한 통합 테스트가 없다. `ExpiryStatus` CASE 분기별 검증 테스트 추가를 권장한다.
