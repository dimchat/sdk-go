package cpu

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
	. "github.com/dimchat/sdk-go/dimp/cpu"
)

/**
 *  Group command: INVITE
 */
type InviteCommandProcessor struct {
	GroupCommandProcessor
}

func (gpu *InviteCommandProcessor) callReset(cmd Command, rMsg ReliableMessage) Content {
	processor := CommandProcessorGetByName(RESET)
	processor.SetMessenger(gpu.Messenger())
	return processor.Execute(cmd, rMsg)
}

func (gpu *InviteCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) Content {
	//gCmd := cmd.(*GroupCommand)
	facebook := gpu.Facebook()

	// 0. check group
	group := cmd.Group()
	owner := facebook.GetOwner(group)
	members := facebook.GetMembers(group)
	if owner == nil || members == nil || len(members) == 0 {
		// NOTICE: group membership lost?
		//         reset group members
		return gpu.callReset(cmd, rMsg)
	}

	// 1. check permission
	sender := rMsg.Sender()
	if contains(sender, members) == false {
		// not a member? check assistants
		assistants := facebook.GetAssistants(group)
		if assistants == nil || contains(sender, assistants) == false {
			panic(sender.String() + " is not a member/assistant" +
				" of group " + group.String() + ", cannot invite members")
			return nil
		}
	}

	// 2. inviting members
	inviteList := gpu.GetMembers(cmd.(*GroupCommand))
	if inviteList == nil || len(inviteList) == 0 {
		panic("invite command error")
		return nil
	}
	// 2.1. check for reset
	if sender.Equal(owner) && contains(owner, inviteList) {
		// NOTICE: owner invites owner?
		//         it means this should be a 'reset' command
		return gpu.callReset(cmd, rMsg)
	}
	// 2.2. build invited-list
	addedList := make([]string, 0, len(inviteList))
	for _, item := range inviteList {
		if contains(item, members) {
			continue
		}
		// new member found
		addedList = append(addedList, item.String())
		members = append(members, item)
	}
	// 2.3. do invite
	if len(addedList) > 0 {
		if facebook.SaveMembers(members, group) {
			cmd.Set("added", addedList)
		}
	}

	// 3. response (no need to response this group command)
	return nil
}
