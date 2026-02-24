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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/mkm"
)

/**
 *  Entity Database
 */

// Archivist defines the interface for persistent storage of entity metadata and documents
//
// Acts as the database access layer for User/Group entities with core security constraints:
//   - All metadata/documents MUST be verified for authenticity before storage
//   - Provides access to local user IDs (required for message decryption)
type Archivist interface {

	// SaveMeta stores verified entity metadata to persistent storage
	//
	// Security Constraint: Metadata MUST be verified (signature/identity check) before saving
	//
	// Parameters:
	//   - meta - Verified entity metadata to save
	//   - did  - Target entity ID (user/group) associated with the metadata
	// Returns: true if metadata saved successfully, false otherwise
	SaveMeta(meta Meta, did ID) bool

	// SaveDocument stores a verified entity document to persistent storage
	//
	// Security Constraint: Document MUST be verified (signature/identity check) before saving
	//     For users: Saves Visa documents;
	//     For groups: Saves Bulletin documents
	//
	// Parameters:
	//   - document - Verified entity document to save
	//   - did      - Target entity ID (user/group) associated with the document
	// Returns: true if document saved successfully, false otherwise
	SaveDocument(document Document, did ID) bool

	// -------------------------------------------------------------------------
	//  Local User Management (Message Decryption)
	// -------------------------------------------------------------------------

	// LocalUsers retrieves IDs of all local user accounts with private key access
	//
	// Local users are required for decrypting incoming messages (private key ownership)
	//
	// Returns: Slice of local user IDs (empty slice if no local users exist)
	LocalUsers() []ID
}

/**
 *  Entity Factory
 */

// Barrack defines the entity factory interface (user/group pool manager)
//
// Maintains an in-memory cache of User/Group instances to avoid redundant instantiation
// Creates new entity instances only when required dependencies are available (lazy initialization)
type Barrack interface {

	// CacheUser stores a User instance in the in-memory cache
	//
	// Prevents redundant creation of the same user instance for subsequent access
	//
	// Parameters:
	//   - user - User instance to cache (must not be nil)
	CacheUser(user User)

	// CacheGroup stores a Group instance in the in-memory cache
	//
	// Prevents redundant creation of the same group instance for subsequent access
	//
	// Parameters:
	//   - group - Group instance to cache (must not be nil)
	CacheGroup(group Group)

	// GetUser retrieves a cached User instance by ID
	//
	// Returns nil if the user is not in the cache (use CreateUser for lazy initialization)
	//
	// Parameters:
	//   - uid - Target user ID to retrieve
	// Returns: Cached User instance (nil if not found)
	GetUser(uid ID) User

	// GetGroup retrieves a cached Group instance by ID
	//
	// Returns nil if the group is not in the cache (use CreateGroup for lazy initialization)
	//
	// Parameters:
	//   - gid - Target group ID to retrieve
	// Returns: Cached Group instance (nil if not found)
	GetGroup(gid ID) Group

	// CreateUser lazily initializes a User instance for the given ID (factory method)
	//
	// Creates a user only if the required visa.key (encryption key) is available
	//
	// Parameters:
	//   - uid - Target user ID to create
	// Returns: New User instance (nil if visa.key is missing/not ready)
	CreateUser(uid ID) User

	// CreateGroup lazily initializes a Group instance for the given ID (factory method)
	//
	// Creates a group only if the required member list is available (group is operational)
	//
	// Parameters:
	//   - gid - Target group ID to create
	// Returns: New Group instance (nil if members are missing/not ready)
	CreateGroup(gid ID) Group
}
