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

/**
 *  Entity (User/Group)
 *  <p>
 *      Base class of User and Group, ...
 *  </p>
 *
 *  <pre>
 *  properties:
 *      identifier - entity ID
 *      type       - entity type
 *      meta       - meta for generate ID
 *      documents  - visa for user, or bulletin for group
 *  </pre>
 */
type Entity interface {
	Object

	/**
	 *  Get entity ID
	 *
	 * @return ID
	 */
	ID() ID

	/**
	 *  Get ID.type
	 *
	 * @return network type
	 */
	Type() EntityType

	Meta() Meta
	Documents() []Document

	DataSource() EntityDataSource
	SetDataSource(facebook EntityDataSource)
}

/**
 *  Base Entity
 *  ~~~~~~~~~~~
 */
type BaseEntity struct {
	//Entity

	identifier ID

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
