/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
package mkm

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/crypto"
)

/**
 *  User account for communication
 *  <p>
 *      This class is for creating user account
 *  </p>
 *
 *  <pre>
 *  functions:
 *      (User)
 *      1. verify(data, signature) - verify (encrypted content) data and signature
 *      2. encrypt(data)           - encrypt (symmetric key) data
 *
 *      (LocalUser)
 *      3. sign(data)    - calculate signature of (encrypted content) data
 *      4. decrypt(data) - decrypt (symmetric key) data
 *  </pre>
 */
type User interface {
	Entity

	/**
	 *  Get all contacts of the user
	 *
	 * @return contact list
	 */
	Contacts() []ID

	/**
	 *  Get visa.terminal
	 *
	 * @return terminal list
	 */
	Terminals() []string

	/**
	 *  Verify data and signature with user's public keys
	 *
	 * @param data - message data
	 * @param signature - message signature
	 * @return true on correct
	 */
	Verify(data []byte, signature []byte) bool

	/**
	 *  Encrypt data, try visa.key first, if not found, use meta.key
	 *
	 * @param plaintext - message data
	 * @return encrypted data with targets (ID terminals)
	 */
	EncryptBundle(plaintext []byte) EncryptedBundle

	//
	//  Interfaces for Local User
	//

	/**
	 *  Sign data with user's private key
	 *
	 * @param data - message data
	 * @return signature
	 */
	Sign(data []byte) []byte

	/**
	 *  Decrypt data with user's private key(s)
	 *
	 * @param bundle - encrypted data with targets (ID terminals)
	 * @return plain text
	 */
	DecryptBundle(bundle EncryptedBundle) []byte

	//
	//  Interfaces for Visa
	//
	SignVisa(visa Visa) Visa
	VerifyVisa(visa Visa) bool
}

/**
 *  Base User
 *  ~~~~~~~~~
 */
type BaseUser struct {
	BaseEntity
}

func (user *BaseUser) Init(identifier ID) User {
	if user.BaseEntity.Init(identifier) != nil {
	}
	return user
}

// Override
func (user *BaseUser) Contacts() []ID {
	facebook := user.DataSource()
	if facebook == nil {
		//panic("user datasource not set yet")
		return nil
	}
	return facebook.GetContacts(user.ID())
}

// Override
func (user *BaseUser) Terminals() []string {
	docs := user.Documents()
	if docs == nil {
		//panic("failed to get documents")
		return nil
	}
	agent := GetVisaAgent()
	return agent.GetTerminals(docs)
}

// Override
func (user *BaseUser) Verify(data []byte, signature []byte) bool {
	meta := user.Meta()
	docs := user.Documents()
	if meta == nil || docs == nil {
		//panic("user not ready")
		return false
	}
	agent := GetVisaAgent()
	keys := agent.GetVerifyKeys(meta, docs)
	if keys == nil {
		//panic("failed to get verify keys")
		return false
	}
	for _, pubKey := range keys {
		if pubKey.Verify(data, signature) {
			// matched!
			return true
		}
	}
	// signature not match
	// TODO: check whether visa is expired, query new document for this contact
	return false
}

// Override
func (user *BaseUser) EncryptBundle(plaintext []byte) EncryptedBundle {
	meta := user.Meta()
	docs := user.Documents()
	if meta == nil || docs == nil {
		//panic("user not ready")
		return nil
	}
	// NOTICE: meta.key will never be changed, so use Visa.key to encrypt message
	//         is a better way
	agent := GetVisaAgent()
	return agent.EncryptBundle(plaintext, meta, docs)
}

//
//  Interfaces for Local User
//

// Override
func (user *BaseUser) Sign(data []byte) []byte {
	sKey := user.GetPrivateKeyForSignature()
	if sKey == nil {
		//panic("failed to get sign key for user")
		return nil
	}
	return sKey.Sign(data)
}

// Override
func (user *BaseUser) DecryptBundle(bundle EncryptedBundle) []byte {
	// NOTICE: if you provide a public key in visa for encryption,
	//         here you should return the private key paired with visa.key
	dict := bundle.Map()
	var plaintext []byte
	var keys []DecryptKey
	for terminal, ciphertext := range dict {
		// get private keys for terminal
		keys = user.GetPrivateKeysForDecryption(terminal)
		if keys == nil {
			//panic("failed to get decrypt keys for user")
			continue
		}
		// try decrypting it with each private key
		for _, priKey := range keys {
			plaintext = priKey.Decrypt(ciphertext, nil)
			if plaintext != nil && len(plaintext) > 0 {
				// OK!
				return plaintext
			}
		}
	}
	// decryption failed
	// TODO: check whether my visa key is changed, push new visa to this contact
	return nil
}

// Override
func (user *BaseUser) SignVisa(visa Visa) Visa {
	identifier := user.ID()
	helper := GetGeneralAccountHelper()
	docID := helper.GetDocumentID(visa.Map())
	if docID == nil {
		//panic("visa ID not found")
	} else if !docID.Address().Equal(identifier.Address()) {
		//panic("visa ID not matched")
		return nil
	}
	// NOTICE: only sign visa with the private key paired with your meta.key
	sKey := user.GetPrivateKeyForVisaSignature()
	if sKey == nil {
		//panic("failed to get sign key for visa: : " + user.ID().String())
		return nil
	}
	signature := visa.Sign(sKey)
	if signature == nil {
		//panic("failed to sign vis")
		return nil
	}
	return visa
}

// Override
func (user *BaseUser) VerifyVisa(visa Visa) bool {
	// NOTICE: only verify visa with meta.key
	//         (if meta not exists, user won't be created)
	identifier := user.ID()
	helper := GetGeneralAccountHelper()
	docID := helper.GetDocumentID(visa.Map())
	if docID == nil {
		//panic("visa ID not found")
	} else if !docID.Address().Equal(identifier.Address()) {
		//panic("visa ID not matched")
		return false
	}
	meta := user.Meta()
	if meta == nil {
		//panic("failed to get meta")
		return false
	}
	pKey := meta.PublicKey()
	if pKey == nil {
		//panic("failed to get verify key for visa: : " + user.ID().String())
		return false
	}
	return visa.Verify(pKey)
}

//
//  Private Keys
//

// protected
func (user *BaseUser) GetPrivateKeysForDecryption(terminal string) []DecryptKey {
	facebook := user.DataSource()
	if facebook == nil {
		//panic("user datasource not set yet")
		return nil
	}
	identifier := user.ID()
	if terminal == "" || terminal == "*" {
		return facebook.GetPrivateKeysForDecryption(identifier)
	}
	uid := CreateID(identifier.Name(), identifier.Address(), terminal)
	return facebook.GetPrivateKeysForDecryption(uid)
}

// protected
func (user *BaseUser) GetPrivateKeyForSignature() SignKey {
	facebook := user.DataSource()
	if facebook == nil {
		//panic("user datasource not set yet")
		return nil
	}
	identifier := user.ID()
	return facebook.GetPrivateKeyForSignature(identifier)
}

// protected
func (user *BaseUser) GetPrivateKeyForVisaSignature() SignKey {
	facebook := user.DataSource()
	if facebook == nil {
		//panic("user datasource not set yet")
		return nil
	}
	identifier := user.ID()
	return facebook.GetPrivateKeyForVisaSignature(identifier)
}
