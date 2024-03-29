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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Default Meta to build ID with 'name@address'
 *
 *  version:
 *      0x01 - MKM
 *
 *  algorithm:
 *      CT      = fingerprint = sKey.sign(seed);
 *      hash    = ripemd160(sha256(CT));
 *      code    = sha256(sha256(network + hash)).prefix(4);
 *      address = base58_encode(network + hash + code);
 */
type DefaultMeta struct {
	BaseMeta

	// caches
	_addresses map[NetworkType]Address
}

func (meta *DefaultMeta) Init(dict map[string]interface{}) Meta {
	if meta.BaseMeta.Init(dict) != nil {
		meta._addresses = make(map[NetworkType]Address)
	}
	return meta
}

func (meta *DefaultMeta) InitWithType(version MetaType, key VerifyKey, seed string, fingerprint []byte) Meta {
	if meta.BaseMeta.InitWithType(version, key, seed, fingerprint) != nil {
		meta._addresses = make(map[NetworkType]Address)
	}
	return meta
}

//-------- IMeta

func (meta *DefaultMeta) GenerateAddress(network NetworkType) Address {
	// check caches
	address := meta._addresses[network]
	if address == nil {
		// generate and cache it
		address = BTCAddressGenerate(meta.Fingerprint(), network)
		meta._addresses[network] = address
	}
	return address
}

/**
 *  Meta to build BTC address for ID
 *
 *  version:
 *      0x02 - BTC
 *      0x03 - ExBTC
 *
 *  algorithm:
 *      CT      = key.data;
 *      hash    = ripemd160(sha256(CT));
 *      code    = sha256(sha256(network + hash)).prefix(4);
 *      address = base58_encode(network + hash + code);
 */
type BTCMeta struct {
	BaseMeta

	// cached
	_address Address
}

func (meta *BTCMeta) Init(dict map[string]interface{}) Meta {
	if meta.BaseMeta.Init(dict) != nil {
		meta._address = nil
	}
	return meta
}

func (meta *BTCMeta) InitWithType(version MetaType, key VerifyKey, seed string, fingerprint []byte) Meta {
	if meta.BaseMeta.InitWithType(version, key, seed, fingerprint) != nil {
		meta._address = nil
	}
	return meta
}

//-------- IMeta

func (meta *BTCMeta) GenerateAddress(network NetworkType) Address {
	if network != BTCMain {
		return nil
	}
	// check caches
	address := meta._address
	if address == nil {
		// generate and cache it
		key := meta.Key()
		address = BTCAddressGenerate(key.Data(), network)
		meta._address = address
	}
	return address
}

/**
 *  Meta to build ETH address for ID
 *
 *  version:
 *      0x04 - ETH
 *      0x05 - ExETH
 *
 *  algorithm:
 *      CT      = key.data;  // without prefix byte
 *      digest  = keccak256(CT);
 *      address = hex_encode(digest.suffix(20));
 */
type ETHMeta struct {
	BaseMeta

	// cached
	_address Address
}

func (meta *ETHMeta) Init(dict map[string]interface{}) Meta {
	if meta.BaseMeta.Init(dict) != nil {
		meta._address = nil
	}
	return meta
}

func (meta *ETHMeta) InitWithType(version MetaType, key VerifyKey, seed string, fingerprint []byte) Meta {
	if meta.BaseMeta.InitWithType(version, key, seed, fingerprint) != nil {
		meta._address = nil
	}
	return meta
}

//-------- IMeta

func (meta *ETHMeta) GenerateAddress(network NetworkType) Address {
	if network != MAIN {
		return nil
	}
	// check caches
	address := meta._address
	if address == nil {
		// generate and cache it
		key := meta.Key()
		address = ETHAddressGenerate(key.Data())
		meta._address = address
	}
	return address
}

//
//  Factory methods for Meta
//

func NewDefaultMeta(key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(DefaultMeta)
	meta.InitWithType(MKM, key, seed, fingerprint)
	return meta
}

func ParseDefaultMeta(dict map[string]interface{}) Meta {
	meta := new(DefaultMeta)
	meta.Init(dict)
	return meta
}

func NewBTCMeta(version MetaType, key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(BTCMeta)
	meta.InitWithType(version, key, seed, fingerprint)
	return meta
}

func ParseBTCMeta(dict map[string]interface{}) Meta {
	meta := new(BTCMeta)
	meta.Init(dict)
	return meta
}

func NewETHMeta(version MetaType, key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(ETHMeta)
	meta.InitWithType(version, key, seed, fingerprint)
	return meta
}

func ParseETHMeta(dict map[string]interface{}) Meta {
	meta := new(ETHMeta)
	meta.Init(dict)
	return meta
}
