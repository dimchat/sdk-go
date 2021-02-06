/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
package mkm

import (
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Address like BitCoin
 *
 *      data format: "network+digest+code"
 *          network    --  1 byte
 *          digest     -- 20 bytes
 *          code       --  4 bytes
 *
 *      algorithm:
 *          fingerprint = sign(seed, SK);
 *          digest      = ripemd160(sha256(fingerprint));
 *          code        = sha256(sha256(network + digest)).prefix(4);
 *          address     = base58_encode(network + digest + code);
 */
type BTCAddress struct {
	ConstantString
	Address

	_network uint8
}

func NewBTCAddress(address string, network uint8) *BTCAddress {
	return new(BTCAddress).Init(address, network)
}

func (address *BTCAddress) Init(string string, network uint8) *BTCAddress {
	if address.ConstantString.Init(string) != nil {
		address._network = network
	}
	return address
}

func (address *BTCAddress) String() string {
	return address.ConstantString.String()
}

func (address *BTCAddress) Equal(other interface{}) bool {
	return address.ConstantString.Equal(other)
}

func (address *BTCAddress) Network() uint8 {
	return address._network
}

func (address *BTCAddress) IsUser() bool {
	return NetworkTypeIsUser(address._network)
}

func (address *BTCAddress) IsGroup() bool {
	return NetworkTypeIsGroup(address._network)
}

func (address *BTCAddress) IsBroadcast() bool {
	return false
}

/**
 *  Generate address with fingerprint and network ID
 *
 * @param fingerprint = meta.fingerprint or key.data
 * @param network - address type
 * @return Address object
 */
func BTCAddressGenerate(fingerprint []byte, network uint8) *BTCAddress {
	// 1. digest = ripemd160(sha256(fingerprint))
	digest := RIPEMD160(SHA256(fingerprint))
	// 2. head = network + digest
	head := make([]byte, 21)
	head[0] = network
	ArrayCopy(digest, 0, head, 1, 20)
	// 3. cc = sha256(sha256(head)).prefix(4)
	cc := checkCode(head)
	// 4. data = base58_encode(head + cc)
	data := make([]byte, 25)
	ArrayCopy(head, 0, data, 0, 21)
	ArrayCopy(cc, 0, data, 21, 24)
	base58 := Base58Encode(data)
	return NewBTCAddress(base58, network)
}

/**
 *  Parse a string for BTC address
 *
 * @param base58 - address string
 * @return null on error
 */
func BTCAddressParse(base58 string) *BTCAddress {
	// decode
	data := Base58Decode(base58)
	if len(data) != 25 {
		//panic("address length error")
		return nil
	}
	// CheckCode
	prefix := make([]byte, 21)
	suffix := make([]byte, 4)
	ArrayCopy(data, 0, prefix, 0, 21)
	ArrayCopy(data, 21, suffix, 0, 4)
	cc := checkCode(prefix)
	if !ArraysEqual(cc, suffix) {
		panic("address check code error")
		return nil
	}
	network := data[0]
	return NewBTCAddress(base58, network)
}

func checkCode(data []byte) []byte {
	sha256d := SHA256(SHA256(data))
	cc := make([]byte, 4)
	ArrayCopy(sha256d, 0, cc, 0, 4)
	return cc
}
