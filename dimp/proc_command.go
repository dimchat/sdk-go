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

/**
 *  CPU: Command Processing Unit
 */
type CommandProcessor interface {
	ContentProcessor

	Execute(cmd Command, rMsg ReliableMessage) Content
}

//
//  CPU Factories
//
var commandProcessors = make(map[string]CommandProcessor)

func CommandProcessorRegister(command string, cpu CommandProcessor) {
	commandProcessors[command] = cpu
}
func CommandProcessorGet(cmd Command) CommandProcessor {
	return CommandProcessorGetByName(cmd.CommandName())
}
func CommandProcessorGetByName(command string) CommandProcessor {
	return commandProcessors[command]
}

/**
 *  Base Command Processor
 */
type BaseCommandProcessor struct {
	BaseContentProcessor
	CommandProcessor
}

func (cpu *BaseCommandProcessor) Init() *BaseCommandProcessor {
	if cpu.BaseContentProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *BaseCommandProcessor) SetMessenger(messenger IMessenger) {
	cpu.BaseContentProcessor.SetMessenger(messenger)
}

func (cpu *BaseCommandProcessor) Process(content Content, rMsg ReliableMessage) Content {
	cmd, _ := content.(Command)
	// get CPU by command name
	processor := CommandProcessorGet(cmd)
	if processor == nil {
		// check for group command
		_, ok := cmd.(*GroupCommand)
		if ok {
			processor = CommandProcessorGetByName("group")
		}
	}
	if processor == nil {
		processor = cpu
	} else {
		processor.SetMessenger(cpu.Messenger())
	}
	return processor.Execute(cmd, rMsg)
}

func (cpu *BaseCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	text := fmt.Sprintf("Command (name: %s) not support yet!", cmd.CommandName())
	res := NewTextContent(text)
	// check group message
	group := cmd.Group()
	if group != nil {
		res.SetGroup(group)
	}
	return res
}
