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
	. "github.com/dimchat/sdk-go/dimp/cpu"
)

/**
 *  Group command: QUIT
 */
type QuitCommandProcessor struct {
	GroupCommandProcessor
}

func (gpu *QuitCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) Content {
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		panic("group not ready: " + group.String())
		return nil
	}

	// 1. check permission
	sender := rMsg.Sender()
	if owner.Equal(sender) {
		panic("owner cannot quit: " + sender.String() + " -> " + group.String())
		return nil
	}
	assistants := facebook.GetAssistants(group)
	if assistants != nil && contains(sender, assistants) {
		panic("assistant cannot quit: " + sender.String() + " -> " + group.String())
		return nil
	}

	// 2. remove the sender from group members
	if contains(sender, members) {
		members = remove(members, sender)
		facebook.SaveMembers(members, group)
	}

	// 3. response (no need to response this group command)
	return nil
}
