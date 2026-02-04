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

/** Short Keys
<pre>
    ======+==================================================+==================
          |   Message        Content        Symmetric Key    |    Description
    ------+--------------------------------------------------+------------------
    "A"   |                                 "algorithm"      |
    "C"   |   "content"      "command"                       |
    "D"   |   "data"                        "data"           |
    "F"   |   "sender"                                       |   (From)
    "G"   |   "group"        "group"                         |
    "I"   |                                 "iv"             |
    "K"   |   "key", "keys"                                  |
    "M"   |   "meta"                                         |
    "N"   |                  "sn"                            |   (Number)
    "P"   |   "visa"                                         |   (Profile)
    "R"   |   "receiver"                                     |
    "S"   |   ...                                            |
    "T"   |   "type"         "type"                          |
    "V"   |   "signature"                                    |   (Verification)
    "W"   |   "time"         "time"                          |   (When)
    ======+==================================================+==================

    Note:
        "S" - deprecated (ambiguous for "sender" and "signature")
</pre>
*/

var ShortKeys = struct {
	ContentKeys []string
	CryptoKeys  []string
	MessageKeys []string
}{
	ContentKeys: []string{
		"T", "type",
		"N", "sn",
		"W", "time", // When
		"G", "group",
		"C", "command", // Command name
	},
	CryptoKeys: []string{
		"A", "algorithm",
		"D", "data",
		"I", "iv", // Initial Vector
	},
	MessageKeys: []string{
		"F", "sender", // From
		"R", "receiver", // Rcpt to
		"W", "time", // When
		"T", "type",
		"G", "group",
		//------------------
		"K", "key", // or "keys"
		"D", "data",
		"V", "signature", // Verification
		//------------------
		"M", "meta",
		"P", "visa", // Profile
	},
}

//
//  Factory
//

func CreateCompressor() Compressor {
	factory := GetCompressFactory()
	return factory.CreateCompressor()
}

type CompressFactory interface {
	CreateCompressor() Compressor
}

var sharedCompressFactory CompressFactory = &defaultCompressFactory{}

func SetCompressFactory(factory CompressFactory) {
	sharedCompressFactory = factory
}

func GetCompressFactory() CompressFactory {
	return sharedCompressFactory
}

type defaultCompressFactory struct {
	//CompressFactory
}

// Override
func (factory defaultCompressFactory) CreateCompressor() Compressor {
	shortener := NewShortener(ShortKeys.ContentKeys, ShortKeys.CryptoKeys, ShortKeys.MessageKeys)
	return NewCompressor(shortener)
}

func NewShortener(contentKeys, cryptoKeys, messageKeys []string) Shortener {
	return &MessageShortener{
		contentShortKeys: contentKeys,
		cryptoShortKeys:  cryptoKeys,
		messageShortKeys: messageKeys,
	}
}

func NewCompressor(shortener Shortener) Compressor {
	return &MessageCompressor{
		shortener: shortener,
	}
}
