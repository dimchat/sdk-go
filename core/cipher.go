/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
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
package sdk

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

/*  Situations:
                  +-------------+-------------+-------------+-------------+
                  |  receiver   |  receiver   |  receiver   |  receiver   |
                  |     is      |     is      |     is      |     is      |
                  |             |             |  broadcast  |  broadcast  |
                  |    user     |    group    |    user     |    group    |
    +-------------+-------------+-------------+-------------+-------------+
    |             |      A      |             |             |             |
    |             +-------------+-------------+-------------+-------------+
    |    group    |             |      B      |             |             |
    |     is      |-------------+-------------+-------------+-------------+
    |    null     |             |             |      C      |             |
    |             +-------------+-------------+-------------+-------------+
    |             |             |             |             |      D      |
    +-------------+-------------+-------------+-------------+-------------+
    |             |      E      |             |             |             |
    |             +-------------+-------------+-------------+-------------+
    |    group    |             |             |             |             |
    |     is      |-------------+-------------+-------------+-------------+
    |  broadcast  |             |             |      F      |             |
    |             +-------------+-------------+-------------+-------------+
    |             |             |             |             |      G      |
    +-------------+-------------+-------------+-------------+-------------+
    |             |      H      |             |             |             |
    |             +-------------+-------------+-------------+-------------+
    |    group    |             |      J      |             |             |
    |     is      |-------------+-------------+-------------+-------------+
    |    normal   |             |             |      K      |             |
    |             +-------------+-------------+-------------+-------------+
    |             |             |             |             |             |
    +-------------+-------------+-------------+-------------+-------------+
*/

// CipherKeyDestinationForMessage returns destination for cipher key vector: (sender, dest)
func CipherKeyDestinationForMessage(msg Message) ID {
	receiver := msg.Receiver()
	group := ParseID(msg.Get("group"))
	return CipherKeyDestination(receiver, group)
}

func CipherKeyDestination(receiver, group ID) ID {
	if group == nil && receiver.IsGroup() {
		/// Transform:
		///     (B) => (J)
		///     (D) => (G)
		group = receiver
	}
	if group == nil {
		/// A : personal message (or hidden group message)
		/// C : broadcast message for anyone
		return receiver
	}
	if group.IsBroadcast() {
		/// E : unencrypted message for someone
		//      return group as broadcast ID for disable encryption
		/// F : broadcast message for anyone
		/// G : (receiver == group) broadcast group message
		return group
	} else if receiver.IsBroadcast() {
		/// K : unencrypted group message, usually group command
		//      return receiver as broadcast ID for disable encryption
		return receiver
	}
	/// H    : group message split for someone
	/// J    : (receiver == group) non-split group message
	return group
}

// CipherKeyDelegate defines the interface for symmetric cipher key caching and retrieval
//
// Manages direction-specific symmetric keys used for message encryption between entities
// Keys are cached with sender→receiver directionality to enable key reuse for subsequent messages
type CipherKeyDelegate interface {

	// GetCipherKey retrieves a symmetric cipher key for message encryption between sender and receiver
	//
	// If the key doesn't exist and generate=true, creates a new symmetric key (AES/DES/...)
	//
	// Parameters:
	//   - sender   - From where (user or contact ID)
	//   - receiver - To where (contact or user/group ID)
	//   - generate - If true, generates a new key when no cached key exists; if false, returns nil
	// Returns: SymmetricKey for encryption (nil if no key exists)
	GetCipherKey(sender, receiver ID, generate bool) SymmetricKey

	// CacheCipherKey stores a symmetric cipher key for reuse in future message encryption
	//
	// Keys are cached with directionality (sender→receiver) - keys are not bidirectional
	//
	// Parameters:
	//   - sender   - From where (user or contact ID)
	//   - receiver - To where (contact or user/group ID)
	//   - key      - Symmetric key to cache (must not be nil)
	CacheCipherKey(sender, receiver ID, key SymmetricKey)
}
