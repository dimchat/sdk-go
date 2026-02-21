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

/**
 *  Base EncryptedBundle
 */
type UserEncryptedBundle struct {
	//EncryptedBundle

	// terminal -> encrypted key.data
	table map[string][]byte
}

func NewUserEncryptedBundle(capacity int) *UserEncryptedBundle {
	return &UserEncryptedBundle{
		table: make(map[string][]byte, capacity),
	}
}

// Override
func (bundle *UserEncryptedBundle) Map() map[string][]byte {
	return bundle.table
}

// Override
func (bundle *UserEncryptedBundle) String() string {
	clazz := "EncryptedBundle"
	table := bundle.table
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<%s count=%d>\n", clazz, len(table)))
	for target, data := range table {
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
	return len(bundle.table) == 0
}

// Override
func (bundle *UserEncryptedBundle) Contains(key string) bool {
	_, exists := bundle.table[key]
	return exists
}

// Override
func (bundle *UserEncryptedBundle) Get(key string) []byte {
	data, ok := bundle.table[key]
	if !ok {
		return nil
	}
	return data
}

// Override
func (bundle *UserEncryptedBundle) Set(key string, data []byte) {
	if data == nil {
		delete(bundle.table, key)
	} else {
		bundle.table[key] = data
	}
}

// Override
func (bundle *UserEncryptedBundle) Remove(key string) {
	delete(bundle.table, key)
}

// Override
func (bundle *UserEncryptedBundle) Encode(did ID) StringKeyMap {
	helper := GetEncryptedBundleHelper()
	return helper.EncodeBundle(bundle, did)
}
