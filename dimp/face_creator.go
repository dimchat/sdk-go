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
 *  Facebook shadow as EntityCreator
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 * @abstract:
 *      GetLocalUsers() []User
 */
type FacebookCreator struct {
	BarrackShadow
	EntityCreator
}

func (shadow *FacebookCreator) Init(facebook IFacebook) *FacebookCreator {
	if shadow.BarrackShadow.Init(facebook) != nil {
	}
	return shadow
}

//-------- EntityCreator

func (shadow *FacebookCreator) CreateUser(identifier ID) User {
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

func (shadow *FacebookCreator) CreateGroup(identifier ID) Group {
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
		return NewServiceProvider(identifier)
	} else {
		return nil
	}
}
