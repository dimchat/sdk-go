/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
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
package dimp

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

var KEYWORDS = [...]string {
	"all", "everyone", "anyone", "owner", "founder",
	// --------------------------------
	"dkd", "mkm", "dimp", "dim", "dimt",
	"rsa", "ecc", "aes", "des", "btc", "eth",
	// --------------------------------
	"crypto", "key", "symmetric", "asymmetric",
	"public", "private", "secret", "password",
	"id", "address", "meta", "profile", "document",
	"entity", "user", "group", "contact",
	// --------------------------------
	"member", "admin", "administrator", "assistant",
	"main", "polylogue", "chatroom",
	"social", "organization",
	"company", "school", "government", "department",
	"provider", "station", "thing", "robot",
	// --------------------------------
	"message", "instant", "secure", "reliable",
	"envelope", "sender", "receiver", "time",
	"content", "forward", "command", "history",
	"keys", "data", "signature",
	// --------------------------------
	"type", "serial", "sn",
	"text", "file", "image", "audio", "video", "page",
	"handshake", "receipt", "block", "mute",
	"register", "suicide", "found", "abdicate",
	"invite", "expel", "join", "quit", "reset", "query",
	"hire", "fire", "resign",
	// --------------------------------
	"server", "client", "terminal", "local", "remote",
	"barrack", "cache", "transceiver",
	"ans", "facebook", "store", "messenger",
	"root", "supervisor",
}

type AddressNameService interface {

	IsReserved(name string) bool

	Cache(name string, identifier ID) bool

	/**
	 *  Get ID by short name
	 *
	 * @param name - sort name
	 * @return user ID
	 */
	GetID(name string) ID

	/**
	 *  Get all short names with the same ID
	 *
	 * @param identifier - user ID
	 * @return short name list
	 */
	GetNames(identifier ID) []string

	/**
	 *  Save ANS record
	 *
	 * @param name - username
	 * @param identifier - user ID; if empty, means delete this name
	 * @return true on success
	 */
	Save(name string, identifier ID) bool
}

/**
 *  ANS
 */
type AddressNameServer struct {

	_reserved map[string]bool
	_caches map[string]ID
}

func (ans *AddressNameServer) Init() *AddressNameServer {
	// reserved names
	ans._reserved = make(map[string]bool, len(KEYWORDS))
	for _, item := range KEYWORDS {
		ans._reserved[item] = true
	}

	ans._caches = make(map[string]ID, 1024)
	// constant ANS records
	ans.setID("all", EVERYONE)
	ans.setID("everyone", EVERYONE)
	ans.setID("anyone", ANYONE)
	ans.setID("owner", ANYONE)
	ans.setID("founder", FOUNDER)

	return ans
}

func (ans *AddressNameServer) setID(name string, identifier ID) {
	if ValueIsNil(identifier) {
		delete(ans._caches, name)
	} else {
		ans._caches[name] = identifier
	}
}

//-------- IAddressNameService

func (ans *AddressNameServer) IsReserved(name string) bool {
	return ans._reserved[name]
}

func (ans *AddressNameServer) Cache(name string, identifier ID) bool {
	if ans.IsReserved(name) {
		// this name is reserved, cannot register
		return false
	}
	ans.setID(name, identifier)
	return true
}

func (ans *AddressNameServer) GetID(name string) ID {
	return ans._caches[name]
}

func (ans *AddressNameServer) GetNames(identifier ID) []string {
	array := make([]string, 0, 1)
	for key, value := range ans._caches {
		if identifier.Equal(value) {
			array = append(array, key)
		}
	}
	return array
}

func (ans *AddressNameServer) Save(name string, identifier ID) bool {
	// override to save this record into local storage
	return ans.Cache(name, identifier)
}
