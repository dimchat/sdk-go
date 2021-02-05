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

type IFacebook interface {

	/**
	 *  Save meta for entity ID (must verify first)
	 *
	 * @param meta - entity meta
	 * @param identifier - entity ID
	 * @return true on success
	 */
	SaveMeta(meta Meta, identifier ID) bool

	/**
	 *  Save entity document with ID (must verify first)
	 *
	 * @param doc - entity document
	 * @return true on success
	 */
	SaveDocument(doc Document) bool

	/**
	 *  Save members of group
	 *
	 * @param members - member ID list
	 * @param group - group ID
	 * @return true on success
	 */
	SaveMembers(members []ID, group ID) bool
}

type Facebook struct {
	Barrack
	IFacebook
}

/**
 *  Get current user (for signing and sending message)
 *
 * @return User
 */
func (facebook *Facebook) GetCurrentUser() *User {
	users := facebook.GetLocalUsers()
	if users == nil || len(users) == 0 {
		return nil
	}
	return users[0]
}

/**
 *  Document checking
 *
 * @param doc - entity document
 * @return true on accepted
 */
func (facebook *Facebook) CheckDocument(doc Document) bool {
	identifier := doc.ID()
	if identifier == nil {
		return false
	}
	// NOTICE: if this is a bulletin document for group,
	//             verify it with the group owner's meta.key
	//         else (this is a visa document for user)
	//             verify it with the user's meta.key
	var meta Meta
	if identifier.IsGroup() {
		// check by owner
		owner := facebook.GetOwner(identifier)
		if owner == nil {
			if identifier.Type() == POLYLOGUE {
				// NOTICE: if this is a polylogue document,
				//             verify it with the founder's meta.key
				//             (which equals to the group's meta.key)
				meta = facebook.GetMeta(identifier)
			} else {
				// FIXME: owner not found for this group
				return false
			}
		} else {
			meta = facebook.GetMeta(owner)
		}
	} else {
		meta = facebook.GetMeta(identifier)
	}
	return meta != nil && doc.Verify(meta.Key())
}

//-------- group membership

func (facebook *Facebook) IsFounder(member ID, group ID) bool {
	gMeta := facebook.GetMeta(group)
	mMeta := facebook.GetMeta(member)
	return gMeta.MatchKey(mMeta.Key())
}

func (facebook *Facebook) IsOwner(member ID, group ID) bool {
	if group.Type() == POLYLOGUE {
		return facebook.IsFounder(member, group)
	}
	panic("only Polylogue so far")
	return false
}

func (facebook *Facebook) CreateUser(identifier ID) *User {
	if identifier.IsBroadcast() {
		// create user 'anyone@anywhere'
		return NewUser(identifier)
	}
	// make sure meta exists
	// TODO: make sure visa key exists before calling this
	// check user type
	network := identifier.Type()
	if network == MAIN || network == BTCMain {
		return NewUser(identifier)
	}
	var user interface{}
	if network == ROBOT {
		user = NewRobot(identifier)
	} else if network == STATION {
		user = NewStation(identifier, "", 0)
	}
	return user.(*User)
}

func (facebook *Facebook) CreateGroup(identifier ID) *Group {
	if identifier.IsBroadcast() {
		// create group 'everyone@everywhere'
		return NewGroup(identifier)
	}
	// make sure meta exists
	// check group type
	network := identifier.Type()
	var group interface{}
	if network == POLYLOGUE {
		group = NewPolylogue(identifier)
	} else if network == CHATROOM {
		group = NewChatroom(identifier)
	} else if network == PROVIDER {
		group = NewServiceProvider(identifier)
	}
	return group.(*Group)
}
