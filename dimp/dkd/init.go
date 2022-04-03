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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

func NewReceiptCommand(text string, env Envelope, sn uint64, signature []byte) ReceiptCommand {
	cmd := new(BaseReceiptCommand)
	if ValueIsNil(env) {
		cmd.InitWithMessage(text)
	} else {
		cmd.InitWithEnvelope(env, sn, text)
	}
	if !ValueIsNil(signature) {
		cmd.SetSignature(signature)
	}
	return cmd
}

//
//  Handshake command factories
//
func HandshakeCommandStart() HandshakeCommand {
	return new(BaseHandshakeCommand).InitWithMessage("", "")
}

func HandshakeCommandRestart(session string) HandshakeCommand {
	return new(BaseHandshakeCommand).InitWithMessage("", session)
}

func HandshakeCommandAgain(session string) HandshakeCommand {
	return new(BaseHandshakeCommand).InitWithMessage("DIM?", session)
}

func HandshakeCommandSuccess(session string) HandshakeCommand {
	return new(BaseHandshakeCommand).InitWithMessage("DIM!", session)
}

//
//  Mute command factories
//
func MuteCommandQuery() MuteCommand {
	return new(BaseMuteCommand).InitWithList(nil)
}

func MuteCommandRespond(list []ID) MuteCommand {
	return new(BaseMuteCommand).InitWithList(list)
}

//
//  Block command factories
//
func BlockCommandQuery() BlockCommand {
	return new(BaseBlockCommand).InitWithList(nil)
}

func BlockCommandRespond(list []ID) BlockCommand {
	return new(BaseBlockCommand).InitWithList(list)
}

/**
 *  Register command factories
 */
func RegisterExtraCommandFactories() {
	CommandSetFactory(RECEIPT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseReceiptCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(HANDSHAKE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseHandshakeCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(LOGIN, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseLoginCommand)
		cmd.Init(dict)
		return cmd
	}))

	CommandSetFactory(MUTE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseMuteCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(BLOCK, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseBlockCommand)
		cmd.Init(dict)
		return cmd
	}))

	// storage (contacts, private_key)
	CommandSetFactory(STORAGE, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseStorageCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(CONTACTS, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseStorageCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(PRIVATE_KEY, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseStorageCommand)
		cmd.Init(dict)
		return cmd
	}))
}

func init() {
	//
	//  Register core factories
	//
	RegisterContentFactories()
	RegisterCommandFactories()

	//
	//  Register extra command factories
	//
	RegisterExtraCommandFactories()
}
