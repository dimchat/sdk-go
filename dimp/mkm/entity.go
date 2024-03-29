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
package mkm

import (
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
)

type Robot interface {
	User

	Master() ID
}

type Station interface {
	User

	Host() string
	Port() uint16
}

type Provider interface {
	Group

	GetStations() []ID
}

type Polylogue interface {
	Group
}

type Chatroom interface {
	Group

	Admins() []ID
}

type ChatroomDataSource interface {
	/**
	 *  This interface is for getting information for chatroom
	 *  Chatroom admins should be set complying with the consensus algorithm
	 */
	GroupDataSource

	/**
	 *  Get all admins in the chatroom
	 *
	 * @param chatroom - chatroom ID
	 * @return admin ID list
	 */
	GetAdmins(group ID) []ID
}

//
//  Creators
//
func NewUser(identifier ID) User {
	user := new(BaseUser)
	user.Init(identifier)
	return user
}

func NewGroup(identifier ID) Group {
	group := new(BaseGroup)
	group.Init(identifier)
	return group
}

func NewRobot(identifier ID) Robot {
	robot := new(BaseRobot)
	robot.Init(identifier)
	return robot
}

func NewStation(identifier ID, host string, port uint16) Station {
	station := new(BaseStation)
	station.Init(identifier, host, port)
	return station
}

func NewProvider(identifier ID) Provider {
	sp := new(ServiceProvider)
	sp.Init(identifier)
	return sp
}

func NewPolylogue(identifier ID) Polylogue {
	group := new(SimpleGroup)
	group.Init(identifier)
	return group
}

func NewChatroom(identifier ID) Chatroom {
	group := new(BigGroup)
	group.Init(identifier)
	return group
}
