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
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/mkm"
	"reflect"
)

type IFacebook interface {
	IBarrack

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

	/**
	 *  Document checking
	 *
	 * @param doc - entity document
	 * @return true on accepted
	 */
	CheckDocument(doc Document) bool

	/**
	 *  Check whether this user is the group's founder
	 */
	IsFounder(member ID, group ID) bool

	/**
	 *  Check whether this user is the group's owner
	 */
	IsOwner(member ID, group ID) bool

	/**
	 *  Create user with ID
	 */
	CreateUser(identifier ID) User

	/**
	 *  Create group with ID
	 */
	CreateGroup(identifier ID) Group

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return Users with private keys
	 */
	GetLocalUsers() []User

	/**
	 *  Get current user (for signing and sending message)
	 *
	 * @return User with private key
	 */
	GetCurrentUser() User

	/**
	 *  Select local user for receiver
	 *
	 * @param receiver - user/group ID
	 * @return local user
	 */
	SelectLocalUser(receiver ID) User
}

/**
 *  Entity Manager
 *  ~~~~~~~~~~~~~~
 *  Manage meta/document for all users/groups
 *
 * @abstract:
 *      // IEntityDataSource
 *      GetMeta(identifier ID) Meta
 *      GetDocument(identifier ID, docType string) Document
 *      // IUserDataSource
 *      GetContacts(user ID) []ID
 *      GetPrivateKeysForDecryption(user ID) []DecryptKey
 *      GetPrivateKeyForSignature(user ID) SignKey
 *      GetPrivateKeyForVisaSignature(user ID) SignKey
 *
 *      // IFacebookExt
 *      SaveMeta(meta Meta, identifier ID) bool
 *      SaveDocument(doc Document) bool
 *      SaveMembers(members []ID, group ID) bool
 *      GetLocalUsers() []User
 */
type Facebook struct {
	Barrack

	// memory caches
	_users map[ID]User
	_groups map[ID]Group
}

func (facebook *Facebook) Init() *Facebook {
	if facebook.Barrack.Init() != nil {
		facebook._users = make(map[ID]User)
		facebook._groups = make(map[ID]Group)
	}
	return facebook
}

func (facebook *Facebook) Clean() {
	// remove memory caches
	facebook._users = nil
	facebook._groups = nil
	// clean barrack.Source
	self := facebook.Source()
	facebook.SetSource(nil)
	Cleanup(self)
}

