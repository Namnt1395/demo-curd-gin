package constant

const DefaultPageSize = 10
const DefaultPage = 1
const DefaultPageSort = "created_at desc"
const HeaderAcceptLanguage = "Accept-Language"
const DefaultLang = "en"
const DefaultEnv = "PROD"
const CharSetUtf8 = "UTF-8"
const (
	USER_ID      = "user_id"
	COMPANY_ID   = "company_id"
	COMPANY_CODE = "company_code"
	BRANCH_ID    = "branch_id"
	BRANCH_CODE  = "branch_code"
	BIZAPP_ALIAS = "bizapp_alias"
	BIZAPP_ID    = "bizapp_id"
)

const ConfigPath = "./config"
const EnvKey = "ENVIRONMENT"

type Lang string

const (
	En Lang = "en"
	Vi Lang = "vi"
)

const I18nMessage = "messages"

type SecurityAccess string

const (
	AccessHasPermission SecurityAccess = "HasPermission"
	AccessHasRole       SecurityAccess = "HasRole"
	AccessPermitAll     SecurityAccess = "PermitAll"
	AccessDenyAll       SecurityAccess = "DenyAll"
	AccessCustom        SecurityAccess = "Custom"
)
