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
	IBarrack
	EntityManager
}

/**
 *  Delegate for Entity
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Facebook struct {
	Barrack
	EntityManager
}

func (facebook *Facebook) Init() *Facebook {
	if facebook.Barrack.Init() != nil {
	}
	return facebook
}

func (facebook *Facebook) SetShadow(shadow IFacebook) {
	facebook.Barrack.SetShadow(shadow)
}
func (facebook *Facebook) Shadow() IFacebook {
	return facebook.Barrack.Shadow().(IFacebook)
}

//-------- EntityManager

func (facebook *Facebook) GetCurrentUser() User {
	return facebook.Shadow().GetCurrentUser()
}

func (facebook *Facebook) CheckDocument(doc Document) bool {
	return facebook.Shadow().CheckDocument(doc)
}

func (facebook *Facebook) SaveDocument(doc Document) bool {
	return facebook.Shadow().SaveDocument(doc)
}

func (facebook *Facebook) SaveMeta(meta Meta, identifier ID) bool {
	return facebook.Shadow().SaveMeta(meta, identifier)
}

func (facebook *Facebook) SaveMembers(members []ID, group ID) bool {
	return facebook.Shadow().SaveMembers(members, group)
}

func (facebook *Facebook) IsFounder(member ID, group ID) bool {
	return facebook.Shadow().IsFounder(member, group)
}

func (facebook *Facebook) IsOwner(member ID, group ID) bool {
	return facebook.Shadow().IsOwner(member, group)
}
