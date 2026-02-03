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

	//DataSource() GroupDataSource
	//SetDataSource(facebook GroupDataSource)
}

/**
 *  Group Data Source
 *
 *  <pre>
 *  1. founder has the same public key with the group's meta.key
 *  2. owner and members should be set complying with consensus algorithm
 *  </pre>
 */
type GroupDataSource interface {
	EntityDataSource

	/**
	 *  Get founder of the group
	 *
	 * @param group - group ID
	 * @return fonder ID
	 */
	GetFounder(group ID) ID

	/**
	 *  Get current owner of the group
	 *
	 * @param group - group ID
	 * @return owner ID
	 */
	GetOwner(group ID) ID

	/**
	 *  Get all members in the group
	 *
	 * @param group - group ID
	 * @return members list (ID)
	 */
	GetMembers(group ID) []ID
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

//// Override
//func (group *BaseGroup) DataSource() GroupDataSource {
//	facebook := group.BaseEntity.DataSource()
//	if delegate, ok := facebook.(GroupDataSource); ok {
//		return delegate
//	}
//	//panic("not a GroupDataSource")
//	return nil
//}

// Override
func (group *BaseGroup) Founder() ID {
	user := group._founder
	if user == nil {
		ds := group.DataSource()
		if facebook, ok := ds.(GroupDataSource); ok {
			user = facebook.GetFounder(group.ID())
			group._founder = user
		} else {
			//panic("group datasource not set yet")
		}
	}
	return user
}

// Override
func (group *BaseGroup) Owner() ID {
	ds := group.DataSource()
	if facebook, ok := ds.(GroupDataSource); ok {
		return facebook.GetOwner(group.ID())
	}
	//panic("group datasource not set yet")
	return nil
}

// Override
func (group *BaseGroup) Members() []ID {
	ds := group.DataSource()
	if facebook, ok := ds.(GroupDataSource); ok {
		return facebook.GetMembers(group.ID())
	}
	//panic("group datasource not set yet")
	return nil
}
