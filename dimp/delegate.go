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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Cipher Key Delegate
 *
 *  1. get symmetric key for sending message;
 *  2. cache symmetric key for reusing.
 */
type CipherKeyDelegate interface {

	/**
	 *  Get cipher key for encrypt message from 'sender' to 'receiver'
	 *
	 * @param sender - from where (user or contact ID)
	 * @param receiver - to where (contact or user/group ID)
	 * @param generate - generate when key not exists
	 * @return cipher key
	 */
	GetCipherKey(sender, receiver ID, generate bool) SymmetricKey

	/**
	 *  Cache cipher key for reusing, with the direction (from 'sender' to 'receiver')
	 *
	 * @param sender - from where (user or contact ID)
	 * @param receiver - to where (contact or user/group ID)
	 * @param key - cipher key
	 */
	CacheCipherKey(sender, receiver ID, key SymmetricKey)
}

type Cleanable interface {

	/**
	 *  Remove foreign pointers to break circular references
	 */
	Clean()
}

func Cleanup(obj interface{}) {
	target, ok := obj.(Cleanable)
	if ok && target != nil {
		target.Clean()
	}
}

/**
 *  Messenger Shadow
 *  ~~~~~~~~~~~~~~~~
 *
 *  delegate for messenger
 */
type Helper interface {
	Cleanable

	Facebook() IFacebook
	Messenger() IMessenger
}

type TwinsHelper struct {
	_facebook IFacebook
	_messenger IMessenger
}

func (helper *TwinsHelper) Init(facebook IFacebook, messenger IMessenger) Helper {
	helper._facebook = facebook
	helper._messenger = messenger
	return helper
}

//-------- ICleanable

func (helper *TwinsHelper) Clean() {
	// remove the twins
	helper._messenger = nil
	helper._facebook = nil
}

//-------- IHelper

func (helper *TwinsHelper) Facebook() IFacebook {
	return helper._facebook
}

func (helper *TwinsHelper) Messenger() IMessenger {
	return helper._messenger
}
