/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
package protocol

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

const (
	STORAGE = "storage"

	// storage titles (should be encrypted)
	CONTACTS = "contacts"
	PRIVATE_KEY = "private_key"
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "storage",
 *      title   : "key name",  // "contacts", "private_key", ...
 *
 *      data    : "...",       // base64_encode(symmetric)
 *      key     : "...",       // base64_encode(asymmetric)
 *
 *      // -- extra info
 *      //...
 *  }
 */
type StorageCommand interface {
	Command

	ID() ID
	SetID(identifier ID)

	Title() string

	/**
	 *  Encrypted data
	 *  ~~~~~~~~~~~~~~
	 *  encrypted by a random password before upload
	 */
	Data() []byte
	SetData(data []byte)

	/**
	 *  Symmetric key
	 *  ~~~~~~~~~~~~~
	 *  password to decrypt data
	 *  encrypted by user's public key before upload.
	 *  this should be empty when the storage data is "private_key".
	 */
	Key() []byte
	SetKey(key []byte)

	// Decrypt key data to password with private key.
	DecryptKey(privateKey DecryptKey) SymmetricKey

	// Decrypt data with password.
	DecryptWithSymmetricKey(password SymmetricKey) []byte

	// Decrypt key data to password with private key,
	// after that, decrypt data with the password.
	DecryptWithPrivateKey(privateKey DecryptKey) []byte
}
