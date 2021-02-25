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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/protocol"
)

/**
 *  Register command factories
 */
func BuildExtraCommandFactories() {
	CommandRegister(RECEIPT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseReceiptCommand).Init(dict)
	}))
	CommandRegister(HANDSHAKE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseHandshakeCommand).Init(dict)
	}))
	CommandRegister(LOGIN, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseLoginCommand).Init(dict)
	}))

	CommandRegister(MUTE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(MuteCommand).Init(dict)
	}))
	CommandRegister(BLOCK, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BlockCommand).Init(dict)
	}))

	// storage (contacts, private_key)
	CommandRegister(STORAGE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseStorageCommand).Init(dict)
	}))
	CommandRegister(CONTACTS, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseStorageCommand).Init(dict)
	}))
	CommandRegister(PRIVATE_KEY, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseStorageCommand).Init(dict)
	}))
}

func init() {
	BuildExtraCommandFactories()
}
