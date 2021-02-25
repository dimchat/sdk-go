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
type DefaultMetaShadow struct {
	BaseMetaShadow

	// caches
	_addresses map[uint8]Address
}

func (shadow *DefaultMetaShadow) Init(meta IMetaExt) *DefaultMetaShadow {
	if shadow.BaseMetaShadow.Init(meta) != nil {
		shadow._addresses = make(map[uint8]Address)
	}
	return shadow
}

func (shadow *DefaultMetaShadow) GenerateAddress(network uint8) Address {
	meta := shadow.Meta()
	// check caches
	address := shadow._addresses[network]
	if address == nil && meta.IsValid() {
		// generate and cache it
		address = BTCAddressGenerate(meta.Fingerprint(), network)
		shadow._addresses[network] = address
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
type BTCMetaShadow struct {
	BaseMetaShadow

	// cached
	_address Address
}

func (shadow *BTCMetaShadow) Init(meta IMetaExt) *BTCMetaShadow {
	if shadow.BaseMetaShadow.Init(meta) != nil {
		shadow._address = nil
	}
	return shadow
}

func (shadow *BTCMetaShadow) GenerateAddress(network uint8) Address {
	if network != BTCMain {
		return nil
	}
	meta := shadow.Meta()
	// check caches
	address := shadow._address
	if address == nil && meta.IsValid() {
		// generate and cache it
		key := meta.Key()
		pKey, ok := key.(CryptographyKey)
		if ok {
			address = BTCAddressGenerate(pKey.Data(), network)
			shadow._address = address
		}
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
type ETHMetaShadow struct {
	BaseMetaShadow

	// cached
	_address Address
}

func (shadow *ETHMetaShadow) Init(meta IMetaExt) *ETHMetaShadow {
	if shadow.BaseMetaShadow.Init(meta) != nil {
		shadow._address = nil
	}
	return shadow
}

func (shadow *ETHMetaShadow) GenerateAddress(network uint8) Address {
	if network != MAIN {
		return nil
	}
	meta := shadow.Meta()
	// check caches
	address := shadow._address
	if address == nil && meta.IsValid() {
		// generate and cache it
		key := meta.Key()
		pKey, ok := key.(CryptographyKey)
		if ok {
			address = ETHAddressGenerate(pKey.Data())
			shadow._address = address
		}
	}
	return address
}

//
//  Factory methods for Meta
//

func NewDefaultMeta(key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(BaseMeta).InitWithType(MKM, key, seed, fingerprint)
	meta.SetShadow(new(DefaultMetaShadow).Init(meta))
	return meta
}

func ParseDefaultMeta(dict map[string]interface{}) Meta {
	meta := new(BaseMeta).Init(dict)
	meta.SetShadow(new(DefaultMetaShadow).Init(meta))
	return meta
}

func NewBTCMeta(version uint8, key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(BaseMeta).InitWithType(version, key, seed, fingerprint)
	meta.SetShadow(new(BTCMetaShadow).Init(meta))
	return meta
}

func ParseBTCMeta(dict map[string]interface{}) Meta {
	meta := new(BaseMeta).Init(dict)
	meta.SetShadow(new(BTCMetaShadow).Init(meta))
	return meta
}

func NewETHMeta(version uint8, key VerifyKey, seed string, fingerprint []byte) Meta {
	meta := new(BaseMeta).InitWithType(version, key, seed, fingerprint)
	meta.SetShadow(new(ETHMetaShadow).Init(meta))
	return meta
}

func ParseETHMeta(dict map[string]interface{}) Meta {
	meta := new(BaseMeta).Init(dict)
	meta.SetShadow(new(ETHMetaShadow).Init(meta))
	return meta
}
