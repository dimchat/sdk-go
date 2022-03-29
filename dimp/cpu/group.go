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
	"fmt"
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrGroupEmpty = "Group empty."
	FmtGrpCmdNotSupport = "Group command (name: %s) not support yet!"

	StrInviteCmdError = "Invite command error."
	StrInviteNotAllowed = "Sorry, you are not allowed to invite new members into this group."

	StrExpelCmdError = "Expel command error."
	StrExpelNotAllowed = "Sorry, you are not allowed to expel member from this group."
	StrCannotExpelOwner = "Group owner cannot be expelled."

	StrOwnerCannotQuit = "Sorry, group owner cannot quit."
	StrAssistantCannotQuit = "Sorry, group assistant cannot quit."

	StrResetCmdError = "Reset command error."
	StrResetNotAllowed = "Sorry, you are not allowed to reset this group."

	StrQueryNotAllowed = "Sorry, you are not allowed to query this group."
)

type GroupCommandProcessor struct {
	HistoryCommandProcessor
}

func NewGroupCommandProcessor(facebook IFacebook, messenger IMessenger) *GroupCommandProcessor {
	cpu := new(GroupCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (gpu *GroupCommandProcessor) GetMembers(cmd GroupCommand) []ID {
	// get from members
	members := cmd.Members()
	if members == nil {
		// get from 'member'
		member := cmd.Member()
		if member != nil {
			members = []ID{member}
		}
	}
	return members
}

//-------- ICommandProcessor

func (gpu *GroupCommandProcessor) Execute(cmd Command, _ ReliableMessage) []Content {
	text := fmt.Sprintf(FmtGrpCmdNotSupport, cmd.CommandName())
	return gpu.RespondText(text, cmd.Group())
}

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

func (gpu *InviteCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) []Content {
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
	if contains(sender, members) == false {
		// not a member? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
			return gpu.RespondText(StrInviteNotAllowed, group)
		}
	}

	// 2. inviting members
	inviteList := gpu.GetMembers(cmd.(GroupCommand))
	if inviteList == nil || len(inviteList) == 0 {
		return gpu.RespondText(StrInviteCmdError, group)
	}
	// 2.1. check for reset
	if sender.Equal(owner) && contains(owner, inviteList) {
		// NOTICE: owner invites owner?
		//         it means this should be a 'reset' command
		return gpu.temporarySave(cmd.(GroupCommand), rMsg.Sender())
	}
	// 2.2. build invited-list
	addedList := make([]string, 0, len(inviteList))
	for _, item := range inviteList {
		if contains(item, members) {
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

/**
 *  Group command: "expel"
 *  ~~~~~~~~~~~~~~~~~~~~~~
 */
type ExpelCommandProcessor struct {
	GroupCommandProcessor
}

func NewExpelCommandProcessor(facebook IFacebook, messenger IMessenger) *ExpelCommandProcessor {
	cpu := new(ExpelCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//-------- ICommandProcessor

func (gpu *ExpelCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) []Content {
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
	if owner.Equal(sender) == false {
		// not the owner? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
			return gpu.RespondText(StrExpelNotAllowed, group)
		}
	}

	// 2. expelling members
	expelList := gpu.GetMembers(cmd.(GroupCommand))
	if expelList == nil || len(expelList) == 0 {
		return gpu.RespondText(StrExpelCmdError, group)
	}
	// 2.1. check owner
	if contains(owner, expelList) {
		return gpu.RespondText(StrCannotExpelOwner, group)
	}
	// 2.2. build expelled-list
	removedList := make([]string, 0, len(expelList))
	for _, item := range expelList {
		if contains(item, members) == false {
			continue
		}
		// removing member found
		removedList = append(removedList, item.String())
		members = remove(members, item)
	}
	// 2.3. do expel
	if len(removedList) > 0 {
		if facebook.SaveMembers(members, group) {
			cmd.Set("removed", removedList)
		}
	}

	// 3. response (no need to response this group command)
	return nil
}

/**
 *  Group command: "quit"
 *  ~~~~~~~~~~~~~~~~~~~~~
 */
type QuitCommandProcessor struct {
	GroupCommandProcessor
}

func NewQuitCommandProcessor(facebook IFacebook, messenger IMessenger) *QuitCommandProcessor {
	cpu := new(QuitCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (gpu *QuitCommandProcessor) RemoveAssistant(cmd QuitCommand, _ ReliableMessage) []Content {
	// NOTICE: group assistant should be retried by the owner
	return gpu.RespondText(StrAssistantCannotQuit, cmd.Group())
}

//-------- ICommandProcessor

func (gpu *QuitCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) []Content {
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
	if owner.Equal(sender) {
		return gpu.RespondText(StrOwnerCannotQuit, group)
	}
	assistants := facebook.GetAssistants(group)
	if assistants != nil && contains(sender, assistants) {
		return gpu.RemoveAssistant(cmd.(QuitCommand), rMsg)
	}

	// 2. remove the sender from group members
	if contains(sender, members) {
		members = remove(members, sender)
		facebook.SaveMembers(members, group)
	}

	// 3. response (no need to response this group command)
	return nil
}

/**
 *  Group command: "reset"
 *  ~~~~~~~~~~~~~~~~~~~~~~
 */
type ResetCommandProcessor struct {
	GroupCommandProcessor
}

func NewResetCommandProcessor(facebook IFacebook, messenger IMessenger) *ResetCommandProcessor {
	cpu := new(ResetCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (gpu *ResetCommandProcessor) QueryOwner(owner, group ID) {
	// TODO: send QueryCommand to the owner
}

func (gpu *ResetCommandProcessor) temporarySave(cmd GroupCommand, sender ID) []Content {
	facebook := gpu.Facebook()
	group := cmd.Group()
	// check whether the owner contained in the new members
	newMembers := gpu.GetMembers(cmd)
	if newMembers == nil || len(newMembers) == 0 {
		return gpu.RespondText(StrResetCmdError, group)
	}
	for _, item := range newMembers {
		if facebook.GetMeta(item) == nil {
			// TODO: waiting for member's meta?
			continue
		} else if facebook.IsOwner(item, group) == false {
			// not owner, skip it
			continue
		}
		// it's a full list, save it now
		if facebook.SaveMembers(newMembers, group) {
			if sender.Equal(item) == false {
				// NOTICE: to prevent counterfeit,
				//         query the owner for newest member-list
				gpu.QueryOwner(item, group)
			}
		}
		// response (no need to respond this group command)
		return nil
	}
	// NOTICE: this is a partial member-list
	//         query the sender for full-list
	return gpu.RespondContent(NewQueryCommand(group))
}

//-------- ICommandProcessor

func (gpu *ResetCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) []Content {
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		// FIXME: group info lost?
		// FIXME: how to avoid strangers impersonating group member?
		return gpu.temporarySave(cmd.(GroupCommand), rMsg.Sender())
	}

	// 1. check permission
	sender := rMsg.Sender()
	if owner.Equal(sender) == false {
		// not the owner? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
			return gpu.RespondText(StrResetNotAllowed, group)
		}
	}

	// 2. resetting members
	newMembers := gpu.GetMembers(cmd.(GroupCommand))
	if newMembers == nil || len(newMembers) == 0 {
		return gpu.RespondText(StrResetCmdError, group)
	}
	// 2.1. check owner
	if contains(owner, newMembers) == false {
		return gpu.RespondText(StrResetCmdError, group)
	}
	// 2.2. build expelled-list
	removedList := make([]string, 0, len(members))
	for _, item := range members {
		if contains(item, newMembers) {
			continue
		}
		// removing member found
		removedList = append(removedList, item.String())
	}
	// 2.3. build invited-list
	addedList := make([]string, 0, len(newMembers))
	for _, item := range newMembers {
		if contains(item, members) {
			continue
		}
		// adding member found
		addedList = append(addedList, item.String())
	}
	// 2.4. do reset
	if len(addedList) > 0 || len(removedList) > 0 {
		if facebook.SaveMembers(members, group) {
			if len(addedList) > 0 {
				cmd.Set("added", addedList)
			}
			if len(removedList) > 0 {
				cmd.Set("removed", removedList)
			}
		}
	}

	// 3. response (no need to response this group command)
	return nil
}

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

func (gpu *QueryCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) []Content {
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
	if contains(sender, members) == false {
		// not a member? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
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
