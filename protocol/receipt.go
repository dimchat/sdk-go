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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
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
type ReceiptCommand struct {
	BaseCommand

	// original message info
	_envelope Envelope
}

func (cmd *ReceiptCommand) Init(dict map[string]interface{}) *ReceiptCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._envelope = nil
	}
	return cmd
}

func (cmd *ReceiptCommand) InitWithMessage(text string) *ReceiptCommand {
	if cmd.BaseCommand.InitWithCommand(RECEIPT) != nil {
		cmd.Set("message", text)
		cmd._envelope = nil
	}
	return cmd
}

func (cmd *ReceiptCommand) InitWithEnvelope(env Envelope, sn uint32, text string) *ReceiptCommand {
	if cmd.BaseCommand.InitWithCommand(RECEIPT) != nil {
		// envelope of the message responding to
		cmd.SetEnvelope(env)
		cmd._envelope = env
		// SerialNumber of the message responding to
		if sn > 0 {
			cmd.Set("sn", sn)
		}
		if text != "" {
			cmd.Set("message", text)
		}
	}
	return cmd
}

func (cmd *ReceiptCommand) Message() string {
	text := cmd.Get("message")
	if text == nil {
		return ""
	}
	return text.(string)
}

func (cmd *ReceiptCommand) Envelope() Envelope {
	if cmd._envelope == nil {
		env := cmd.Get("envelope")
		if env == nil {
			sender := cmd.Get("sender")
			receiver := cmd.Get("receiver")
			if sender != nil && receiver != nil {
				env = cmd.GetMap(false)
			}
		}
		cmd._envelope = EnvelopeParse(env)
	}
	return cmd._envelope
}
func (cmd *ReceiptCommand) SetEnvelope(env Envelope) {
	if env == nil {
		cmd.Set("sender", nil)
		cmd.Set("receiver", nil)
		//cmd.Set("time", nil)
	} else {
		cmd.Set("sender", env.Sender().String())
		cmd.Set("receiver", env.Receiver().String())
		when := env.Time()
		if when.IsZero() == false {
			cmd.Set("time", when.Unix())
		}
	}
}

func (cmd *ReceiptCommand) Signature() []byte {
	signature := cmd.Get("signature")
	if signature == nil {
		return nil
	}
	return Base64Decode(signature.(string))
}
func (cmd *ReceiptCommand) SetSignature(signature []byte) {
	if signature == nil {
		cmd.Set("signature", nil)
	} else {
		cmd.Set("signature", Base64Encode(signature))
	}
}
