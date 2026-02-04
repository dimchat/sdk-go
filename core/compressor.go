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

import (
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

type Compressor interface {

	//
	//  Compress Content
	//
	CompressContent(content StringKeyMap, key StringKeyMap) []byte
	ExtractContent(data []byte, key StringKeyMap) StringKeyMap

	//
	//  Compress SymmetricKey
	//
	CompressSymmetricKey(key StringKeyMap) []byte
	ExtractSymmetricKey(data []byte) StringKeyMap

	//
	//  Compress ReliableMessage
	//
	CompressReliableMessage(msg StringKeyMap) []byte
	ExtractReliableMessage(data []byte) StringKeyMap
}

type MessageCompressor struct {
	//Compressor

	shortener Shortener
}

func (compressor *MessageCompressor) Init(shortener Shortener) Compressor {
	compressor.shortener = shortener
	return compressor
}

// Override
func (compressor *MessageCompressor) CompressContent(content StringKeyMap, _ StringKeyMap) []byte {
	content = compressor.shortener.CompressContent(content)
	return jsonEncode(content)
}

// Override
func (compressor *MessageCompressor) ExtractContent(data []byte, _ StringKeyMap) StringKeyMap {
	info := jsonDecode(data)
	if info != nil {
		info = compressor.shortener.ExtractContent(info)
	}
	return info
}

// Override
func (compressor *MessageCompressor) CompressSymmetricKey(key StringKeyMap) []byte {
	key = compressor.shortener.CompressSymmetricKey(key)
	return jsonEncode(key)
}

// Override
func (compressor *MessageCompressor) ExtractSymmetricKey(data []byte) StringKeyMap {
	info := jsonDecode(data)
	if info != nil {
		info = compressor.shortener.ExtractSymmetricKey(info)
	}
	return info
}

// Override
func (compressor *MessageCompressor) CompressReliableMessage(msg StringKeyMap) []byte {
	msg = compressor.shortener.CompressReliableMessage(msg)
	return jsonEncode(msg)
}

// Override
func (compressor *MessageCompressor) ExtractReliableMessage(data []byte) StringKeyMap {
	info := jsonDecode(data)
	if info != nil {
		info = compressor.shortener.ExtractReliableMessage(info)
	}
	return info
}

func jsonEncode(dict StringKeyMap) []byte {
	text := JSONEncodeMap(dict)
	return UTF8Encode(text)
}

func jsonDecode(data []byte) StringKeyMap {
	text := UTF8Decode(data)
	if text == "" {
		return nil
	}
	return JSONDecodeMap(text)
}
