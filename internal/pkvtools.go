// Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
)

// GetCheckSum generates hex encoded hash sum of the given bytes.
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

// KeyChecksum validates that the checksum of the product key is correct
// Note that this only validates wether the key was entered correctly,
// not if it is a valid generated product key.
func KeyChecksum(key string) (err error) {

	if len(key) != 24 {
		return errors.New("invalid key length")
	}

	k := strings.Replace(key, "-", "", -1)

	b, err := hex.DecodeString(k)

	if err != nil {
		return
	}

	return checkKeyChecksumByte(b)

}

// GetKeyByte generates a a byte representing a key part.
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

// KeyPart validates wether the chosen key part was generated with the KeyMatrix
func KeyPart(key string, part int, f, s, t uint8, bl []uint64) (err error) {

	if err = KeyChecksum(key); err != nil {
		return err
	}

	k := strings.Replace(key, "-", "", -1)

	// hex.decode was already tested at KeyChecksum
	b, _ := hex.DecodeString(k)

	d, _ := binary.Uvarint(b)

	for _, b := range bl {
		if d == b {
			return errors.New("key blacklisted")
		}
	}

	err = errors.New("invalid key")
	if b[part+3] != GetKeyByte(d, f, s, t) {
		return
	}

	return nil
}