func (facebook *Facebook) self() IFacebook {
	self := facebook.Source()
	return self.(IFacebook)
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (facebook *Facebook) ReduceMemory() int {
	finger := 0
	finger = thanos(facebook._users, finger)
	finger = thanos(facebook._groups, finger)
	return finger >> 1
}

func thanos(planet interface{}, finger int) int {
	if reflect.TypeOf(planet).Kind() != reflect.Map {
		panic(planet)
		return finger
	}
	dict := reflect.ValueOf(planet)
	keys := dict.MapKeys()
	for _, key := range keys {
		finger++
		if (finger & 1) == 1 {
			// kill it
			dict.SetMapIndex(key, reflect.Value{})
		}
		// let it go
	}
	return finger
}

func (facebook *Facebook) cacheUser(user User) {
	if user.DataSource() == nil {
		user.SetDataSource(facebook.Source())
	}
	facebook._users[user.ID()] = user
}

func (facebook *Facebook) cacheGroup(group Group) {
	if group.DataSource() == nil {
		group.SetDataSource(facebook.Source())
	}
	facebook._groups[group.ID()] = group
}

//-------- IFacebookExt

//func (facebook *Facebook) SaveMeta(meta Meta, identifier ID) bool {
//	panic("not implemented")
//}
//
//func (facebook *Facebook) SaveDocument(doc Document) bool {
//	panic("not implemented")
//}
//
//func (facebook *Facebook) SaveMembers(members []ID, group ID) bool {
//	panic("not implemented")
//}

func (facebook *Facebook) CheckDocument(doc Document) bool {
	identifier := doc.ID()
	if identifier == nil {
		return false
	}
	// NOTICE: if this is a bulletin document for group,
	//             verify it with the group owner's meta.key
	//         else (this is a visa document for user)
	//             verify it with the user's meta.key
	self := facebook.Source()

	var meta Meta
	if identifier.IsGroup() {
		// check by owner
		owner := self.GetOwner(identifier)
		if owner == nil {
			if identifier.Type() == POLYLOGUE {
				// NOTICE: if this is a polylogue document,
				//             verify it with the founder's meta.key
				//             (which equals to the group's meta.key)
				meta = self.GetMeta(identifier)
			} else {
				// FIXME: owner not found for this group
				return false
			}
		} else {
			meta = self.GetMeta(owner)
		}
	} else {
		meta = self.GetMeta(identifier)
	}
	return meta != nil && doc.Verify(meta.Key())
}

func (facebook *Facebook) IsFounder(member ID, group ID) bool {
	self := facebook.Source()
	gMeta := self.GetMeta(group)
	mMeta := self.GetMeta(member)
	return MetaMatchKey(gMeta, mMeta.Key())
}

func (facebook *Facebook) IsOwner(member ID, group ID) bool {
	if group.Type() == POLYLOGUE {
		return facebook.self().IsFounder(member, group)
	}
	panic("only Polylogue so far")
	return false
}

func (facebook *Facebook) CreateUser(identifier ID) User {
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
	if network == ROBOT {
		return NewRobot(identifier)
	} else if network == STATION {
		return NewStation(identifier, "", 0)
	} else {
		return nil
	}
}

func (facebook *Facebook) CreateGroup(identifier ID) Group {
	if identifier.IsBroadcast() {
		// create group 'everyone@everywhere'
		return NewGroup(identifier)
	}
	// make sure meta exists
	// check group type
	network := identifier.Type()
	if network == POLYLOGUE {
		return NewPolylogue(identifier)
	} else if network == CHATROOM {
		return NewChatroom(identifier)
	} else if network == PROVIDER {
		return NewProvider(identifier)
	} else {
		return nil
	}
}

//func (facebook *Facebook) GetLocalUsers() []User {
//	panic("not implemented")
//}

func (facebook *Facebook) GetCurrentUser() User {
	users := facebook.self().GetLocalUsers()
	if users == nil || len(users) == 0 {
		return nil
	} else {
		return users[0]
	}
}

func (facebook *Facebook) SelectLocalUser(receiver ID) User {
	self := facebook.self()
	users := self.GetLocalUsers()
	if users == nil || len(users) == 0 {
		panic("local users should not be empty")
	} else if receiver.IsBroadcast() {
		// broadcast message can decrypt by anyone, so just return current user
		return users[0]
	}
	if receiver.IsGroup() {
		// group message (recipient not designated)
		members := self.GetMembers(receiver)
		if members == nil || len(members) == 0 {
			// TODO: group not ready, waiting for group info
			return nil
		}
		for _, item := range users {
			if item == nil {
				continue
			}
			for _, member := range members {
				if item.ID().Equal(member) {
					// DISCUSS: set this item to be current user?
					return item
				}
			}
		}
	} else {
		// 1. personal message
		// 2. split group message
		for _, item := range users {
			if item == nil {
				continue
			}
			if item.ID().Equal(receiver) {
				// DISCUSS: set this item to be current user?
				return item
			}
		}
	}
	return nil
}

//-------- IEntityDelegate

func (facebook *Facebook) GetUser(identifier ID) User {
	// 1. get from user cache
	user := facebook._users[identifier]
	if user == nil {
		// 2. create user and cache it
		user = facebook.self().CreateUser(identifier)
		if user != nil {
			facebook.cacheUser(user)
		}
	}
	return user
}

func (facebook *Facebook) GetGroup(identifier ID) Group {
	// 1. get from group cache
	// 1. get from user cache
	group := facebook._groups[identifier]
	if group == nil {
		// 2. create group and cache it
		group = facebook.self().CreateGroup(identifier)
		if group != nil {
			facebook.cacheGroup(group)
		}
	}
	return group
}
