/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Group for organizing users
 *
 * <pre>
 *  roles:
 *      founder
 *      owner
 *      members
 *      administrators - Optional
 *  </pre>
 */
type Group interface {
	Entity

	Founder() ID

	Owner() ID

	// NOTICE: the owner must be a member
	//         (usually the first one)
	Members() []ID
}

func NewGroup(gid ID) Group {
	group := &BaseGroup{}
	return group.Init(gid)
}

/**
 *  Base Group
 *  ~~~~~~~~~~
 */
type BaseGroup struct {
	BaseEntity

	// once the group founder is set, it will never change
	_founder ID
}

func (group *BaseGroup) Init(gid ID) Group {
	if group.BaseEntity.Init(gid) != nil {
		// lazy load
		group._founder = nil
	}
	return group
}

// Override
func (group *BaseGroup) Founder() ID {
	user := group._founder
	if user == nil {
		facebook := group.DataSource()
		if facebook == nil {
			//panic("group datasource not set yet")
			return nil
		}
		user = facebook.GetFounder(group.ID())
		group._founder = user
	}
	return user
}

// Override
func (group *BaseGroup) Owner() ID {
	facebook := group.DataSource()
	if facebook == nil {
		//panic("group datasource not set yet")
		return nil
	}
	return facebook.GetOwner(group.ID())
}

// Override
func (group *BaseGroup) Members() []ID {
	facebook := group.DataSource()
	if facebook == nil {
		//panic("group datasource not set yet")
		return nil
	}
	return facebook.GetMembers(group.ID())
}
