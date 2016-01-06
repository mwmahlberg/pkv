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

/*
Package generate implements the generator part of pkv.
*/
package internal

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"io"

	//	pkv "gopkg.in/mwmahlberg/pkv.v1/verify"
)

const (
	minSeed = 1
	maxSeed = 2097151
)

// The KeyMatrix is the private part of your product keys.
type KeyMatrix struct {
	Matrix [][3]uint8 `json:"matrix"`
}

// NewKey generates a new KeyMatrix using sufficiently secure random numbers.
func NewKey(rnd io.Reader) KeyMatrix {

	pk := KeyMatrix{}
	var m [][3]uint8
	m = make([][3]uint8, 4)

	for i := 0; i < 4; i++ {

		rb := make([]byte, 3)

		_, err := rnd.Read(rb)

		m[i][0] = rb[0]
		m[i][1] = rb[1]
		m[i][2] = rb[2]

		if err != nil {
			panic(err)
		}

	}
	pk.Matrix = m
	return pk
}

func (pk *KeyMatrix) iv(k, v int) uint8 {
	return pk.Matrix[k][v]
}

// GetKey generates a new product key based on the seed provided
func (pk *KeyMatrix) GetKey(seed uint64) (string, error) {

	if seed < minSeed {
		return "", fmt.Errorf("invalid value for seed: %d (<%d)", seed, minSeed)
	} else if seed > maxSeed {
		return "", fmt.Errorf("invalid value for seed: %d (>%d)", seed, maxSeed)
	}

	k := make([]byte, 10)
	binary.PutUvarint(k, seed)

	k[3] = GetKeyByte(seed, pk.iv(0, 0), pk.iv(0, 1), pk.iv(0, 2))

	k[4] = GetKeyByte(seed, pk.iv(1, 0), pk.iv(1, 1), pk.iv(1, 2))
	k[5] = GetKeyByte(seed, pk.iv(2, 0), pk.iv(2, 1), pk.iv(2, 2))
	k[6] = GetKeyByte(seed, pk.iv(3, 0), pk.iv(3, 1), pk.iv(3, 2))
	cs := GetCheckSum(k[:7])
	k[7] = cs[0]
	k[8] = cs[1]
	k[9] = cs[2]

	key := fmt.Sprintf("%X", k)

	re := regexp.MustCompile(".{4}")
	parts := re.FindAllString(key, -1)
	return strings.Join(parts, "-"), nil
}

// CheckCompleteKey validates each key part of a product key
func CheckCompleteKey(key string, matrix [][3]uint8) (err error) {

	k := strings.Replace(key, "-", "", -1)

	b, err := hex.DecodeString(k)

	if err != nil {
		return err
	}

	d, _ := binary.Uvarint(b)

	err = errors.New("invalid key")

	for i := 0; i < 4; i++ {
		if b[i+3] != GetKeyByte(d, matrix[i][0], matrix[i][1], matrix[i][2]) {
			return fmt.Errorf("invalid key %d", i+1)
		}
	}

	return nil
}

