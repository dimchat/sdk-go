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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// Entity defines the core interface for all identifiable entities (User/Group)
//
// Serves as the base interface for user and group accounts, providing:
//   - Unique entity identification (ID and type)
//   - Access to identity metadata (Meta) and profile documents (Visa/Bulletin)
//   - Data source integration for extended entity information
//
// Core Properties:
//   - identifier: Unique entity ID (user/group address)
//   - type: Entity type classification (user/group/bot/etc.)
//   - meta: Identity metadata used to generate the entity ID
//   - documents: Profile documents (Visa for users, Bulletin for groups)
type Entity interface {
	Object

	// ID returns the unique identifier of the entity (user/group address)
	//
	// This is the primary key for all entity operations
	ID() ID

	// Type returns the entity type classification (EntityType)
	//
	// Also known as ID.type - identifies if the entity is a user, group, bot, etc.
	Type() EntityType

	// Meta returns the identity metadata used to generate the entity ID
	//
	// Contains cryptographic keys and seed data for identity verification
	Meta() Meta

	// Documents returns the profile documents associated with the entity
	//
	// For users:
	//     Visa documents (encryption keys, nickname, avatar, etc.)
	// For groups:
	//     Bulletin documents (title, founder, etc.)
	Documents() []Document

	// DataSource returns the entity data source for extended information
	//
	// Provides access to dynamic data (group members, user contacts, keys, ...)
	DataSource() EntityDataSource
	SetDataSource(facebook EntityDataSource)
}

// BaseEntity is the base implementation of the Entity interface
//
// Provides core entity functionality (ID storage and data source management)
// Used as the foundation for BaseUser and BaseGroup implementations
type BaseEntity struct {
	//Entity

	// identifier stores the unique entity ID (user/group address)
	identifier ID

	// facebook stores the data source for extended entity information
	//
	// Provides access to dynamic data not stored in the base entity
	facebook EntityDataSource
}

func NewBaseEntity(did ID) *BaseEntity {
	return &BaseEntity{
		identifier: did,
		facebook:   nil,
	}
}

// Override
func (entity *BaseEntity) Equal(other interface{}) bool {
	e, ok := other.(Entity)
	if ok {
		// compare pointers
		if entity == other {
			return true
		}
		// compare ids
		other = e.ID()
	}
	return entity.identifier.Equal(other)
}

// Override
func (entity *BaseEntity) DataSource() EntityDataSource {
	return entity.facebook
}

// Override
func (entity *BaseEntity) SetDataSource(facebook EntityDataSource) {
	entity.facebook = facebook
}

// Override
func (entity *BaseEntity) ID() ID {
	return entity.identifier
}

// Override
func (entity *BaseEntity) Type() EntityType {
	return entity.identifier.Type()
}

// Override
func (entity *BaseEntity) Meta() Meta {
	facebook := entity.DataSource()
	if facebook == nil {
		//panic("entity datasource not set yet")
		return nil
	}
	return facebook.GetMeta(entity.identifier)
}

// Override
func (entity *BaseEntity) Documents() []Document {
	facebook := entity.DataSource()
	if facebook == nil {
		//panic("entity datasource not set yet")
		return nil
	}
	return facebook.GetDocuments(entity.identifier)
}
