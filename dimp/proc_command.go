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
	"fmt"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
)

var (
	FmtCmdNotSupport = "Command (name: %s) not support yet!"
)

/**
 *  CPU: Command Processing Unit
 */
type CommandProcessor struct {
	ICommandProcessor
	ContentProcessor
}
type ICommandProcessor interface {
	IContentProcessor
	ICommandProcessorExt
}
type ICommandProcessorExt interface {

	/**
	 *  Execute command
	 *
	 * @param cmd     - command received
	 * @param message - reliable message
	 * @return contents responding to msg.sender
	 */
	Execute(cmd Command, rMsg ReliableMessage) []Content
}

func NewCommandProcessor(facebook IFacebook, messenger IMessenger) * CommandProcessor {
	cpu := new(CommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (cpu *CommandProcessor) Process(content Content, rMsg ReliableMessage) []Content {
	cmd, _ := content.(Command)
	return cpu.Execute(cmd, rMsg)
}

func (cpu *CommandProcessor) Execute(cmd Command, _ ReliableMessage) []Content {
	text := fmt.Sprintf(FmtCmdNotSupport, cmd.CommandName())
	return cpu.RespondText(text, cmd.Group())
}
