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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/cpu"
)

/**
 *  Group command: RESET
 */
type ResetCommandProcessor struct {
	GroupCommandProcessor
}

func (gpu *ResetCommandProcessor) temporarySave(cmd *GroupCommand, sender ID) Content {
	facebook := gpu.Facebook()
	group := cmd.Group()
	// check whether the owner contained in the new members
	newMembers := gpu.GetMembers(cmd)
	for _, item := range newMembers {
		if facebook.IsOwner(item, group) {
			// it's a full list, save it now
			if facebook.SaveMembers(newMembers, group) {
				if sender.Equal(item) == false {
					// NOTICE: to prevent counterfeit,
					//         query the owner for newest member-list
					query := new(QueryCommand).InitWithGroup(group)
					gpu.Messenger().SendContent(nil, item, query, nil, 1)
				}
			}
			// response (no need to respond this group command)
			return nil
		}
	}
	// NOTICE: this is a partial member-list
	//         query the sender for full-list
	return new(QueryCommand).InitWithGroup(group)
}

func (gpu *ResetCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) Content {
	//gCmd := cmd.(*GroupCommand)
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		// FIXME: group info lost?
		// FIXME: how to avoid strangers impersonating group member?
		return gpu.temporarySave(cmd.(*GroupCommand), rMsg.Sender())
	}

	// 1. check permission
	sender := rMsg.Sender()
	if owner.Equal(sender) == false {
		// not the owner? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
			panic(sender.String() + " is not the owner/assistant" +
				" of group " + group.String() + ", cannot reset members")
			return nil
		}
	}

	// 2. resetting members
	newMembers := gpu.GetMembers(cmd.(*GroupCommand))
	if newMembers == nil || len(newMembers) == 0 {
		panic("group command error")
		return nil
	}
	// 2.1. check owner
	if contains(owner, newMembers) == false {
		panic("cannot expel owner(" + owner.String() + ") of group: " + group.String())
		return nil
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
