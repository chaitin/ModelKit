package consts

type ConfigType string

const (
	ConfigTypeContinue ConfigType = "continue"
	ConfigTypeRooCode  ConfigType = "roo-code"
)

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)
