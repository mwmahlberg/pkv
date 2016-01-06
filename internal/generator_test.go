package internal

import (
	"crypto/rand"
	"errors"
	"testing"
)

func TestValidKey(t *testing.T) {
	matrix := NewKey(rand.Reader)

	key, _ := matrix.GetKey(123456)

	if err := CheckCompleteKey(key, matrix.Matrix); err != nil {
		t.Errorf("Key validation failed")
	}
}

func TestInvalidKey(t *testing.T) {

	matrix := NewKey(rand.Reader)

	key, _ := matrix.GetKey(123456)

	r := []rune(key)

	r[6] = '0'
	r[7] = '0'

	k := string(r)

	if err := CheckCompleteKey(k, matrix.Matrix); err == nil {
		t.Errorf("Invalid key '&s' was accepted", k)
	}
}

func TestSeedSmall(t *testing.T) {

	matrix := NewKey(rand.Reader)

	if _, err := matrix.GetKey(minSeed - 1); err == nil {
		t.Errorf("seed lower than minimum seed was accepted: %d", minSeed-1)
	}

}

func TestSeedBig(t *testing.T) {

	matrix := NewKey(rand.Reader)

	if _, err := matrix.GetKey(maxSeed + 1); err == nil {
		t.Errorf("seed higher than maximum seed was accepted: %d", maxSeed+1)
	}
}

func TestSeedRange(t *testing.T) {

	if testing.Short() {
		return
	}

	matrix := NewKey(rand.Reader)
	for i := minSeed; i <= maxSeed; i++ {
		key, err := matrix.GetKey(uint64(i))
		if err != nil {
			t.Errorf("product key creation failed for seed %s: %v", i, err)
		}

		if err := CheckCompleteKey(key, matrix.Matrix); err != nil {
			t.Errorf("product key verification failed for seed %s:%v\nkey:%s", i, err, key)
		}
	}
}

func TestInvalidKeyInput(t *testing.T) {

	matrix := NewKey(rand.Reader)

	invalidFormattedKey := "1-2-3-4-5"

	if err := CheckCompleteKey(invalidFormattedKey, matrix.Matrix); err == nil {
		t.Errorf("Undecodable key string was accepted")
	}

}

type MockRand struct{}

func (mock MockRand) Read(b []byte) (n int, err error) {
	return 0, errors.New("expected error")
}
func TestNewKeyRandomGeneratorError(t *testing.T) {
	recovered := false

	func() {
		defer func() {
			recovered = true
			recover()
		}()
		_ = NewKey(MockRand{})
	}()

	if !recovered {
		t.Error("NewKey does not panic on error while reading from random")
	}
	//

}
