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
package crypto

import (
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

type EncryptedBundleHelper interface {

	/**
	 *  Encode key data
	 *
	 * @param bundle - encrypted key data with targets (ID terminals)
	 * @param did    - user ID
	 * @return encoded key data with targets (ID + terminals)
	 */
	EncodeBundle(bundle EncryptedBundle, did ID) StringKeyMap

	/**
	 *  Decode key data from 'message.keys'
	 *
	 * @param encodedKeys - encoded key data with targets (ID + terminals)
	 * @param did         - receiver ID
	 * @param terminals   - visa terminals
	 * @return encrypted key data with targets (ID terminals)
	 */
	DecodeBundle(encodedKeys StringKeyMap, did ID, terminals []string) EncryptedBundle
}

type DefaultBundleHelper struct {
	//EncryptedBundleHelper
}

// Override
func (helper DefaultBundleHelper) EncodeBundle(bundle EncryptedBundle, did ID) StringKeyMap {
	identifier := IDConcat(did.Name(), did.Address(), "")
	encodedKeys := NewMap()
	var base64 string
	dict := bundle.Map()
	for target, data := range dict {
		// encode data
		base64 = Base64Encode(data)
		if target == "" || target == "*" {
			target = identifier
		} else {
			target = identifier + "/" + target
		}
		// insert to 'message.keys' with ID + terminal
		encodedKeys[target] = base64
	}
	// OK
	return encodedKeys
}

// Override
func (helper DefaultBundleHelper) DecodeBundle(encodedKeys StringKeyMap, did ID, terminals []string) EncryptedBundle {
	bundle := NewEncryptedBundle()
	//
	//  0. ID string without terminal
	//
	identifier := IDConcat(did.Name(), did.Address(), "")
	var target string
	var base64 interface{}
	var ted TransportableData
	var data []byte
	for _, item := range terminals {
		if item == "" {
			target = "*"
		} else {
			target = item
		}
		//
		//  1. get encoded data with target (ID + terminal)
		//
		if target == "*" {
			base64 = encodedKeys[identifier]
		} else {
			base64 = encodedKeys[identifier+"/"+target]
		}
		if base64 == nil {
			// key data not found
			continue
		}
		//
		//  2. decode data
		//
		ted = ParseTransportableData(base64)
		if ted == nil {
			//panic("invalid transport data")
			continue
		}
		data = ted.Bytes()
		if data == nil || len(data) == 0 {
			//panic("data is nil")
			continue
		}
		//
		//  3. put data for target (ID terminal)
		//
		bundle.Set(target, data)
	}
	// OK
	return bundle
}

var sharedEncryptedBundleHelper EncryptedBundleHelper = &DefaultBundleHelper{}

func SetEncryptedBundleHelper(helper EncryptedBundleHelper) {
	sharedEncryptedBundleHelper = helper
}

func GetEncryptedBundleHelper() EncryptedBundleHelper {
	return sharedEncryptedBundleHelper
}
