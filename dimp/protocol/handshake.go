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
	"fmt"
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

func (state HandshakeState) String() string {
	switch state {
	case HandshakeInit:
		return "HandshakeInit"
	case HandshakeStart:
		return "HandshakeStart"
	case HandshakeAgain:
		return "HandshakeAgain"
	case HandshakeRestart:
		return "HandshakeRestart"
	case HandshakeSuccess:
		return "HandshakeSuccess"
	default:
		return fmt.Sprintf("HandshakeState(%d)", state)
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
type HandshakeCommand interface {
	Command

	Message() string
	Session() string
	State() HandshakeState
}
