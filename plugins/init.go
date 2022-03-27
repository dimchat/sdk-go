/* license: https://mit-license.org
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
package plugins

import (
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/plugins/crypto"
	. "github.com/dimchat/sdk-go/plugins/digest"
	. "github.com/dimchat/sdk-go/plugins/mkm"
)

/**
 *  Data Digesters
 */
func RegisterDataDigesters()  {
	RIPEMD160SetDigester(new(RIPEMD160Digester))
	KECCAK256SetDigester(new(KECCAK256Digester))
}

/**
 *  Address factory
 */
func BuildAddressFactory() AddressFactory {
	factory := AddressGetFactory()
	if factory == nil {
		factory = new(GeneralAddressFactory).Init(func(address string) Address {
			if len(address) == 42 {
				return ETHAddressParse(address)
			} else {
				return BTCAddressParse(address)
			}
		})
		AddressSetFactory(factory)
	}
	return factory
}

/**
 *  Meta factories
 */
func BuildMetaFactories() {
	MetaSetFactory(MKM, NewGeneralMetaFactory(MKM))
	MetaSetFactory(BTC, NewGeneralMetaFactory(BTC))
	MetaSetFactory(ExBTC, NewGeneralMetaFactory(ExBTC))
	MetaSetFactory(ETH, NewGeneralMetaFactory(ETH))
	MetaSetFactory(ExETH, NewGeneralMetaFactory(ExETH))
}

/**
 *  Document factories
 */
func BuildDocumentFactories() {
	DocumentSetFactory("*", NewGeneralDocumentFactory("*"))
	DocumentSetFactory(VISA, NewGeneralDocumentFactory(VISA))
	DocumentSetFactory(PROFILE, NewGeneralDocumentFactory(PROFILE))
	DocumentSetFactory(BULLETIN, NewGeneralDocumentFactory(BULLETIN))
}

func init() {
	RegisterDataDigesters()

	BuildSymmetricKeyFactories()
	BuildPrivateKeyFactories()
	BuildPublicKeyFactories()

	BuildAddressFactory()

	BuildMetaFactories()
	BuildDocumentFactories()
}
