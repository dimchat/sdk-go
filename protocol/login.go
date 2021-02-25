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
package protocol

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "login",
 *      time    : 0,
 *      //---- client info ----
 *      ID       : "{UserID}",
 *      device   : "DeviceID",  // (optional)
 *      agent    : "UserAgent", // (optional)
 *      //---- server info ----
 *      station  : {
 *          ID   : "{StationID}",
 *          host : "{IP}",
 *          port : 9394
 *      },
 *      provider : {
 *          ID   : "{SP_ID}"
 *      }
 *  }
 */
type LoginCommand interface {
	Command
	ILoginCommand
}
type ILoginCommand interface {

	/**
	 *  User ID
	 */
	ID() ID

	/**
	 *  Device ID
	 */
	Device() string
	SetDevice(device string)

	/**
	 *  User Agent
	 */
	Agent() string
	SetAgent(agent string)

	/**
	 *  DIM Network Station
	 */
	StationInfo() map[string]interface{}
	SetStationInfo(station map[string]interface{})

	/**
	 *  DIM Network Service Provider
	 */
	ProviderInfo() map[string]interface{}
	SetProviderInfo(sp map[string]interface{})
}

//
//  Login command implementation
//
type BaseLoginCommand struct {
	BaseCommand
	ILoginCommand
}

func (cmd *BaseLoginCommand) Init(dict map[string]interface{}) *BaseLoginCommand {
	if cmd.BaseCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *BaseLoginCommand) InitWithID(identifier ID) *BaseLoginCommand {
	if cmd.BaseCommand.InitWithCommand(LOGIN) != nil {
		cmd.Set("ID", identifier.String())
	}
	return cmd
}

//-------- ILoginCommand

func (cmd *BaseLoginCommand) ID() ID {
	return IDParse(cmd.Get("ID"))
}

func (cmd *BaseLoginCommand) Device() string {
	device := cmd.Get("device")
	if device == nil {
		return ""
	}
	return device.(string)
}
func (cmd *BaseLoginCommand) SetDevice(device string) {
	if device == "" {
		cmd.Set("device", nil)
	} else {
		cmd.Set("device", device)
	}
}

func (cmd *BaseLoginCommand) Agent() string {
	agent := cmd.Get("agent")
	if agent == nil {
		return ""
	}
	return agent.(string)
}
func (cmd *BaseLoginCommand) SetAgent(agent string) {
	if agent == "" {
		cmd.Set("agent", nil)
	} else {
		cmd.Set("agent", agent)
	}
}

//
//  Server Info
//

func (cmd *BaseLoginCommand) StationInfo() map[string]interface{} {
	station := cmd.Get("station")
	if station == nil {
		return nil
	}
	return station.(map[string]interface{})
}
func (cmd *BaseLoginCommand) SetStationInfo(station map[string]interface{}) {
	if station == nil {
		cmd.Set("station", nil)
	} else {
		cmd.Set("station", station)
	}
}

func (cmd *BaseLoginCommand) ProviderInfo() map[string]interface{} {
	sp := cmd.Get("provider")
	if sp == nil {
		return nil
	}
	return sp.(map[string]interface{})
}
func (cmd *BaseLoginCommand) SetProviderInfo(sp map[string]interface{}) {
	if sp == nil {
		cmd.Set("provider", nil)
	} else {
		cmd.Set("provider", sp)
	}
}
