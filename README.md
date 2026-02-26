# Decentralized Instant Messaging (Go SDK)

[![License](https://img.shields.io/github/license/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/blob/main/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/dimchat/sdk-go/pulls)
[![Platform](https://img.shields.io/github/go-mod/go-version/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/wiki)
[![Issues](https://img.shields.io/github/issues/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/issues)
[![Repo Size](https://img.shields.io/github/repo-size/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/archive/refs/heads/main.zip)
[![Tags](https://img.shields.io/github/tag/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/tags)

[![Watchers](https://img.shields.io/github/watchers/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/watchers)
[![Forks](https://img.shields.io/github/forks/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/forks)
[![Stars](https://img.shields.io/github/stars/dimchat/sdk-go)](https://github.com/dimchat/sdk-go/stargazers)
[![Followers](https://img.shields.io/github/followers/dimchat)](https://github.com/orgs/dimchat/followers)

## Dependencies

* Latest Versions

| Name | Version | Description |
|------|---------|-------------|
| [Ming Ke Ming (名可名)](https://github.com/dimchat/mkm-go) | [![Tags](https://img.shields.io/github/tag/dimchat/mkm-go)](https://github.com/dimchat/mkm-go/tags) | Decentralized User Identity Authentication |
| [Dao Ke Dao (道可道)](https://github.com/dimchat/dkd-go) | [![Tags](https://img.shields.io/github/tag/dimchat/dkd-go)](https://github.com/dimchat/dkd-go/tags) | Universal Message Module |
| [DIMP (去中心化通讯协议)](https://github.com/dimchat/core-go) | [![Tags](https://img.shields.io/github/tag/dimchat/core-go)](https://github.com/dimchat/core-go/tags) | Decentralized Instant Messaging Protocol |

## Extensions

### Content

extends [CustomizedContent](https://github.com/dimchat/core-go#extends-content)

### ContentProcessor

```go
import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
)

// CustomizedContentProcessor is the ContentProcessor implementation for CustomizedContent
//
// # Integrates with CustomizedContentFilter and CustomizedContentHandler to process app-specific content
//
// Serves as the bridge between the core message processing system and application-specific logic
type CustomizedContentProcessor struct {
	*BaseContentProcessor
}

// Override
func (cpu *CustomizedContentProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	customized, ok := content.(CustomizedContent)
	if !ok {
		//panic("customized content error")
		return nil
	}
	// get handler for 'app' & 'mod'
	filter := GetCustomizedContentFilter()
	handler := filter.FilterContent(customized, rMsg)
	if handler == nil {
		//panic("should not happen")
		return nil
	}
	// handle the action
	messenger := cpu.Messenger
	return handler.HandleContent(customized, rMsg, messenger)
}
```

CustomizedContentHandler + CustomizedContentFilter

```go
import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimchat/sdk-go/sdk"
)

/*
Command Transform:

	+===============================+===============================+
	|      Customized Content       |      Group Query Command      |
	+-------------------------------+-------------------------------+
	|   "type" : i2s(0xCC)          |   "type" : i2s(0x88)          |
	|   "sn"   : 123                |   "sn"   : 123                |
	|   "time" : 123.456            |   "time" : 123.456            |
	|   "app"  : "chat.dim.group"   |                               |
	|   "mod"  : "history"          |                               |
	|   "act"  : "query"            |                               |
	|                               |   "command"   : "query"       |
	|   "group"     : "{GROUP_ID}"  |   "group"     : "{GROUP_ID}"  |
	|   "last_time" : 0             |   "last_time" : 0             |
	+===============================+===============================+
*/
type GroupHistoryHandler struct {
	BaseCustomizedHandler
}

// Override
func (handler GroupHistoryHandler) HandleContent(content CustomizedContent, rMsg ReliableMessage, messenger Messenger) []Content {
	if content.Group() == nil {
		return handler.RespondReceipt("group command error.", rMsg.Envelope(), content, nil)
	}
	act := content.Action()
	if act == "query" {
		return handler.transformQueryCommand(content, rMsg, messenger)
	}
	//panic("unknown action: " + act)
	return handler.BaseCustomizedHandler.HandleContent(content, rMsg, messenger)
}

func (handler GroupHistoryHandler) transformQueryCommand(content CustomizedContent, rMsg ReliableMessage, messenger Messenger) []Content {
	//info := content.CopyMap(false)
	//info["type"] = ContentType.COMMAND
	//info["command"] = "query"
	//query := ParseContent(info)
	//if command, ok := query.(QueryCommand); ok {
	//	return messenger.ProcessContent(command, rMsg)
	//}
	return handler.RespondReceipt("Query Command error.", rMsg.Envelope(), content, nil)
}

type AppCustomizedFilter struct {
	//CustomizedContentFilter

	defaultHandler CustomizedContentHandler
	handlers       map[string]CustomizedContentHandler
}

func NewAppCustomizedFilter() *AppCustomizedFilter {
	return &AppCustomizedFilter{
		defaultHandler: &BaseCustomizedHandler{},
		handlers:       make(map[string]CustomizedContentHandler, 8),
	}
}

func (filter *AppCustomizedFilter) SetContentHandler(app, mod string, handler CustomizedContentHandler) {
	key := app + ":" + mod
	filter.handlers[key] = handler
}

// protected
func (filter *AppCustomizedFilter) GetContentHandler(app, mod string) CustomizedContentHandler {
	key := app + ":" + mod
	return filter.handlers[key]
}

// Override
func (filter *AppCustomizedFilter) FilterContent(content CustomizedContent, _ ReliableMessage) CustomizedContentHandler {
	app := content.Application()
	mod := content.Module()
	handler := filter.GetContentHandler(app, mod)
	if handler == nil {
		return filter.defaultHandler
	}
	return handler
}

func init() {
	filter := NewAppCustomizedFilter()

	// 'chat.dim.group:history'
	handler := &GroupHistoryHandler{}
	filter.SetContentHandler("chat.dim.group", "history", handler)

	SetCustomizedContentFilter(filter)
}
```

### ContentProcessorCreator

```go
import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimchat/sdk-go/dkd"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  CPU Creator
 *  ~~~~~~~~~~~
 *  Delegate for CPU factory
 */

type ClientContentProcessorCreator struct {
	*BaseContentProcessorCreator
}

//-------- IProcessorCreator

func (creator *ClientContentProcessorCreator) CreateContentProcessor(msgType MessageType) ContentProcessor {
	switch msgType {
	// application customized
	case "application", "customized":
		return NewCustomizedContentProcessor(creator.Facebook, creator.Messenger)
	
	// ...
	}
	// others
	return creator.BaseContentProcessorCreator.CreateContentProcessor(msgType)
}

func (creator *ClientContentProcessorCreator) CreateCommandProcessor(msgType MessageType, cmdName string) ContentProcessor {
	switch cmdName {
	// handshake command
	case HANDSHAKE:
		return NewHandshakeCommandProcessor(creator.Facebook, creator.Messenger)
	
	// ...
	}
	// others
	return creator.BaseContentProcessorCreator.CreateCommandProcessor(msgType, cmdName)
}

//
//  Factories
//

func NewCustomizedContentProcessor(facebook Facebook, messenger Messenger) *CustomizedContentProcessor {
	return &CustomizedContentProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}

func NewHandshakeCommandProcessor(facebook Facebook, messenger Messenger) ContentProcessor {
	return &HandshakeCommandProcessor{
		BaseCommandProcessor: NewBaseCommandProcessor(facebook, messenger),
	}
}
```

## Usage

To let your **AppCustomizedProcessor** start to work,
you must override ```BaseContentProcessorCreator``` for message types:

1. ContentType.APPLICATION 
2. ContentType.CUSTOMIZED

and then set your **creator** for ```GeneralContentProcessorFactory``` in the ```MessageProcessor```.

----

Copyright &copy; 2020-2026 Albert Moky
[![Followers](https://img.shields.io/github/followers/moky)](https://github.com/moky?tab=followers)
