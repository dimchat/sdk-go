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
package dkd

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "block",
 *      list    : []      // block-list
 *  }
 */
type BaseBlockCommand struct {
	BaseCommand

	// block-list
	_list []ID
}

func (cmd *BaseBlockCommand) Init(dict map[string]interface{}) *BaseBlockCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._list = nil
	}
	return cmd
}

func (cmd *BaseBlockCommand) InitWithList(list []ID) *BaseBlockCommand {
	if cmd.BaseCommand.InitWithCommand(BLOCK) != nil {
		if !ValueIsNil(list) {
			cmd.SetBlockList(list)
		}
	}
	return cmd
}

func (cmd *BaseBlockCommand) BlockList() []ID {
	if cmd._list == nil {
		list := cmd.Get("list")
		if list != nil {
			cmd._list = IDConvert(list)
		}
	}
	return cmd._list
}

func (cmd *BaseBlockCommand) SetBlockList(list []ID) {
	if ValueIsNil(list) {
		cmd.Remove("list")
	} else {
		cmd.Set("list", IDRevert(list))
	}
	cmd._list = list
}
