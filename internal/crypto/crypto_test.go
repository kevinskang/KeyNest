package crypto_test

import (
	"strings"
	"testing"

	"KeyNest/internal/crypto"
)

func TestGenerateSalt_Unique(t *testing.T) {
	s1, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	s2, _ := crypto.GenerateSalt()
	if s1 == s2 {
		t.Error("두 번의 GenerateSalt 결과가 동일해서는 안 됩니다")
	}
}

func TestGenerateSalt_Length(t *testing.T) {
	salt, err := crypto.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	// 32바이트 → hex 64자
	if len(salt) != 64 {
		t.Errorf("salt 길이: got %d, want 64", len(salt))
	}
	// 올바른 hex 형식 (0-9, a-f)
	for _, c := range salt {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("salt에 유효하지 않은 문자 포함: %c", c)
			break
		}
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	k1, err := crypto.DeriveKey("password123", salt)
	if err != nil {
		t.Fatalf("DeriveKey: %v", err)
	}
	k2, _ := crypto.DeriveKey("password123", salt)

	if len(k1) != 32 {
		t.Errorf("키 길이: got %d, want 32", len(k1))
	}
	for i := range k1 {
		if k1[i] != k2[i] {
			t.Error("동일한 입력에서 다른 키가 파생됩니다")
			break
		}
	}
}

func TestDeriveKey_DifferentSalt(t *testing.T) {
	s1, _ := crypto.GenerateSalt()
	s2, _ := crypto.GenerateSalt()
	k1, _ := crypto.DeriveKey("password", s1)
	k2, _ := crypto.DeriveKey("password", s2)

	same := true
	for i := range k1 {
		if k1[i] != k2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("다른 salt로 같은 키가 파생됩니다")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	cases := []string{
		"ghp_xxxxxxxxxxxxxxxxxxxx",
		"sk-proj-1234567890abcdef",
		"",
		strings.Repeat("A", 1024), // 긴 키값
	}
	salt, _ := crypto.GenerateSalt()
	key, _ := crypto.DeriveKey("TestPassword!", salt)

	for _, plaintext := range cases {
		encrypted, err := crypto.Encrypt(key, plaintext)
		if err != nil {
			t.Errorf("Encrypt(%q): %v", plaintext, err)
			continue
		}
		decrypted, err := crypto.Decrypt(key, encrypted)
		if err != nil {
			t.Errorf("Decrypt(%q): %v", plaintext, err)
			continue
		}
		if decrypted != plaintext {
			t.Errorf("복호화 결과 불일치: got %q, want %q", decrypted, plaintext)
		}
	}
}

func TestEncrypt_UniqueCiphertext(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	key, _ := crypto.DeriveKey("password", salt)

	enc1, _ := crypto.Encrypt(key, "동일한 평문")
	enc2, _ := crypto.Encrypt(key, "동일한 평문")

	// nonce가 매번 다르므로 암호문도 달라야 한다
	if enc1 == enc2 {
		t.Error("동일한 평문도 매번 다른 암호문이어야 합니다 (랜덤 nonce)")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	key, _ := crypto.DeriveKey("correct-password", salt)
	wrongKey := make([]byte, 32)

	encrypted, _ := crypto.Encrypt(key, "secret value")
	_, err := crypto.Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Error("잘못된 키로 복호화가 성공해서는 안 됩니다")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	key, _ := crypto.DeriveKey("password", salt)

	_, err := crypto.Decrypt(key, "!!!invalid-base64!!!")
	if err == nil {
		t.Error("잘못된 base64 입력에서 에러가 발생해야 합니다")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	salt, _ := crypto.GenerateSalt()
	key, _ := crypto.DeriveKey("password", salt)

	// 12바이트(nonce) 미만인 base64 입력
	_, err := crypto.Decrypt(key, "c2hvcnQ=") // "short" = 5bytes
	if err == nil {
		t.Error("너무 짧은 암호문에서 에러가 발생해야 합니다")
	}
}
