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
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
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
	IStorageCommand
}
type IStorageCommand interface {

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

//
//  Storage command implementation
//
type BaseStorageCommand struct {
	BaseCommand
	IStorageCommand

	_title string

	_data []byte
	_key []byte

	_plaintext []byte
	_password SymmetricKey
}

func (cmd *BaseStorageCommand) Init(dict map[string]interface{}) *BaseStorageCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._title = ""
		cmd._data = nil
		cmd._key = nil
		cmd._plaintext = nil
		cmd._password = nil
	}
	return cmd
}

func (cmd *BaseStorageCommand) InitWithTitle(title string) *BaseStorageCommand {
	if cmd.BaseCommand.InitWithCommand(STORAGE) != nil {
		cmd._title = title
		cmd._data = nil
		cmd._key = nil
		cmd._plaintext = nil
		cmd._password = nil
	}
	return cmd
}

//-------- IStorageCommand

func (cmd *BaseStorageCommand) ID() ID {
	return IDParse(cmd.Get("ID"))
}
func (cmd *BaseStorageCommand) SetID(identifier ID) {
	if ValueIsNil(identifier) {
		cmd.Remove("ID")
	} else {
		cmd.Set("ID", identifier.String())
	}
}

func (cmd *BaseStorageCommand) Title() string {
	if cmd._title == "" {
		text, ok := cmd.Get("title").(string)
		if !ok {
			// (compatible with v1.0)
			//  contacts command: {
			//      command : 'contacts',
			//      data    : '...',
			//      key     : '...'
			//  }
			text, _ = cmd.Get("command").(string)
		}
		cmd._title = text
	}
	return cmd._title
}

func (cmd *BaseStorageCommand) Data() []byte {
	if cmd._data == nil {
		base64, ok := cmd.Get("data").(string)
		if ok {
			cmd._data = Base64Decode(base64)
		}
	}
	return cmd._data
}
func (cmd *BaseStorageCommand) SetData(data []byte) {
	if ValueIsNil(data) {
		cmd.Remove("data")
	} else {
		cmd.Set("data", Base64Encode(data))
	}
	cmd._data = data
	cmd._plaintext = nil
}

func (cmd *BaseStorageCommand) Key() []byte {
	if cmd._key == nil {
		base64, ok := cmd.Get("key").(string)
		if ok {
			cmd._key = Base64Decode(base64)
		}
	}
	return cmd._key
}
func (cmd *BaseStorageCommand) SetKey(key []byte) {
	if ValueIsNil(key) {
		cmd.Remove("key")
	} else {
		cmd.Set("key", Base64Encode(key))
	}
	cmd._key = key
	// reset password
	cmd._password = nil
}

//-------- Decryption

func (cmd *BaseStorageCommand) DecryptKey(privateKey DecryptKey) SymmetricKey {
	data := cmd.Key()
	if data == nil {
		return nil
	}
	key := privateKey.Decrypt(data)
	if key == nil {
		//panic("failed to decrypt key")
		return nil
	}
	return SymmetricKeyParse(JSONDecodeMap(key))
}

func (cmd *BaseStorageCommand) DecryptWithSymmetricKey(password SymmetricKey) []byte {
	if cmd._plaintext == nil {
		if ValueIsNil(password) {
			panic("symmetric key empty")
			return nil
		}
		data := cmd.Data()
		if data == nil {
			return nil
		}
		cmd._plaintext = password.Decrypt(data)
	}
	return cmd._plaintext
}

func (cmd *BaseStorageCommand) DecryptWithPrivateKey(privateKey DecryptKey) []byte {
	if cmd._password == nil {
		cmd._password = cmd.DecryptKey(privateKey)
	}
	return cmd.DecryptWithSymmetricKey(cmd._password)
}
