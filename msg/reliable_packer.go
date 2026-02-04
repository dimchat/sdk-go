/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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

import . "github.com/dimchat/dkd-go/protocol"

type NetworkMessagePacker struct {
	//ReliableMessagePacker

	_transceiver ReliableMessageDelegate
}

func (packer NetworkMessagePacker) Init(messenger ReliableMessageDelegate) ReliableMessagePacker {
	packer._transceiver = messenger
	return packer
}

func (packer NetworkMessagePacker) Delegate() ReliableMessageDelegate {
	return packer._transceiver
}

// Override
func (packer NetworkMessagePacker) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	transceiver := packer.Delegate()
	if transceiver == nil {
		panic("reliable message delegate not found")
		return nil
	}

	//
	//  0. Decode 'message.data' to encrypted content data
	//
	ciphertext := rMsg.Data()
	if ciphertext == nil || ciphertext.IsEmpty() {
		//panic("failed to decode message data")
		return nil
	}

	//
	//  1. Decode 'message.signature' from String (Base64)
	//
	signature := rMsg.Signature()
	if signature == nil || signature.IsEmpty() {
		//panic("failed to decode message signature")
		return nil
	}

	//
	//  2. Verify the message data and signature with sender's public key
	//
	ok := transceiver.VerifyDataSignature(ciphertext.Bytes(), signature.Bytes(), rMsg)
	if !ok {
		//panic("message signature not match")
		return nil
	}

	// OK, pack message
	info := rMsg.CopyMap(false)
	delete(info, "signature")
	return ParseSecureMessage(info)
}
