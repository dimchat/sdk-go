/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package crypto

import (
	"fmt"
	"strings"

	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  User Encrypted Key Data with Terminals
 */
type EncryptedBundle interface {

	// terminal -> encrypted key.data
	Map() map[string][]byte

	String() string

	IsEmpty() bool

	Contains(key string) bool

	/**
	 *  Get encrypted key data for terminal
	 *
	 * @param terminal - ID terminal
	 * @return encrypted key data
	 */
	Get(key string) []byte

	/**
	 *  Put encrypted key data for terminal
	 *
	 * @param terminal - ID terminal
	 * @param data     - encrypted key data
	 */
	Set(terminal string, data []byte)

	/**
	 *  Remove encrypted key data for terminal
	 *
	 * @param terminal - ID terminal
	 * @return removed data
	 */
	Remove(terminal string)

	/**
	 *  Encode key data
	 *
	 * @param did - user ID
	 * @return encoded key data with targets (ID + terminals)
	 */
	Encode(did ID) StringKeyMap
}

/**
 *  Decode key data from 'message.keys'
 *
 * @param encodedKeys - encoded key data with targets (ID + terminals)
 * @param did         - receiver ID
 * @param terminals   - visa terminals
 * @return encrypted key data with targets (ID terminals)
 */
func DecodeEncryptedBundle(encodedKeys StringKeyMap, did ID, terminals []string) EncryptedBundle {
	helper := GetEncryptedBundleHelper()
	return helper.DecodeBundle(encodedKeys, did, terminals)
}

func NewEncryptedBundle() EncryptedBundle {
	return &UserEncryptedBundle{
		_map: make(map[string][]byte),
	}
}

/**
 *  Base EncryptedBundle
 */
type UserEncryptedBundle struct {
	//EncryptedBundle

	// terminal -> encrypted key.data
	_map map[string][]byte
}

func (bundle *UserEncryptedBundle) Init() EncryptedBundle {
	bundle._map = make(map[string][]byte)
	return bundle
}

// Override
func (bundle *UserEncryptedBundle) Map() map[string][]byte {
	return bundle._map
}

// Override
func (bundle *UserEncryptedBundle) String() string {
	clazz := "EncryptedBundle"
	info := bundle._map
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<%s count=%d>\n", clazz, len(info)))
	for target, data := range info {
		//if target == "" || data == nil {
		//	continue
		//}
		sb.WriteString(fmt.Sprintf("\t\"%s\": %d byte(s)\n", target, len(data)))
	}
	sb.WriteString(fmt.Sprintf("</%s>", clazz))
	return sb.String()
}

// Override
func (bundle *UserEncryptedBundle) IsEmpty() bool {
	return len(bundle._map) == 0
}

// Override
func (bundle *UserEncryptedBundle) Contains(key string) bool {
	_, exists := bundle._map[key]
	return exists
}

// Override
func (bundle *UserEncryptedBundle) Get(key string) []byte {
	return bundle._map[key]
}

// Override
func (bundle *UserEncryptedBundle) Set(key string, data []byte) {
	if data == nil {
		delete(bundle._map, key)
	} else {
		bundle._map[key] = data
	}
}

// Override
func (bundle *UserEncryptedBundle) Remove(key string) {
	delete(bundle._map, key)
}

// Override
func (bundle *UserEncryptedBundle) Encode(did ID) StringKeyMap {
	helper := GetEncryptedBundleHelper()
	return helper.EncodeBundle(bundle, did)
}
