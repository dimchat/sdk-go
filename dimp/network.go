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
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Robot User
 */
type Robot struct {
	User
}

func NewRobot(identifier ID) *Robot {
	return new(Robot).Init(identifier)
}

func (user *Robot) Init(identifier ID) *Robot {
	if user.User.Init(identifier) != nil {
	}
	return user
}

/**
 *  DIM Server
 */
type Station struct {
	User

	_host string
	_port uint16
}

func NewStation(identifier ID, host string, port uint16) *Station {
	return new(Station).Init(identifier, host, port)
}

func (server *Station) Init(identifier ID, host string, port uint16) *Station {
	if server.User.Init(identifier) != nil {
		server._host = host
		server._port = port
	}
	return server
}

func (server *Station) Host() string {
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

func (server *Station) Port() uint16 {
	if server._port == 0 {
		doc := server.GetDocument("*")
		if doc != nil {
			port := doc.GetProperty("port")
			if port != nil {
				server._port = port.(uint16)
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
	Group
}

func NewServiceProvider(identifier ID) *ServiceProvider {
	return new(ServiceProvider).Init(identifier)
}

func (sp *ServiceProvider) Init(identifier ID) *ServiceProvider {
	if sp.Group.Init(identifier) != nil {
	}
	return sp
}

func (sp ServiceProvider) GetStations() []ID {
	return sp.GetMembers()
}
