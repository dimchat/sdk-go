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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrGroupEmpty = "Group empty."
	FmtGrpCmdNotSupport = "Group command (name: %s) not support yet!"
)

type GroupCommandProcessor struct {
	HistoryCommandProcessor
}

func NewGroupCommandProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := new(GroupCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//func (cpu *GroupCommandProcessor) Init(facebook IFacebook, messenger IMessenger) ContentProcessor {
//	if cpu.HistoryCommandProcessor.Init(facebook, messenger) != nil {
//	}
//	return cpu
//}

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

//-------- IContentProcessor

func (gpu *GroupCommandProcessor) Process(content Content, _ ReliableMessage) []Content {
	cmd, _ := content.(GroupCommand)
	text := fmt.Sprintf(FmtGrpCmdNotSupport, cmd.CommandName())
	return gpu.RespondText(text, cmd.Group())
}

//
//  Utils for Group command Processing Units
//

func find(list []ID, id ID) int {
	for index, item := range list {
		if id.Equal(item) {
			return index
		}
	}
	return -1
}

func contains(list []ID, id ID) bool {
	return find(list, id) != -1
}

func remove(list []ID, item ID) []ID {
	pos := find(list, item)
	if pos < 0 {
		return list
	} else if pos == 0 {
		return list[1:]
	}
	length := len(list) - 1
	if pos == length {
		return list[:length]
	}
	out := make([]ID, length)
	index := 0
	for ; index < pos; index++ {
		out[index] = list[index]
	}
	for ; index < length; index++ {
		out[index] = list[index+1]
	}
	return out
}
