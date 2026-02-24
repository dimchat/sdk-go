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

import . "github.com/dimchat/mkm-go/protocol"

// Group defines the interface for group entities (collections of users)
//
// # Extends Entity with group-specific functionality for membership and ownership management
//
// Core Group Roles (hierarchy):
//   - founder: Original creator of the group (immutable once set)
//   - owner: Current administrator of the group (may change via consensus)
//   - members: Regular users in the group (owner must be a member)
//   - administrators: Optional privileged members (extended role)
type Group interface {
	Entity

	// Founder returns the unique ID of the group's original creator
	//
	// This value is immutable - once set, it never changes
	Founder() ID

	// Owner returns the ID of the current group owner
	//
	// The owner may be different from the founder (ownership transfer)
	Owner() ID

	// Members returns the list of all user IDs in the group
	//
	// Important: The group owner MUST be included in the members list (usually first)
	Members() []ID
}

// BaseGroup is the base implementation of the Group interface
//
// Extends BaseEntity with group-specific fields (immutable founder ID)
type BaseGroup struct {
	//Group
	*BaseEntity

	// founder stores the immutable ID of the group's original creator
	//
	// Once set during group creation, this value can never be changed
	founder ID
}

func NewBaseGroup(gid ID) *BaseGroup {
	return &BaseGroup{
		BaseEntity: NewBaseEntity(gid),
		// lazy load
		founder: nil,
	}
}

// Override
func (group *BaseGroup) Founder() ID {
	user := group.founder
	if user == nil {
		facebook := group.DataSource()
		if facebook == nil {
			//panic("group datasource not set yet")
			return nil
		}
		user = facebook.GetFounder(group.ID())
		group.founder = user
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
