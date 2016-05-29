package internal

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestKeyPartValid(t *testing.T) {

	matrix := NewKey(rand.Reader)

	key, _ := matrix.GetKey(123456)

	for i := 0; i < 4; i++ {
		iv := matrix.Matrix[i]

		if err := KeyPart(key, i, iv[0], iv[1], iv[2], []uint64{}); err != nil {
			t.Errorf("Verification of key part %d failed for key '%s' with seed %d and iv %v", i+1, key, 123456, iv)
		}
	}

}

func TestKeyPartInvalidChecksum(t *testing.T) {

	matrix := NewKey(rand.Reader)
	iv := matrix.Matrix[0]
	key, _ := matrix.GetKey(123456)

	r := []rune(key)

	for i := 20; i <= 23; i++ {
		if r[i] != '0' {
			r[i] = '0'
		} else {
			r[i] = 'F'
		}
	}

	key = string(r)

	if err := KeyPart(key, 0, iv[0], iv[1], iv[2], []uint64{}); err == nil {
		t.Errorf("Invalid checksum was accepted")
	}

}

func TestKeyPartKeyBlacklisted(t *testing.T) {

	matrix := NewKey(rand.Reader)
	iv := matrix.Matrix[0]
	key, _ := matrix.GetKey(123456)

	if err := KeyPart(key, 0, iv[0], iv[1], iv[2], []uint64{123456}); err == nil {
		t.Error("Blacklisted Key was accepted")
	}
}

func TestInvalidKeyPart(t *testing.T) {

	matrix := NewKey(rand.Reader)
	iv := matrix.Matrix[0]

	key, _ := matrix.GetKey(123456)

	r := []rune(key)

	for i := 6; i <= 7; i++ {
		if r[i] != '0' {
			r[i] = '0'
		} else {
			r[i] = 'F'
		}
	}

	// Need to create a valid checksum, otherwise KeyPart will fail early
	k := strings.Replace(string(r), "-", "", -1)

	buf, _ := hex.DecodeString(k)

	cs := GetCheckSum(buf[:7])
	buf[7] = cs[0]
	buf[8] = cs[1]
	buf[9] = cs[2]

	tampered := fmt.Sprintf("%X", buf)

	re := regexp.MustCompile(".{4}")
	parts := re.FindAllString(tampered, -1)

	if err := KeyPart(strings.Join(parts, "-"), 0, iv[0], iv[1], iv[2], []uint64{}); err == nil {
		t.Error("Tampered key was accepted")
	}
}

func TestChecksumInvalidKeyLength(t *testing.T) {
	key := "1234-1234-1234-1234-123"

	if err := KeyChecksum(key); err == nil {
		t.Error("Invalid key length was accepted")
	}
}

func TestChecksumInvalidCharacter(t *testing.T) {
	key := "12#4-1234-1234-1234-123@"

	if err := KeyChecksum(key); err == nil {
		t.Error("Invalid characters were accepted")
	}
}
