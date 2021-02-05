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
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Simple group chat
 */
type Polylogue struct {
	Group
}

func NewPolylogue(identifier ID) *Polylogue {
	return new(Polylogue).Init(identifier)
}

func (group *Polylogue) Init(identifier ID) *Polylogue {
	if group.Group.Init(identifier) != nil {
	}
	return group
}

func (group *Polylogue) GetOwner() ID {
	owner := group.Group.GetOwner()
	if owner == nil {
		// polylogue's owner is its founder
		owner = group.GetFounder()
	}
	return owner
}

/**
 *  Big group with admins
 */
type Chatroom struct {
	Group
}

func NewChatroom(identifier ID) *Chatroom {
	return new(Chatroom).Init(identifier)
}

func (group *Chatroom) Init(identifier ID) *Chatroom {
	if group.Group.Init(identifier) != nil {
	}
	return group
}

func (group Chatroom) DataSource() ChatroomDataSource {
	delegate := group.Group.DataSource()
	return delegate.(ChatroomDataSource)
}

func (group Chatroom) GetAdmins() []ID {
	return group.DataSource().GetAdmins(group.ID())
}

/**
 *  This interface is for getting information for chatroom
 *  Chatroom admins should be set complying with the consensus algorithm
 */
type ChatroomDataSource interface {
	GroupDataSource

	/**
	 *  Get all admins in the chatroom
	 *
	 * @param chatroom - chatroom ID
	 * @return admin ID list
	 */
	GetAdmins(group ID) []ID
}
