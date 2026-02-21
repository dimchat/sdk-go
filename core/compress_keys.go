/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
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
package sdk

import . "github.com/dimchat/mkm-go/types"

type Shortener interface {

	/**
	 *  Compress Content
	 */
	CompressContent(content StringKeyMap) StringKeyMap
	ExtractContent(content StringKeyMap) StringKeyMap

	/**
	 *  Compress SymmetricKey
	 */
	CompressSymmetricKey(key StringKeyMap) StringKeyMap
	ExtractSymmetricKey(key StringKeyMap) StringKeyMap

	/**
	 *  Compress ReliableMessage
	 */
	CompressReliableMessage(msg StringKeyMap) StringKeyMap
	ExtractReliableMessage(msg StringKeyMap) StringKeyMap
}

type MessageShortener struct {
	//Shortener

	contentShortKeys []string
	cryptoShortKeys  []string
	messageShortKeys []string
}

func NewMessageShortener(contentKeys, cryptoKeys, messageKeys []string) *MessageShortener {
	return &MessageShortener{
		contentShortKeys: contentKeys,
		cryptoShortKeys:  cryptoKeys,
		messageShortKeys: messageKeys,
	}
}

// protected
func (shortener *MessageShortener) MoveKey(from, to string, info StringKeyMap) {
	value := info[from]
	if value != nil {
		delete(info, from)
		info[to] = value
	}
}

// protected
func (shortener *MessageShortener) ShortenKeys(keys []string, info StringKeyMap) {
	size := len(keys)
	for i := 1; i < size; i += 2 {
		shortener.MoveKey(keys[i], keys[i-1], info)
	}
}

// protected
func (shortener *MessageShortener) RestoreKeys(keys []string, info StringKeyMap) {
	size := len(keys)
	for i := 1; i < size; i += 2 {
		shortener.MoveKey(keys[i-1], keys[i], info)
	}
}

// Override
func (shortener *MessageShortener) CompressContent(content StringKeyMap) StringKeyMap {
	shortener.ShortenKeys(shortener.contentShortKeys, content)
	return content
}

// Override
func (shortener *MessageShortener) ExtractContent(content StringKeyMap) StringKeyMap {
	shortener.RestoreKeys(shortener.contentShortKeys, content)
	return content
}

// Override
func (shortener *MessageShortener) CompressSymmetricKey(key StringKeyMap) StringKeyMap {
	shortener.ShortenKeys(shortener.cryptoShortKeys, key)
	return key
}

// Override
func (shortener *MessageShortener) ExtractSymmetricKey(key StringKeyMap) StringKeyMap {
	shortener.RestoreKeys(shortener.cryptoShortKeys, key)
	return key
}

// Override
func (shortener *MessageShortener) CompressReliableMessage(msg StringKeyMap) StringKeyMap {
	shortener.ShortenKeys(shortener.messageShortKeys, msg)
	return msg
}

// Override
func (shortener *MessageShortener) ExtractReliableMessage(msg StringKeyMap) StringKeyMap {
	shortener.RestoreKeys(shortener.messageShortKeys, msg)
	return msg
}
