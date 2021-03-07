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

type GroupCommandProcessor struct {
	HistoryCommandProcessor
}

func (gpu *GroupCommandProcessor) Init() *GroupCommandProcessor {
	if gpu.HistoryCommandProcessor.Init() != nil {
	}
	return gpu
}

func (gpu *GroupCommandProcessor) Process(content Content, rMsg ReliableMessage) Content {
	cmd, _ := content.(Command)
	// get CPU by command name
	processor := CommandProcessorGet(cmd)
	if processor == nil {
		processor = gpu
	} else {
		processor.SetMessenger(gpu.Messenger())
	}
	return processor.Execute(cmd, rMsg)
}

func (gpu *GroupCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	text := fmt.Sprintf("Group command (name: %s) not support yet!", cmd.CommandName())
	res := NewTextContent(text)
	res.SetGroup(cmd.Group())
	return res
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
