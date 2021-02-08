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
package protocol

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

const MUTE = "mute"

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "mute",
 *      list    : []      // mute-list
 *  }
 */
type MuteCommand struct {
	BaseCommand

	// mute-list
	_list []ID
}

func (cmd *MuteCommand) Init(dict map[string]interface{}) *MuteCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._list = nil
	}
	return cmd
}

func (cmd *MuteCommand) InitWithList(list []ID) *MuteCommand {
	if cmd.BaseCommand.InitWithCommand(BLOCK) != nil {
		if list != nil {
			cmd.SetList(list)
		}
	}
	return cmd
}

func (cmd *MuteCommand) List() []ID {
	if cmd._list == nil {
		list := cmd.Get("list")
		if list != nil {
			cmd._list = IDConvert(list)
		}
	}
	return cmd._list
}

func (cmd *MuteCommand) SetList(list []ID) {
	if list == nil {
		cmd.Set("list", nil)
	} else {
		cmd.Set("list", IDRevert(list))
	}
	cmd._list = list
}

//
//  Factories
//
func MuteCommandQuery() *MuteCommand {
	return new(MuteCommand).InitWithList(nil)
}

func MuteCommandRespond(list []ID) *MuteCommand {
	return new(MuteCommand).InitWithList(list)
}
