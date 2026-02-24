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
	. "github.com/dimchat/sdk-go/mkm"
)

// ITwinsHelper defines the interface for bridging Facebook (account) and Messenger (message) systems
//
// Core responsibility: Resolve local user identities for message processing by combining
// account information (Facebook) with message context (Messenger)
type ITwinsHelper interface {

	// SelectLocalUser resolves the local User instance for a given receiver ID
	//
	// Uses Facebook.SelectUser to find the local user ID, then retrieves the User instance
	// from Facebook's Barrack (entity cache)
	//
	// Parameters:
	//   - receiver - Target receiver ID (user/group/broadcast)
	// Returns: Local User instance (nil if no matching local user found)
	SelectLocalUser(receiver ID) User
}

// TwinsHelper is the concrete implementation of the ITwinsHelper interface
//
// Acts as the core integration layer between account management (Facebook)
// and message processing (Messenger) subsystems
type TwinsHelper struct {
	//ITwinsHelper

	// Facebook provides access to account/entity operations
	//
	// Used to resolve local user identities for message processing
	Facebook Facebook

	// Messenger provides access to message processing operations
	//
	// Used to handle message encryption/decryption for resolved local users
	Messenger Messenger
}

func NewTwinsHelper(facebook Facebook, messenger Messenger) *TwinsHelper {
	return &TwinsHelper{
		Facebook:  facebook,
		Messenger: messenger,
	}
}

// protected
func (helper *TwinsHelper) SelectLocalUser(receiver ID) User {
	facebook := helper.Facebook
	var me ID
	if receiver.IsBroadcast() {
		// broadcast message can be decrypted by anyone
		me = facebook.SelectUser(receiver)
	} else if receiver.IsUser() {
		// check local users
		me = facebook.SelectUser(receiver)
	} else if receiver.IsGroup() {
		// check local users for the group members
		members := facebook.GetMembers(receiver)
		// the messenger will check group info before decrypting message,
		// so we can trust that the group's meta & members MUST exist here.
		if len(members) == 0 {
			//panic("failed to get group members")
			return nil
		}
		me = facebook.SelectMember(members)
	} else {
		//panic("unknown receiver")
		return nil
	}
	if me == nil {
		// not for me?
		return nil
	}
	return facebook.GetUser(me)
}
