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
package dimp

import (
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Entity Manager
 *  ~~~~~~~~~~~~~~
 */
type EntityManager interface {

	/**
	 *  Get current user (for signing and sending message)
	 *
	 * @return User
	 */
	GetCurrentUser() User

	/**
	 *  Document checking
	 *
	 * @param doc - entity document
	 * @return true on accepted
	 */
	CheckDocument(doc Document) bool

	/**
	 *  Save entity document with ID (must verify first)
	 *
	 * @param doc - entity document
	 * @return true on success
	 */
	SaveDocument(doc Document) bool

	/**
	 *  Save meta for entity ID (must verify first)
	 *
	 * @param meta - entity meta
	 * @param identifier - entity ID
	 * @return true on success
	 */
	SaveMeta(meta Meta, identifier ID) bool

	/**
	 *  Save members of group
	 *
	 * @param members - member ID list
	 * @param group - group ID
	 * @return true on success
	 */
	SaveMembers(members []ID, group ID) bool

	IsFounder(member ID, group ID) bool
	IsOwner(member ID, group ID) bool
}

/**
 *  Delegate for Entity
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Facebook struct {
	Barrack
	EntityManager

	_manager EntityManager
}

func (facebook *Facebook) Init() *Facebook {
	if facebook.Barrack.Init() != nil {
		facebook._manager = nil
	}
	return facebook
}

func (facebook *Facebook) SetDelegate(delegate EntityDelegate) {
	facebook.Barrack.SetDelegate(delegate)
	handler, ok := delegate.(EntityHandler)
	if ok {
		facebook.Barrack.SetHandler(handler)
	}
	manager, ok := delegate.(EntityManager)
	if ok {
		facebook._manager = manager
	}
}

func (facebook *Facebook) SetManager(manager EntityManager) {
	facebook._manager = manager
}
func (facebook *Facebook) Manager() EntityManager {
	return facebook._manager
}

//-------- EntityManager

func (facebook *Facebook) GetCurrentUser() User {
	return facebook.Manager().GetCurrentUser()
}

func (facebook *Facebook) CheckDocument(doc Document) bool {
	return facebook.Manager().CheckDocument(doc)
}

func (facebook *Facebook) SaveDocument(doc Document) bool {
	return facebook.Manager().SaveDocument(doc)
}

func (facebook *Facebook) SaveMeta(meta Meta, identifier ID) bool {
	return facebook.Manager().SaveMeta(meta, identifier)
}

func (facebook *Facebook) SaveMembers(members []ID, group ID) bool {
	return facebook.Manager().SaveMembers(members, group)
}

func (facebook *Facebook) IsFounder(member ID, group ID) bool {
	return facebook.Manager().IsFounder(member, group)
}

func (facebook *Facebook) IsOwner(member ID, group ID) bool {
	return facebook.Manager().IsOwner(member, group)
}
