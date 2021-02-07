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
package dimp

import . "github.com/dimchat/dkd-go/protocol"

/**
 *  Messenger Delegate
 *  ~~~~~~~~~~~~~~~~~~
 */
type MessengerDelegate interface {

	/**
	 *  Upload encrypted data to CDN
	 *
	 * @param data - encrypted file data
	 * @param iMsg - instant message
	 * @return download URL
	 */
	UploadData(data []byte, iMsg InstantMessage) string

	/**
	 *  Download encrypted data from CDN
	 *
	 * @param url - download URL
	 * @param iMsg - instant message
	 * @return encrypted file data
	 */
	DownloadData(url string, iMsg InstantMessage) []byte

	/**
	 *  Send out a data package onto network
	 *
	 * @param data - package data
	 * @param handler - completion handler
	 * @param priority - task priority
	 * @return true on success
	 */
	SendPackage(data []byte, handler MessengerCompletionHandler, priority int) bool
}

/**
 *  Messenger DataSource
 *  ~~~~~~~~~~~~~~~~~~~~
 */
type MessengerDataSource interface {

	/**
	 * Save the message into local storage
	 *
	 * @param iMsg - instant message
	 * @return true on success
	 */
	SaveMessage(iMsg InstantMessage) bool

	/**
	 *  Suspend the sending message for the receiver's meta & visa,
	 *  or group meta when received new message
	 *
	 * @param iMsg - instant message to be sent
	 */
	SuspendInstantMessage(iMsg InstantMessage)

	/**
	 *  Suspend the received message for the sender's meta
	 *
	 * @param rMsg - message received from network
	 */
	SuspendReliableMessage(rMsg ReliableMessage)
}

/**
 *  Messenger Completion Handler
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
type MessengerCompletionHandler interface {

	OnSuccess()

	OnFailed(err error)
}

/**
 *  Messenger Callback
 *  ~~~~~~~~~~~~~~~~~~
 */
type MessengerCallback interface {

	OnFinished(result interface{}, err error)
}
