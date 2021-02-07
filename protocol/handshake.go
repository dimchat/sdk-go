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
)

type HandshakeState uint8

const (
	HandshakeInit HandshakeState = iota
	HandshakeStart    // C -> S, without session key(or session expired)
	HandshakeAgain    // S -> C, with new session key
	HandshakeRestart  // C -> S, with new session key
	HandshakeSuccess  // S -> C, handshake accepted
)

func getState(message string, session string) HandshakeState {
	// check message text
	if message == "" {
		return HandshakeInit
	}
	if message == "DIM!" || message == "OK!" {
		return HandshakeSuccess
	}
	if message == "DIM?" {
		return HandshakeAgain
	}
	// check session key
	if session == "" {
		return HandshakeStart
	} else {
		return HandshakeRestart
	}
}

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "handshake",    // command name
 *      message : "Hello world!",
 *      session : "{SESSION_KEY}" // session key
 *  }
 */
type HandshakeCommand struct {
	BaseCommand

	_message string
	_session string
	_state HandshakeState
}

func (cmd *HandshakeCommand) Init(dict map[string]interface{}) *HandshakeCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// message text
		text := dict["message"]
		if text == nil {
			cmd._message = ""
		} else {
			cmd._message = text.(string)
		}
		// session key
		key := dict["session"]
		if key == nil {
			cmd._session = ""
		} else {
			cmd._session = key.(string)
		}
		// check state
		cmd._state = getState(cmd._message, cmd._session)
	}
	return cmd
}

func (cmd *HandshakeCommand) InitWithMessage(message string, session string) *HandshakeCommand {
	if cmd.BaseCommand.InitWithCommand(HANDSHAKE) != nil {
		// message text
		if message == "" {
			message = "Hello world!"
		}
		cmd.Set("message", message)
		cmd._message = message
		// session key
		if session != "" {
			cmd.Set("session", session)
		}
		cmd._session = session
		// check state
		cmd._state = getState(message, session)
	}
	return cmd
}

//
//  Factories
//
func HandshakeCommandStart() *HandshakeCommand {
	return new(HandshakeCommand).InitWithMessage("", "")
}

func HandshakeCommandRestart(session string) *HandshakeCommand {
	return new(HandshakeCommand).InitWithMessage("", session)
}

func HandshakeCommandAgain(session string) *HandshakeCommand {
	return new(HandshakeCommand).InitWithMessage("DIM?", session)
}

func HandshakeCommandSuccess(session string) *HandshakeCommand {
	return new(HandshakeCommand).InitWithMessage("DIM!", session)
}
