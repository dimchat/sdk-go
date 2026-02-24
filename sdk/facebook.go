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
package sdk

import (
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/mkm"
)

// Facebook defines the core interface for all account/entity-related operations in the system
//
// Acts as the unified entry point for entity management, combining:
//   - EntityDelegate: User/Group instance creation/lookup
//   - EntityDataSource: Metadata/document retrieval for users/groups
//   - Persistent storage: Saving verified entity metadata and documents
//   - Local user selection: Resolving local user identities for message reception
type Facebook interface {
	EntityDelegate
	EntityDataSource
	//UserDataSource
	//GroupDataSource

	// SaveMeta stores verified entity metadata to persistent storage
	//
	// Wraps Archivist.SaveMeta with entity validation and access control
	//
	// Parameters:
	//   - meta - Verified entity metadata to save
	//   - did  - Target entity ID (user/group) associated with the metadata
	// Returns: true if metadata saved successfully, false otherwise
	SaveMeta(meta Meta, did ID) bool

	// SaveDocument stores a verified entity document to persistent storage
	//
	// Wraps Archivist.SaveDocument with entity validation and access control
	//     For users: Saves Visa documents;
	//     For groups: Saves Bulletin documents
	//
	// Parameters:
	//   - document - Verified entity document to save
	//   - did      - Target entity ID (user/group) associated with the document
	// Returns: true if document saved successfully, false otherwise
	SaveDocument(document Document, did ID) bool

	// SelectUser identifies the local user account for a given receiver ID
	//
	// Resolves which local user should process incoming messages for the target receiver
	// Used for personal messages or broadcast messages targeting local users
	//
	// Parameters:
	//   - receiver - Target receiver ID (user/broadcast address)
	// Returns: Local user ID (nil if no matching local user found)
	SelectUser(receiver ID) ID

	// SelectMember identifies the local user account within a group member list
	//
	// Finds which group member corresponds to a local user (for group message decryption)
	//
	// Parameters:
	//   - members - Slice of group member IDs
	// Returns: Local user ID from the member list (nil if no local user is a group member)
	SelectMember(members []ID) ID
}

// BaseFacebook is the base implementation of the Facebook interface
//
// Provides core dependencies and infrastructure for account/entity operations
type BaseFacebook struct {
	//Facebook

	// Archivist provides persistent storage access for entity metadata/documents
	//
	// Used by SaveMeta and SaveDocument for data persistence
	Archivist Archivist

	// Barrack manages in-memory cache of User/Group instances (entity factory)
	//
	// Used for lazy initialization and reuse of entity instances
	Barrack Barrack

	// DataSource provides access to entity metadata and documents
	//
	// Backend implementation for EntityDataSource interface methods
	DataSource EntityDataSource
}

func NewBaseFacebook(db EntityDataSource) *BaseFacebook {
	return &BaseFacebook{
		DataSource: db,
		Archivist:  nil,
		Barrack:    nil,
	}
}

func (facebook *BaseFacebook) SaveMeta(meta Meta, did ID) bool {
	archivist := facebook.Archivist
	return archivist.SaveMeta(meta, did)
}

func (facebook *BaseFacebook) SaveDocument(document Document, did ID) bool {
	archivist := facebook.Archivist
	return archivist.SaveDocument(document, did)
}

func (facebook *BaseFacebook) SelectUser(receiver ID) ID {
	archivist := facebook.Archivist
	users := archivist.LocalUsers()
	if len(users) == 0 {
		//panic("local users should not be empty")
		return nil
	} else if receiver.IsBroadcast() {
		// broadcast message can be decrypted by anyone, so
		// just return current user here
		return users[0]
	}
	// personal message
	for _, item := range users {
		if receiver.Equal(item) {
			// DISCUSS: set this item to be current user?
			return item
		}
	}
	// not for me?
	return nil
}

func (facebook *BaseFacebook) SelectMember(members []ID) ID {
	archivist := facebook.Archivist
	users := archivist.LocalUsers()
	if len(users) == 0 {
		//panic("local users should not be empty")
		return nil
	}
	// group message (recipient not designated)
	for _, item := range users {
		for _, m := range members {
			if m.Equal(item) {
				// DISCUSS: set this item to be current user?
				return item
			}
		}
	}
	// not for me?
	return nil
}

//-------- EntityDelegate

// Override
func (facebook *BaseFacebook) GetUser(uid ID) User {
	barrack := facebook.Barrack
	if barrack == nil {
		//panic("barrack not ready")
		return nil
	}
	// get from user cache
	user := barrack.GetUser(uid)
	if user == nil {
		// create user and cache it
		user = barrack.CreateUser(uid)
		if user != nil {
			barrack.CacheUser(user)
		}
	}
	return user
}

// Override
func (facebook *BaseFacebook) GetGroup(gid ID) Group {
	barrack := facebook.Barrack
	if barrack == nil {
		//panic("barrack not ready")
		return nil
	}
	// get from group cache
	group := barrack.GetGroup(gid)
	if group == nil {
		// create group and cache it
		group = barrack.CreateGroup(gid)
		if group != nil {
			barrack.CacheGroup(group)
		}
	}
	return group
}

//-------- EntityDataSource

// Override
func (facebook *BaseFacebook) GetMeta(did ID) Meta {
	db := facebook.DataSource
	return db.GetMeta(did)
}

// Override
func (facebook *BaseFacebook) GetDocuments(did ID) []Document {
	db := facebook.DataSource
	return db.GetDocuments(did)
}

//-------- GroupDataSource

// Override
func (facebook *BaseFacebook) GetFounder(gid ID) ID {
	db := facebook.DataSource
	return db.GetFounder(gid)
}

// Override
func (facebook *BaseFacebook) GetOwner(gid ID) ID {
	db := facebook.DataSource
	return db.GetOwner(gid)
}

// Override
func (facebook *BaseFacebook) GetMembers(gid ID) []ID {
	db := facebook.DataSource
	return db.GetMembers(gid)
}

//-------- UserDataSource

// Override
func (facebook *BaseFacebook) GetContacts(uid ID) []ID {
	db := facebook.DataSource
	return db.GetContacts(uid)
}

// Override
func (facebook *BaseFacebook) GetPrivateKeysForDecryption(uid ID) []DecryptKey {
	db := facebook.DataSource
	return db.GetPrivateKeysForDecryption(uid)
}

// Override
func (facebook *BaseFacebook) GetPrivateKeyForSignature(uid ID) SignKey {
	db := facebook.DataSource
	return db.GetPrivateKeyForSignature(uid)
}

// Override
func (facebook *BaseFacebook) GetPrivateKeyForVisaSignature(uid ID) SignKey {
	db := facebook.DataSource
	return db.GetPrivateKeyForVisaSignature(uid)
}
