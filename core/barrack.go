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
 *  Database Access
 */
type Archivist interface {

	/**
	 *  Save meta for entity ID (must verify first)
	 *
	 * @param did  - entity ID
	 * @param meta - entity meta
	 * @return true on success
	 */
	SaveMeta(meta Meta, did ID) bool

	/**
	 *  Save entity document with ID (must verify first)
	 *
	 * @param did      - entity ID
	 * @param document - entity document
	 * @return true on success
	 */
	SaveDocument(document Document, did ID) bool

	//
	//  Local Users
	//

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return users with private key
	 */
	LocalUsers() []ID
}

/**
 *  Entity Factory
 *  <p>
 *      Entity pool to manage User/Group instances
 *  </p>
 */
type Barrack interface {

	// memory cache
	CacheUser(user User)
	CacheGroup(group Group)

	GetUser(uid ID) User
	GetGroup(gid ID) Group

	/**
	 *  Create user when visa.key exists
	 *
	 * @param uid - user ID
	 * @return user, null on not ready
	 */
	CreateUser(uid ID) User

	/**
	 *  Create group when members exist
	 *
	 * @param gid - group ID
	 * @return group, null on not ready
	 */
	CreateGroup(gid ID) Group
}
