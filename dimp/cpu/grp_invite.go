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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrInviteCmdError = "Invite command error."
	StrInviteNotAllowed = "Sorry, you are not allowed to invite new members into this group."
)

/**
 *  Group command: "invite"
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type InviteCommandProcessor struct {
	ResetCommandProcessor
}

func NewInviteCommandProcessor(facebook IFacebook, messenger IMessenger) *InviteCommandProcessor {
	cpu := new(InviteCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//-------- ICommandProcessor

func (gpu *InviteCommandProcessor) Process(content Content, rMsg ReliableMessage) []Content {
	cmd, _ := content.(GroupCommand)
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		// NOTICE: group membership lost?
		//         reset group members
		return gpu.temporarySave(cmd.(GroupCommand), rMsg.Sender())
	}

	// 1. check permission
	sender := rMsg.Sender()
	if !contains(members, sender) {
		// not a member? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || !contains(assistants, sender) {
			return gpu.RespondText(StrInviteNotAllowed, group)
		}
	}

	// 2. inviting members
	inviteList := gpu.GetMembers(cmd.(GroupCommand))
	if inviteList == nil || len(inviteList) == 0 {
		return gpu.RespondText(StrInviteCmdError, group)
	}
	// 2.1. check for reset
	if sender.Equal(owner) && contains(inviteList, owner) {
		// NOTICE: owner invites owner?
		//         it means this should be a 'reset' command
		return gpu.temporarySave(cmd.(GroupCommand), rMsg.Sender())
	}
	// 2.2. build invited-list
	addedList := make([]string, 0, len(inviteList))
	for _, item := range inviteList {
		if contains(members, item) {
			continue
		}
		// new member found
		addedList = append(addedList, item.String())
		members = append(members, item)
	}
	// 2.3. do invite
	if len(addedList) > 0 {
		if facebook.SaveMembers(members, group) {
			cmd.Set("added", addedList)
		}
	}

	// 3. response (no need to response this group command)
	return nil
}
