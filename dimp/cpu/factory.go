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
package cpu

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

type CPFactory struct {
	TwinsHelper

	_creator ContentProcessorCreator

	_contentProcessors map[ContentType]ContentProcessor
	_commandProcessors map[string]ContentProcessor
}

func (factory *CPFactory) Init(facebook IFacebook, messenger IMessenger) *CPFactory {
	if factory.TwinsHelper.Init(facebook, messenger) != nil {
		factory._contentProcessors = make(map[ContentType]ContentProcessor)
		factory._commandProcessors = make(map[string]ContentProcessor)
		factory._creator = nil
	}
	return factory
}

func (factory *CPFactory) Clean() {
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

func (factory *CPFactory) Creator() ContentProcessorCreator {
	return factory._creator
}
func (factory *CPFactory) SetCreator(self ContentProcessorCreator) {
	factory._creator = self
}

func (factory *CPFactory) ContentProcessorByType(msgType ContentType) ContentProcessor {
	return factory._contentProcessors[msgType]
}
func (factory *CPFactory) SetContentProcessorByType(msgType ContentType, cpu ContentProcessor) {
	factory._contentProcessors[msgType] = cpu
}

func (factory *CPFactory) CommandProcessorByName(cmdName string) ContentProcessor {
	return factory._commandProcessors[cmdName]
}
func (factory *CPFactory) SetCommandProcessorByName(cmdName string, cpu ContentProcessor) {
	factory._commandProcessors[cmdName] = cpu
}

//-------- IContentProcessorFactory

func (factory *CPFactory) GetProcessor(content Content) ContentProcessor {
	var cpu ContentProcessor
	msgType := content.Type()
	cmd, ok := content.(Command)
	if ok /*&& cmd != nil */{
		cmdName := cmd.CommandName()
		// command processor
		cpu = factory.GetCommandProcessor(msgType, cmdName)
		if cpu != nil {
			return cpu
		}
		_, ok = content.(GroupCommand)
		if ok {
			// group command processor
			cpu = factory.GetCommandProcessor(msgType, "group")
			if cpu != nil {
				return cpu
			}
		}
	}
	// content processor
	cpu = factory.GetContentProcessor(msgType)
	if cpu == nil {
		// default content processor
		cpu = factory.GetContentProcessor(0)
	}
	return cpu
}

func (factory *CPFactory) GetContentProcessor(msgType ContentType) ContentProcessor {
	cpu := factory.ContentProcessorByType(msgType)
	if cpu == nil {
		cpu = factory.Creator().CreateContentProcessor(msgType)
		if cpu != nil {
			factory.SetContentProcessorByType(msgType, cpu)
		}
	}
	return cpu
}

func (factory *CPFactory) GetCommandProcessor(msgType ContentType, cmdName string) ContentProcessor {
	cpu := factory.CommandProcessorByName(cmdName)
	if cpu == nil {
		cpu = factory.Creator().CreateCommandProcessor(msgType, cmdName)
		if cpu != nil {
			factory.SetCommandProcessorByName(cmdName, cpu)
		}
	}
	return cpu
}
