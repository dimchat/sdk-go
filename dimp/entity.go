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
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Robot User
 */
type Robot struct {
	IRobot
	BaseUser
}
type IRobot interface {

	Master() ID
}

func (user *Robot) Init(identifier ID) *Robot {
	if user.BaseUser.Init(identifier) != nil {
	}
	return user
}

//-------- IRobot

func (user *Robot) Master() ID {
	doc := user.Visa()
	if doc == nil {
		return nil
	}
	return IDParse(doc.Get("master"))
}

/**
 *  DIM Server
 */
type Station struct {
	IStation
	BaseUser

	_host string
	_port uint16
}
type IStation interface {

	Host() string
	Port() uint16
}

func (server *Station) Init(identifier ID, host string, port uint16) *Station {
	if server.BaseUser.Init(identifier) != nil {
		server._host = host
		server._port = port
	}
	return server
}

//-------- IStation

func (server *Station) Host() string {
	if server._host == "" {
		doc := server.GetDocument("*")
		if doc != nil {
			host, ok := doc.GetProperty("host").(string)
			if ok {
				server._host = host
			}
		}
		if server._host == "" {
			server._host = "0.0.0.0"
		}
	}
	return server._host
}

func (server *Station) Port() uint16 {
	if server._port == 0 {
		doc := server.GetDocument("*")
		if doc != nil {
			port, ok := doc.GetProperty("port").(uint16)
			if ok {
				server._port = port
			}
		}
		if server._port == 0 {
			server._port = 9394
		}
	}
	return server._port
}

/**
 *  DIM Station Owner
 */
type ServiceProvider struct {
	IServiceProvider
	BaseGroup
}
type IServiceProvider interface {

	GetStations() []ID
}

func (sp *ServiceProvider) Init(identifier ID) *ServiceProvider {
	if sp.BaseGroup.Init(identifier) != nil {
	}
	return sp
}

//-------- IServiceProvider

func (sp *ServiceProvider) GetStations() []ID {
	return sp.Members()
}

/**
 *  Simple group chat
 */
type Polylogue struct {
	IPolylogue
	BaseGroup
}
type IPolylogue interface {

	Owner() ID
}

func (group *Polylogue) Init(identifier ID) *Polylogue {
	if group.BaseGroup.Init(identifier) != nil {
	}
	return group
}

//-------- IPolylogue

func (group *Polylogue) Owner() ID {
	owner := group.BaseGroup.Owner()
	if owner == nil {
		// polylogue owner is its founder
		owner = group.Founder()
	}
	return owner
}

/**
 *  Big group with admins
 */
type Chatroom struct {
	IChatroom
	BaseGroup
}
type IChatroom interface {

	Admins() []ID
}

func (group *Chatroom) Init(identifier ID) *Chatroom {
	if group.BaseGroup.Init(identifier) != nil {
	}
	return group
}

//-------- IChatroom

func (group *Chatroom) Admins() []ID {
	delegate := group.DataSource().(ChatroomDataSource)
	return delegate.GetAdmins(group.ID())
}

/**
 *  This interface is for getting information for chatroom
 *  Chatroom admins should be set complying with the consensus algorithm
 */
type ChatroomDataSource interface {
	IChatroomDataSource
	GroupDataSource
}
type IChatroomDataSource interface {

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
func NewUser(identifier ID) *BaseUser {
	return new(BaseUser).Init(identifier)
}

func NewGroup(identifier ID) *BaseGroup {
	return new(BaseGroup).Init(identifier)
}

func NewRobot(identifier ID) *Robot {
	return new(Robot).Init(identifier)
}

func NewStation(identifier ID, host string, port uint16) *Station {
	return new(Station).Init(identifier, host, port)
}

func NewServiceProvider(identifier ID) *ServiceProvider {
	return new(ServiceProvider).Init(identifier)
}

func NewPolylogue(identifier ID) *Polylogue {
	return new(Polylogue).Init(identifier)
}

func NewChatroom(identifier ID) *Chatroom {
	return new(Chatroom).Init(identifier)
}
