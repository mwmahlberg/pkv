package verify

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
)

func GetCheckSum(i []byte) []byte {
	left := 0x0056
	right := 0x00AF
	s := []byte(i)

	for i := 0; i < len(s); i++ {
		right = right + int(s[i])
		if right > 0x00FF {
			right = right - 0x00FF
		}
		left = left + right
		if left > 0x00FF {
			left = left - 0x00FF
		}
	}
	sum := (left << 8) + right

	b := make([]byte, 3, 3)
	binary.PutVarint(b, int64(sum))

	return b
}

func checkKeyChecksumByte(k []byte) (err error) {

	g := GetCheckSum(k[:7])
	if !reflect.DeepEqual(k[7:], g) {
		err = errors.New("invalid checksum")
	}

	return
}

func KeyChecksum(key string) (err error) {
	
	if len(key) != 24 {
		return errors.New("Invalid key length")
	}
	
	k := strings.Replace(key, "-", "", -1)

	b, err := hex.DecodeString(k)

	return checkKeyChecksumByte(b)

}

func GetKeyByte(seed uint64, a, b, c uint8) byte {

	var r uint64

	a = a % 25
	b = b % 3

	if a%2 == 0 {
		r = ((seed >> a) & 0x000000FF) ^ ((seed >> b) | uint64(c))
	} else {
		r = ((seed >> a) & 0x000000FF) ^ ((seed >> b) & uint64(c))
	}
	r = r & 0xFF
	var bar []byte
	buf := bytes.NewBuffer(bar)
	binary.Write(buf, binary.LittleEndian, r)

	return buf.Bytes()[0]
}

func KeyPart(key string, part int, f, s, t uint8, bl []uint64) (err error) {

	if err := KeyChecksum(key); err != nil {
		return err
	}

	k := strings.Replace(key, "-", "", -1)

	b, err := hex.DecodeString(k)

	if err != nil {
		return err
	}

	d, _ := binary.Uvarint(b)
	
	for _,b := range bl {
		if d == b {
			return errors.New("Key blacklisted")
		}
	} 

	err = errors.New("invalid key")
	if b[part+3] != GetKeyByte(d, f, s, t) {
		return
	}

	return nil
}