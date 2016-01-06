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
package generate

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	pkv "gopkg.in/mwmahlberg/pkv.v1/verify"
)

// The KeyMatrix is the private part of your product keys.
type KeyMatrix struct {
	Matrix [][3]uint8 `json:"matrix"`
}

// NewKey generates a new KeyMatrix using sufficiently secure random numbers.
func NewKey() KeyMatrix {

	pk := KeyMatrix{}
	var m [][3]uint8
	m = make([][3]uint8, 4)

	for i := 0; i < 4; i++ {

		rb := make([]byte, 3)

		_, err := rand.Read(rb)

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
func (pk *KeyMatrix) GetKey(seed uint64) string {
	k := make([]byte, 10)
	binary.PutUvarint(k, seed)

	k[3] = pkv.GetKeyByte(seed, pk.iv(0, 0), pk.iv(0, 1), pk.iv(0, 2))

	k[4] = pkv.GetKeyByte(seed, pk.iv(1, 0), pk.iv(1, 1), pk.iv(1, 2))
	k[5] = pkv.GetKeyByte(seed, pk.iv(2, 0), pk.iv(2, 1), pk.iv(2, 2))
	k[6] = pkv.GetKeyByte(seed, pk.iv(3, 0), pk.iv(3, 1), pk.iv(3, 2))
	cs := pkv.GetCheckSum(k[:7])
	k[7] = cs[0]
	k[8] = cs[1]
	k[9] = cs[2]

	key := fmt.Sprintf("%X", k)

	re := regexp.MustCompile(".{4}")
	parts := re.FindAllString(key, -1)
	return strings.Join(parts, "-")
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
		if b[i+3] != pkv.GetKeyByte(d, matrix[i][0], matrix[i][1], matrix[i][2]) {
			return fmt.Errorf("invalid key %d", i+1)
		}
	}

	return nil
}
