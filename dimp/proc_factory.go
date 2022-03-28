/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2022 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
)

type ProcessorFactory struct {
	IProcessorFactory
	MessengerHelper

	_creator IProcessorCreator

	_contentProcessors map[uint8]IContentProcessor
	_commandProcessors map[string]ICommandProcessor
}
type IProcessorFactory interface {

	GetProcessor(content Content) IContentProcessor

	GetContentProcessor(msgType uint8) IContentProcessor

	GetCommandProcessor(msgType uint8, cmdName string) ICommandProcessor
}

type IProcessorCreator interface {

	/**
	 *  Create content processor with type
	 *
	 * @param msgType - content type
	 * @return ContentProcessor
	 */
	CreateContentProcessor(msgType uint8) IContentProcessor

	/**
	 *  Create command processor with name
	 *
	 * @param msgType - content type
	 * @param cmdName - command name
	 * @return CommandProcessor
	 */
	CreateCommandProcessor(msgType uint8, cmdName string) ICommandProcessor
}

func (factory *ProcessorFactory) Init(facebook IFacebook, messenger IMessenger) *ProcessorFactory {
	if factory.MessengerHelper.Init(facebook, messenger) != nil {
		factory._creator = nil
		factory._contentProcessors = make(map[uint8]IContentProcessor)
		factory._commandProcessors = make(map[string]ICommandProcessor)
	}
	return factory
}

func (factory *ProcessorFactory) Creator() IProcessorCreator {
	return factory._creator
}
func (factory *ProcessorFactory) SetCreator(self IProcessorCreator) {
	factory._creator = self
}

func (factory *ProcessorFactory) ContentProcessorByType(msgType uint8) IContentProcessor {
	return factory._contentProcessors[msgType]
}
func (factory *ProcessorFactory) SetContentProcessor(msgType uint8, cpu IContentProcessor) {
	factory._contentProcessors[msgType] = cpu
}

func (factory *ProcessorFactory) CommandProcessorByName(cmdName string) ICommandProcessor {
	return factory._commandProcessors[cmdName]
}
func (factory *ProcessorFactory) SetCommandProcessor(cmdName string, cpu ICommandProcessor) {
	factory._commandProcessors[cmdName] = cpu
}

//-------- IProcessorFactory

func (factory *ProcessorFactory) GetProcessor(content Content) IContentProcessor {
	cmd, ok := content.(ICommand)
	if ok {
		return factory.GetCommandProcessor(content.Type(), cmd.CommandName())
	} else {
		return factory.GetContentProcessor(content.Type())
	}
}

func (factory *ProcessorFactory) GetContentProcessor(msgType uint8) IContentProcessor {
	cpu := factory.ContentProcessorByType(msgType)
	if cpu == nil {
		cpu = factory.Creator().CreateContentProcessor(msgType)
		if cpu != nil {
			factory.SetContentProcessor(msgType, cpu)
		}
	}
	return cpu
}

func (factory *ProcessorFactory) GetCommandProcessor(msgType uint8, cmdName string) ICommandProcessor {
	cpu := factory.CommandProcessorByName(cmdName)
	if cpu == nil {
		cpu = factory.Creator().CreateCommandProcessor(msgType, cmdName)
		if cpu != nil {
			factory.SetCommandProcessor(cmdName, cpu)
		}
	}
	return cpu
}
