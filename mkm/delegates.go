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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

// EntityDelegate defines the interface for creating User/Group entity instances by ID
//
// Acts as a factory interface for entity instantiation, abstracting entity creation logic
type EntityDelegate interface {

	// GetUser creates/retrieves a User instance for the given user ID
	//
	// Parameters:
	//   - uid - Target user ID to create/retrieve
	// Returns: User instance (nil if the user ID is invalid/unknown)
	GetUser(uid ID) User

	// GetGroup creates/retrieves a Group instance for the given group ID
	//
	// Parameters:
	//   - gid - Target group ID to create/retrieve
	// Returns: Group instance (nil if the group ID is invalid/unknown)
	GetGroup(gid ID) Group
}

// EntityDataSource defines the core interface for accessing entity metadata and documents
//
// Combines GroupDataSource and UserDataSource to provide unified access to:
//  1. Metadata (user/group identity metadata generated with private keys)
//  2. Documents (visa/bulletin documents for entity profiles)
//  3. Cryptographic keys (verification/encryption keys from meta/documents)
//
// Key Metadata Rules:
//  1. User meta: Generated with the user's private key
//  2. Group meta: Generated with the founder's private key
//  3. Meta key: Used to verify messages sent by the user/group founder
//  4. Visa key: Used to encrypt messages for the target user
type EntityDataSource interface {
	GroupDataSource
	UserDataSource

	// GetMeta retrieves the metadata (Meta) for a given entity ID (user/group)
	//
	// Parameters:
	//   - did - Target entity ID (user or group)
	// Returns: Meta instance (nil if no metadata exists for the entity)
	GetMeta(did ID) Meta

	// GetDocuments retrieves all profile documents for a given entity ID (user/group)
	//
	// For users: Returns Visa documents; For groups: Returns Bulletin documents
	//
	// Parameters:
	//   - did - Target entity ID (user or group)
	// Returns: Slice of Document instances (empty slice if no documents exist)
	GetDocuments(did ID) []Document
}

// GroupDataSource defines the interface for accessing group-specific metadata
//
// Provides access to group ownership and membership information with core rules:
//  1. Group founder's public key matches the group meta.key (identity verification)
//  2. Group owner and members must be managed according to consensus algorithms
type GroupDataSource interface {

	// GetFounder retrieves the unique ID of the group's founder/creator
	//
	// The founder's private key is used to generate the group's initial metadata
	//
	// Parameters:
	//   - group - Target group ID
	// Returns: Founder's user ID
	GetFounder(group ID) ID

	// GetOwner retrieves the current owner of the group
	//
	// The owner may differ from the founder (ownership transfer via consensus)
	//
	// Parameters:
	//   - group - Target group ID
	// Returns: Current owner's user ID
	GetOwner(group ID) ID

	// GetMembers retrieves all member IDs in the target group
	//
	// Member list must comply with group consensus algorithm rules
	//
	// Parameters:
	//   - group - Target group ID
	// Returns: Slice of member user IDs (empty slice if group has no members)
	GetMembers(group ID) []ID
}

// UserDataSource defines the interface for accessing user-specific data and cryptographic keys
//
// Manages user contacts and private keys for encryption/decryption/signature operations
//
// Core Cryptographic Key Rules:
//
//	Encryption/Decryption:
//	    1. Public encryption key: Uses visa.key if available, falls back to meta.key
//	    2. Private decryption keys: Paired with [visa.key, meta.key] (multiple keys supported)
//
//	Signature/Verification:
//	    3. Private signature key: Paired with either visa.key or meta.key
//	    4. Public verification keys: Combined set of [visa.key, meta.key]
//
//	Visa Document Signing:
//	    5. Private visa signature key: Exclusively paired with meta.key
//	    6. Public visa verification key: Only meta.key (visa authenticity validation)
type UserDataSource interface {

	// GetContacts retrieves the list of contact IDs for the target user
	//
	// Contacts represent other users/groups the user has added to their address book
	//
	// Parameters:
	//   - user - Target user ID
	// Returns: Slice of contact entity IDs (empty slice if no contacts)
	GetContacts(user ID) []ID

	// GetPrivateKeysForDecryption retrieves private keys for message decryption
	//
	// These keys are paired with the user's public keys (visa.key and meta.key)
	//
	// Parameters:
	//   - user - Target user ID
	// Returns: Slice of DecryptKey instances (empty slice if no decryption keys)
	GetPrivateKeysForDecryption(user ID) []DecryptKey

	// GetPrivateKeyForSignature retrieves the private key for message signing
	//
	// This key is paired with either the user's visa.key or meta.key (asymmetric signature)
	//
	// Parameters:
	//   - user - Target user ID
	// Returns: SignKey instance (nil if no signature key exists)
	GetPrivateKeyForSignature(user ID) SignKey

	// GetPrivateKeyForVisaSignature retrieves the private key for signing visa documents
	//
	// This key is exclusively paired with meta.key (visa authenticity validation)
	//
	// Parameters:
	//   - user - Target user ID
	// Returns: SignKey instance (nil if no visa signature key exists)
	GetPrivateKeyForVisaSignature(user ID) SignKey
}
