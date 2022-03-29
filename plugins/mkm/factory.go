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
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Meta Factory
 *  ~~~~~~~~~~~~
 */
type GeneralMetaFactory struct {
	MetaFactory

	_type uint8
}

func NewGeneralMetaFactory(version uint8) *GeneralMetaFactory {
	return new(GeneralMetaFactory).Init(version)
}

func (factory *GeneralMetaFactory) Init(version uint8) *GeneralMetaFactory {
	factory._type = version
	return factory
}

//-------- IMetaFactory

func (factory *GeneralMetaFactory) CreateMeta(key VerifyKey, seed string, fingerprint []byte) Meta {
	switch factory._type {
	case MKM:
		return NewDefaultMeta(key, seed, fingerprint)
	case BTC:
		return NewBTCMeta(BTC, key, seed, fingerprint)
	case ExBTC:
		return NewBTCMeta(ExBTC, key, seed, fingerprint)
	case ETH:
		return NewETHMeta(ETH, key, seed, fingerprint)
	case ExETH:
		return NewETHMeta(ExETH, key, seed, fingerprint)
	default:
		return nil
	}
}

func (factory *GeneralMetaFactory) GenerateMeta(sKey SignKey, seed string) Meta {
	priKey, ok := sKey.(PrivateKey)
	if !ok {
		return nil
	}
	pubKey := priKey.PublicKey()
	if pubKey == nil {
		return nil
	}

	var fingerprint []byte
	if len(seed) == 0 {
		fingerprint = nil
	} else {
		fingerprint = sKey.Sign(UTF8Encode(seed))
	}
	return factory.CreateMeta(pubKey, seed, fingerprint)
}

func (factory *GeneralMetaFactory) ParseMeta(meta map[string]interface{}) Meta {
	version := MetaGetType(meta)
	switch version {
	case MKM:
		return ParseDefaultMeta(meta)
	case BTC:
		return ParseBTCMeta(meta)
	case ExBTC:
		return ParseBTCMeta(meta)
	case ETH:
		return ParseETHMeta(meta)
	case ExETH:
		return ParseETHMeta(meta)
	default:
		return nil
	}
}

/**
 *  Document Factory
 *  ~~~~~~~~~~~~~~~~
 */
type GeneralDocumentFactory struct {
	IDocumentFactory

	_type string
}

func NewGeneralDocumentFactory(docType string) *GeneralDocumentFactory {
	return new(GeneralDocumentFactory).Init(docType)
}

func (factory *GeneralDocumentFactory) Init(docType string) *GeneralDocumentFactory {
	factory._type = docType
	return factory
}

func (factory *GeneralDocumentFactory) getDocType(identifier ID) string {
	if factory._type == "*" {
		if identifier.IsGroup() {
			return BULLETIN
		}
		if identifier.IsUser() {
			return VISA
		}
		return PROFILE
	}
	return factory._type
}

//-------- IDocumentFactory

func (factory *GeneralDocumentFactory) CreateDocument(identifier ID, data []byte, signature []byte) Document {
	docType := factory.getDocType(identifier)
	if docType == VISA {
		return NewVisa(identifier, data, signature)
	} else if docType == BULLETIN {
		return NewBulletin(identifier, data, signature)
	} else {
		return NewDocument(identifier, data, signature)
	}
}

func (factory *GeneralDocumentFactory) ParseDocument(doc map[string]interface{}) Document {
	identifier := DocumentGetID(doc)
	if identifier == nil {
		return nil
	}
	docType := DocumentGetType(doc)
	if docType == "" {
		if identifier.IsGroup() {
			docType = BULLETIN
		} else if identifier.IsUser() {
			docType = VISA
		} else {
			docType = PROFILE
		}
	}
	if docType == VISA {
		return ParseVisa(doc)
	}
	if docType == BULLETIN {
		return ParseBulletin(doc)
	}
	return ParseDocument(doc)
}

//
//  Factory methods for Document
//

func NewDocument(identifier ID, data []byte, signature []byte) Document {
	return new(BaseDocument).InitWithData(identifier, data, signature)
}

func ParseDocument(dict map[string]interface{}) Document {
	return new(BaseDocument).Init(dict)
}

func NewVisa(identifier ID, data []byte, signature []byte) Visa {
	return new(BaseVisa).InitWithData(identifier, data, signature)
}

func ParseVisa(dict map[string]interface{}) Visa {
	return new(BaseVisa).Init(dict)
}

func NewBulletin(identifier ID, data []byte, signature []byte) Bulletin {
	return new(BaseBulletin).InitWithData(identifier, data, signature)
}

func ParseBulletin(dict map[string]interface{}) Bulletin {
	return new(BaseBulletin).Init(dict)
}
