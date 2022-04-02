/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
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
package mkm

import (
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	"strings"
)

/**
 *  Address like Ethereum
 *
 *      data format: "0x{address}"
 *
 *      algorithm:
 *          fingerprint = PK.data;
 *          digest      = keccak256(fingerprint);
 *          address     = hex_encode(digest.suffix(20));
 */
type ETHAddress struct {
	ConstantString
}

func NewETHAddress(address string) *ETHAddress {
	return new(ETHAddress).Init(address)
}

func (address *ETHAddress) Init(string string) *ETHAddress {
	if address.ConstantString.Init(string) != nil {
	}
	return address
}

//-------- IAddress

func (address *ETHAddress) Network() NetworkType {
	return MAIN
}

func (address *ETHAddress) IsUser() bool {
	return true
}

func (address *ETHAddress) IsGroup() bool {
	return false
}

func (address *ETHAddress) IsBroadcast() bool {
	return false
}

/**
 *  Generate ETH address with key.data
 *
 * @param fingerprint = key.data
 * @return Address object
 */
func ETHAddressGenerate(fingerprint []byte) *ETHAddress {
	if len(fingerprint) == 65 {
		fingerprint = fingerprint[1:]
	}
	// 1. digest = keccak256(fingerprint);
	digest := KECCAK256(fingerprint)
	// 2. address = hex_encode(digest.suffix(20));
	address := "0x" + eip55(HexEncode(digest[32-20:]))
	return NewETHAddress(address)
}

/**
 *  Parse a string for ETH address
 *
 * @param address - address string
 * @return null on error
 */
func ETHAddressParse(address string) *ETHAddress {
	if isETH(address) {
		return NewETHAddress(address)
	} else {
		return nil
	}
}

// https://eips.ethereum.org/EIPS/eip-55
func eip55(hex string) string {
	sb := make([]byte, 40)
	utf8 := UTF8Encode(hex)
	hash := KECCAK256(utf8)
	var ch byte
	var i uint8
	for i = 0; i < 40; i++ {
		ch = utf8[i]
		if ch > '9' {
			// check for each 4 bits in the hash table
			// if the first bit is '1',
			//     change the character to uppercase
			ch -= (hash[i >> 1] << (i << 2 & 4) & 0x80) >> 2
		}
		sb[i] = ch
	}
	return UTF8Decode(sb)
}

func isETH(address string) bool {
	if len(address) != 42 {
		return false
	}
	if address[0] != '0' || address[1] != 'x' {
		return false
	}
	var ch byte
	for i := 2; i < 42; i++ {
		ch = address[i]
		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		// unexpected character
		return false
	}
	return true
}

func GetValidateETHAddressString(address string) string {
	if isETH(address) {
		lower := strings.ToLower(address[2:])
		return "0x" + eip55(lower)
	}
	return ""
}

func IsValidateETHAddressString(address string) bool {
	validate := GetValidateETHAddressString(address)
	return validate == address
}
