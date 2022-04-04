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
package cpu

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrQueryNotAllowed = "Sorry, you are not allowed to query this group."
)

/**
 *  Group command: "query"
 *  ~~~~~~~~~~~~~~~~~~~~~~
 */
type QueryCommandProcessor struct {
	GroupCommandProcessor
}

func NewQueryCommandProcessor(facebook IFacebook, messenger IMessenger) *QueryCommandProcessor {
	cpu := new(QueryCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//-------- ICommandProcessor

func (gpu *QueryCommandProcessor) Process(content Content, rMsg ReliableMessage) []Content {
	cmd, _ := content.(GroupCommand)
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		return gpu.RespondText(StrGroupEmpty, group)
	}

	// 1. check permission
	sender := rMsg.Sender()
	if !contains(members, sender) {
		// not a member? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || !contains(assistants, sender) {
			return gpu.RespondText(StrQueryNotAllowed, group)
		}
	}

	// 2. respond
	user := facebook.GetCurrentUser()
	if user.ID().Equal(owner) {
		return gpu.RespondContent(NewResetCommand(group, members))
	} else {
		return gpu.RespondContent(NewInviteCommand(group, members))
	}
}
