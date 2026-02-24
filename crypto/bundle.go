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

// EncryptedBundle defines the interface for terminal-specific encrypted key data
//
// Maps ID.terminals to their corresponding encrypted symmetric key data
// Used for secure distribution of symmetric keys to multiple user terminals
type EncryptedBundle interface {

	// Map returns the raw terminal-to-encrypted-data mapping
	// Key: ID.terminal string, Value: Encrypted key data ([]byte)
	Map() map[string][]byte

	// String returns a human-readable string representation of the bundle
	String() string

	// IsEmpty checks if the bundle contains any encrypted key data
	//
	// Returns: true if no terminal entries exist, false otherwise
	IsEmpty() bool

	// Contains checks if the bundle has encrypted data for a specific terminal
	//
	// Parameters:
	//   - terminal - ID.terminal to check
	// Returns: true if data exists for the terminal, false otherwise
	Contains(terminal string) bool

	// Get retrieves encrypted key data for a specific terminal
	//
	// Parameters:
	//  - terminal - ID.terminal to get data for
	// Returns: Encrypted key data ([]byte) or nil if terminal not found
	Get(terminal string) []byte

	// Set stores encrypted key data for a specific terminal
	//
	// Parameters:
	//   - terminal - ID.terminal to associate with
	//   - data     - Encrypted key data to store (must not be nil)
	Set(terminal string, data []byte)

	// Remove deletes encrypted key data for a specific terminal
	//
	// Parameters:
	//   - terminal - ID.terminal to remove data for
	Remove(terminal string)

	// Encode serializes the bundle into a StringKeyMap for network transmission
	//
	// Structures the data with user ID and terminal-specific encrypted keys
	//
	// Parameters:
	//   - did - User ID associated with the encrypted bundle
	// Returns: Encoded StringKeyMap (compatible with "message.keys" field)
	Encode(did ID) StringKeyMap
}

// DecodeEncryptedBundle parses encrypted key data from a "message.keys" StringKeyMap
//
// # Reconstructs an EncryptedBundle from encoded terminal-specific key data
//
// Parameters:
//   - encodedKeys - Encoded key data (from "message.keys" field)
//   - did         - Receiver's user ID (to validate key ownership)
//   - terminals   - List of valid terminals to extract data for
//
// Returns: Decoded EncryptedBundle (empty bundle if no valid data found)
func DecodeEncryptedBundle(encodedKeys StringKeyMap, did ID, terminals []string) EncryptedBundle {
	helper := GetEncryptedBundleHelper()
	return helper.DecodeBundle(encodedKeys, did, terminals)
}

// UserEncryptedBundle is the concrete implementation of the EncryptedBundle interface
//
// Uses a map to store terminal-to-encrypted-key-data associations
type UserEncryptedBundle struct {
	//EncryptedBundle

	// table stores the core terminal-to-encrypted-data mapping
	//
	// Key: terminal string, Value: Encrypted symmetric key data ([]byte)
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
