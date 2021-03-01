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
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/dimp"
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,  // the same serial number with the original message
 *
 *      command   : "receipt",
 *      message   : "...",
 *      // -- extra info
 *      sender    : "...",
 *      receiver  : "...",
 *      time      : 0,
 *      signature : "..." // the same signature with the original message
 *  }
 */
func receipt(message string) Command {
	cmd := new(BaseCommand).InitWithCommand(RECEIPT)
	cmd.Set("message", message)
	ObjectRetain(cmd)
	ObjectAutorelease(cmd)
	return cmd
}

//
//  Register content processors
//
func BuildContentProcessors() {
	ContentProcessorRegister(FORWARD, new(ForwardContentProcessor).Init())

	fpu := new(BaseFileContentProcessor).Init()
	ContentProcessorRegister(FILE, fpu)
	ContentProcessorRegister(IMAGE, fpu)
	ContentProcessorRegister(AUDIO, fpu)
	ContentProcessorRegister(VIDEO, fpu)

	ContentProcessorRegister(COMMAND, new(BaseCommandProcessor).Init())
	ContentProcessorRegister(HISTORY, new(HistoryCommandProcessor).Init())

	ContentProcessorRegister(0, new(BaseContentProcessor).Init())
}

//
//  Register command processors
//
func BuildCommandProcessors() {
	CommandProcessorRegister(META, new(MetaCommandProcessor).Init())
	CommandProcessorRegister(DOCUMENT, new(DocumentCommandProcessor).Init())
}

func init() {
	BuildContentProcessors()
	BuildCommandProcessors()
}
