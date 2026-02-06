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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/mkm"
)

type IFacebook interface {
	EntityDelegate
	EntityDataSource
	//UserDataSource
	//GroupDataSource

	SaveMeta(meta Meta, did ID) bool
	SaveDocument(document Document, did ID) bool

	/**
	 *  Select local user for receiver
	 *
	 * @param receiver - user/broadcast ID
	 * @return local user
	 */
	SelectUser(receiver ID) ID

	/**
	 *  Select local user for group members
	 *
	 * @param members - group members
	 * @return local user
	 */
	SelectMember(members []ID) ID
}

type Facebook struct {
	//IFacebook

	// protected
	Barrack Barrack
	// public
	Archivist Archivist
}

//func (facebook *Facebook) Init(archivist Archivist, barrack Barrack) IFacebook {
//	facebook.Barrack = barrack
//	facebook.Archivist = archivist
//	return facebook
//}

// Override
func (facebook *Facebook) SaveMeta(meta Meta, did ID) bool {
	archivist := facebook.Archivist
	return archivist.SaveMeta(meta, did)
}

// Override
func (facebook *Facebook) SaveDocument(document Document, did ID) bool {
	archivist := facebook.Archivist
	return archivist.SaveDocument(document, did)
}

// Override
func (facebook *Facebook) SelectUser(receiver ID) ID {
	archivist := facebook.Archivist
	users := archivist.LocalUsers()
	if users == nil || len(users) == 0 {
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

// Override
func (facebook *Facebook) SelectMember(members []ID) ID {
	archivist := facebook.Archivist
	users := archivist.LocalUsers()
	if users == nil || len(users) == 0 {
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

//-------- IEntityDelegate

// Override
func (facebook *Facebook) GetUser(uid ID) User {
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
func (facebook *Facebook) GetGroup(gid ID) Group {
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
