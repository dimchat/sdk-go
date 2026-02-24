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
package sdk

import . "github.com/dimchat/dkd-go/protocol"

// Processor defines the interface for processing messages at all format levels
//
// Handles message processing at binary, signed, encrypted, plaintext, and content levels
// Each method returns multiple response messages to support complex workflows (replies, broadcasts, etc.)
type Processor interface {

	// ProcessPackage processes raw binary network data and returns response data packages
	//
	// Orchestrates full inbound flow: Deserialize → Verify → Decrypt → Process → Encrypt → Sign → Serialize
	//
	// Parameters:
	//   - data - Raw binary data package received from network
	// Returns: Slice of binary response data packages (empty slice if no response)
	ProcessPackage(data []byte) [][]byte

	// ProcessReliableMessage processes a signed ReliableMessage and returns response messages
	//
	// Typically verifies the message first before delegating to ProcessSecureMessage
	//
	// Parameters:
	//   - rMsg - Signed ReliableMessage to process
	// Returns: Slice of signed ReliableMessage responses (empty slice if no response)
	ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage

	// ProcessSecureMessage processes an encrypted SecureMessage and returns response messages
	//
	// Typically decrypts the message first before delegating to ProcessInstantMessage
	//
	// Parameters:
	//   - sMsg - Encrypted SecureMessage to process
	//   - rMsg - Original signed ReliableMessage (for context/metadata)
	// Returns: Slice of encrypted SecureMessage responses (empty slice if no response)
	ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage

	// ProcessInstantMessage processes a plaintext InstantMessage and returns response messages
	//
	// Extracts content and delegates to ProcessContent for business logic processing
	//
	// Parameters:
	//   - iMsg - Plaintext InstantMessage to process
	//   - rMsg - Original signed ReliableMessage (for context/metadata)
	// Returns: Slice of plaintext InstantMessage responses (empty slice if no response)
	ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage

	// ProcessContent processes raw message content and returns response contents
	//
	// Core business logic handler (e.g., command execution, message routing, content transformation)
	//
	// Parameters:
	//   - content - Raw message content to process (text/file/command/...)
	//   - rMsg    - Original signed ReliableMessage (for sender/receiver context)
	// Returns: Slice of response Content objects (empty slice if no response)
	ProcessContent(content Content, rMsg ReliableMessage) []Content
}
