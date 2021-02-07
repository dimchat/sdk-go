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
	. "github.com/dimchat/mkm-go/format"
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
type StorageCommand struct {
	BaseCommand

	_title string

	_data []byte
	_key []byte

	_plaintext []byte
	_password SymmetricKey
}

func (cmd *StorageCommand) Init(dict map[string]interface{}) *StorageCommand {
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

func (cmd *StorageCommand) InitWithTitle(title string) *StorageCommand {
	if cmd.BaseCommand.InitWithCommand(STORAGE) != nil {
		cmd._title = title
		cmd._data = nil
		cmd._key = nil
		cmd._plaintext = nil
		cmd._password = nil
	}
	return cmd
}

func (cmd *StorageCommand) ID() ID {
	return IDParse(cmd.Get("ID"))
}
func (cmd *StorageCommand) SetID(identifier ID) {
	if identifier == nil {
		cmd.Set("ID", nil)
	} else {
		cmd.Set("ID", identifier.String())
	}
}

func (cmd *StorageCommand) Title() string {
	if cmd._title == "" {
		title := cmd.Get("title")
		if title == nil {
			// (compatible with v1.0)
			//  contacts command: {
			//      command : 'contacts',
			//      data    : '...',
			//      key     : '...'
			//  }
			title = cmd.Get("command")
		}
		if title != nil {
			cmd._title = title.(string)
		}
	}
	return cmd._title
}

//
//  Encrypted data
//      encrypted by a random password before upload
//
func (cmd *StorageCommand) Data() []byte {
	if cmd._data == nil {
		data := cmd.Get("data")
		if data != nil {
			cmd._data = Base64Decode(data.(string))
		}
	}
	return cmd._data
}
func (cmd *StorageCommand) SetData(data []byte) {
	if data == nil {
		cmd.Set("data", nil)
	} else {
		cmd.Set("data", Base64Encode(data))
	}
	cmd._data = data
	cmd._plaintext = nil
}

//
//  Symmetric key
//      password to decrypt data
//      encrypted by user's public key before upload.
//      this should be empty when the storage data is "private_key".
//
func (cmd *StorageCommand) Key() []byte {
	if cmd._key == nil {
		key := cmd.Get("key")
		if key != nil {
			cmd._key = Base64Decode(key.(string))
		}
	}
	return cmd._key
}
func (cmd *StorageCommand) SetKey(key []byte) {
	if key == nil {
		cmd.Set("key", nil)
	} else {
		cmd.Set("key", Base64Encode(key))
	}
	cmd._key = key
	cmd._password = nil
}

//-------- Decryption

func (cmd *StorageCommand) DecryptWithSymmetricKey(password SymmetricKey) []byte {
	if cmd._plaintext == nil {
		if password == nil {
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

func (cmd *StorageCommand) DecryptWithPrivateKey(privateKey DecryptKey) []byte {
	if cmd._password == nil {
		cmd._password = cmd.DecryptKey(privateKey)
	}
	return cmd.DecryptWithSymmetricKey(cmd._password)
}

func (cmd *StorageCommand) DecryptKey(privateKey DecryptKey) SymmetricKey {
	data := cmd.Key()
	if data == nil {
		return nil
	}
	key := privateKey.Decrypt(data)
	if key == nil {
		//panic("failed to decrypt key")
		return nil
	}
	return SymmetricKeyParse(JSONDecode(key))
}
