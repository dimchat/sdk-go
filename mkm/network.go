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

/**
 *  Robot User
 *  ~~~~~~~~~~
 */
type BaseRobot struct {
	BaseUser
}

//func (user *BaseRobot) Init(identifier ID) Robot {
//	if user.BaseUser.Init(identifier) != nil {
//	}
//	return user
//}

//-------- IRobot

func (user *BaseRobot) Master() ID {
	doc := user.Visa()
	if doc == nil {
		return nil
	}
	return IDParse(doc.Get("master"))
}

/**
 *  DIM Server
 *  ~~~~~~~~~~
 */
type BaseStation struct {
	BaseUser

	_host string
	_port uint16
}

func (server *BaseStation) Init(identifier ID, host string, port uint16) Station {
	if server.BaseUser.Init(identifier) != nil {
		server._host = host
		server._port = port
	}
	return server
}

//-------- IStation

func (server *BaseStation) Host() string {
	if server._host == "" {
		doc := server.GetDocument("*")
		if doc != nil {
			host := doc.GetProperty("host")
			if host != nil {
				server._host = host.(string)
			}
		}
		if server._host == "" {
			server._host = "0.0.0.0"
		}
	}
	return server._host
}

func (server *BaseStation) Port() uint16 {
	if server._port == 0 {
		doc := server.GetDocument("*")
		if doc != nil {
			port := doc.GetProperty("port")
			if port != nil {
				server._port = uint16(port.(float64))
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
 *  ~~~~~~~~~~~~~~~~~
 */
type ServiceProvider struct {
	BaseGroup
}

//func (sp *ServiceProvider) Init(identifier ID) Provider {
//	if sp.BaseGroup.Init(identifier) != nil {
//	}
//	return sp
//}

//-------- IServiceProvider

func (sp *ServiceProvider) GetStations() []ID {
	return sp.Members()
}
