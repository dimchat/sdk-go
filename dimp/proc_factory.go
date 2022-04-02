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

type ProcessorCreator interface {

	/**
	 *  Create content processor with type
	 *
	 * @param msgType - content type
	 * @return ContentProcessor
	 */
	CreateContentProcessor(msgType ContentType) ContentProcessor

	/**
	 *  Create command processor with name
	 *
	 * @param msgType - content type
	 * @param cmdName - command name
	 * @return CommandProcessor
	 */
	CreateCommandProcessor(msgType ContentType, cmdName string) CommandProcessor
}

type ProcessorFactory interface {

	/**
	 *  Get processor for content
	 */
	GetProcessor(content Content) ContentProcessor

	/**
	 *  Get processor for content type
	 */
	GetContentProcessor(msgType ContentType) ContentProcessor

	/**
	 *  Get processor for command name
	 */
	GetCommandProcessor(msgType ContentType, cmdName string) CommandProcessor
}

type CPUFactory struct {
	TwinsHelper

	_creator ProcessorCreator

	_contentProcessors map[ContentType]ContentProcessor
	_commandProcessors map[string]CommandProcessor
}

func (factory *CPUFactory) Init(facebook IFacebook, messenger IMessenger) *CPUFactory {
	if factory.TwinsHelper.Init(facebook, messenger) != nil {
		factory._contentProcessors = make(map[ContentType]ContentProcessor)
		factory._commandProcessors = make(map[string]CommandProcessor)
		factory._creator = nil
	}
	return factory
}

func (factory *CPUFactory) Clean() {
	// remove processor pools
	factory._contentProcessors = nil
	factory._commandProcessors = nil
	// clean creator
	creator := factory.Creator()
	factory.SetCreator(nil)
	Cleanup(creator)
	// clean twins
	factory.TwinsHelper.Clean()
}

func (factory *CPUFactory) Creator() ProcessorCreator {
	return factory._creator
}
func (factory *CPUFactory) SetCreator(self ProcessorCreator) {
	factory._creator = self
}

func (factory *CPUFactory) ContentProcessorByType(msgType ContentType) ContentProcessor {
	return factory._contentProcessors[msgType]
}
func (factory *CPUFactory) SetContentProcessorByTag(msgType ContentType, cpu ContentProcessor) {
	factory._contentProcessors[msgType] = cpu
}

func (factory *CPUFactory) CommandProcessorByName(cmdName string) CommandProcessor {
	return factory._commandProcessors[cmdName]
}
func (factory *CPUFactory) SetCommandProcessorByName(cmdName string, cpu CommandProcessor) {
	factory._commandProcessors[cmdName] = cpu
}

//-------- IProcessorFactory

func (factory *CPUFactory) GetProcessor(content Content) ContentProcessor {
	cmd, ok := content.(Command)
	if ok {
		return factory.GetCommandProcessor(content.Type(), cmd.CommandName())
	} else {
		return factory.GetContentProcessor(content.Type())
	}
}

func (factory *CPUFactory) GetContentProcessor(msgType ContentType) ContentProcessor {
	cpu := factory.ContentProcessorByType(msgType)
	if cpu == nil {
		cpu = factory.Creator().CreateContentProcessor(msgType)
		if cpu != nil {
			factory.SetContentProcessorByTag(msgType, cpu)
		}
	}
	return cpu
}

func (factory *CPUFactory) GetCommandProcessor(msgType ContentType, cmdName string) CommandProcessor {
	cpu := factory.CommandProcessorByName(cmdName)
	if cpu == nil {
		cpu = factory.Creator().CreateCommandProcessor(msgType, cmdName)
		if cpu != nil {
			factory.SetCommandProcessorByName(cmdName, cpu)
		}
	}
	return cpu
}
