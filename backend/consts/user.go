package consts

const (
	UserActiveKeyFmt  = "user:active:%s"
	AdminActiveKeyFmt = "admin:active:%s"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusLocked   UserStatus = "locked"
)

const (
	SessionName     = "ModelKit_session"
	UserSessionName = "ModelKitser_session"
)

type UserPlatform string

const (
	UserPlatformEmail    UserPlatform = "email"
	UserPlatformDingTalk UserPlatform = "dingtalk"
	UserPlatformCustom   UserPlatform = "custom"
)

type OAuthKind string

const (
	OAuthKindInvite OAuthKind = "invite"
	OAuthKindLogin  OAuthKind = "login"
)

type InviteCodeStatus string

const (
	InviteCodeStatusPending InviteCodeStatus = "pending"
	InviteCodeStatusUsed    InviteCodeStatus = "used"
)

type LoginSource string

const (
	LoginSourcePlugin  LoginSource = "plugin"
	LoginSourceBrowser LoginSource = "browser"
)
