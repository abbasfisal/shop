package logging

type Category string
type SubCategory string
type ExtraKey string

const (
	General         Category = "General"
	IO              Category = "IO"
	Internal        Category = "Internal"
	Database        Category = "Database"
	Redis           Category = "Redis"
	ElasticSearch   Category = "ElasticSearch"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Prometheus      Category = "Prometheus"
	JWT             Category = "JWT"
	Notification    Category = "Notification"
	Twilio          Category = "Twilio"
	Vonage          Category = "Vonage"
	Ultramsg        Category = "Ultramsg"
	Email           Category = "Email"
	Slack           Category = "Slack"
	Google          Category = "Google"
	Facebook        Category = "Facebook"
	Apple           Category = "Apple"
	Queue           Category = "Queue"
	SQS             Category = "SQS"
	Mailchimp       Category = "Mailchimp"
	AWS             Category = "AWS"
)

const (
	InternalInfo SubCategory = "InternalInfo"

	Startup         SubCategory = "Startup"
	ExternalService SubCategory = "ExternalService"

	API                 SubCategory = "API"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"

	DatabaseConnectionError SubCategory = "DatabaseConnectionError"
	DatabaseQueryError      SubCategory = "DatabaseQueryError"
	DatabaseSelect          SubCategory = "DatabaseSelect"
	DatabaseInsert          SubCategory = "DatabaseInsert"
	DatabaseUpdate          SubCategory = "DatabaseUpdate"
	DatabaseDelete          SubCategory = "DatabaseDelete"
	DatabaseRollback        SubCategory = "DatabaseRollback"
	DatabaseMigration       SubCategory = "DatabaseMigration"

	RedisConnection SubCategory = "RedisConnection"
	RedisSet        SubCategory = "RedisSet"
	RedisGet        SubCategory = "RedisGet"
	RedisDel        SubCategory = "RedisDel"
	RedisPing       SubCategory = "RedisPing"

	ElasticSearchQueryError SubCategory = "ElasticSearchQueryError"
	ElasticSearchIndexing   SubCategory = "ElasticSearchIndexing"
	ElasticSearchSearch     SubCategory = "ElasticSearchSearch"
	ElasticSearchSuggest    SubCategory = "ElasticSearchSuggest"
	ElasticSearchUpdate     SubCategory = "ElasticSearchUpdate"
	ElasticSearchDelete     SubCategory = "ElasticSearchDelete"
	ElasticSearchFlush      SubCategory = "ElasticSearchFlush"
	ElasticSearchMapping    SubCategory = "ElasticSearchMapping"

	ValidationFailed SubCategory = "ValidationFailed"

	RequestError SubCategory = "RequestError"

	RemoveFile SubCategory = "RemoveFile"

	JWTGenerate SubCategory = "JWTGenerate"

	NotificationSend SubCategory = "NotificationSend"

	SlackSendMessage SubCategory = "SlackSendMessage"

	TwilioWebhook     SubCategory = "TwilioWebhook"
	TwilioSendSMS     SubCategory = "TwilioSendSMS"
	TwilioCheck       SubCategory = "TwilioCheck"
	TwilioRetrySMS    SubCategory = "TwilioRetrySMS"
	TwilioUpdateSMS   SubCategory = "TwilioUpdateSMS"
	VonageWebhook     SubCategory = "VonageWebhook"
	VonageSendSMS     SubCategory = "VonageSendSMS"
	VonageCheck       SubCategory = "VonageCheck"
	VonageRetrySMS    SubCategory = "VonageRetrySMS"
	VonageUpdateSMS   SubCategory = "VonageUpdateSMS"
	UltramsgWebhook   SubCategory = "UltramsgWebhook"
	UltramsgSend      SubCategory = "UltramsgSend"
	UltramsgCheck     SubCategory = "UltramsgCheck"
	UltramsgUpdateSMS SubCategory = "UltramsgUpdateSMS"
	UltramsgDelete    SubCategory = "UltramsgDelete"
	EmailSend         SubCategory = "EmailSend"

	GoogleLogin   SubCategory = "GoogleLogin"
	FacebookLogin SubCategory = "FacebookLogin"
	AppleLogin    SubCategory = "AppleLogin"

	SQSSend             SubCategory = "SQSSend"
	SQSPublish          SubCategory = "SQSPublish"
	SQSRegisterConsumer SubCategory = "SQSRegisterConsumer"
	SQSConsume          SubCategory = "SQSConsume"

	MailchimpAddToMailingList SubCategory = "MailchimpAddToMailingList"

	AWSInit         SubCategory = "AWSInit"
	AwsS3Connection SubCategory = "AwsS3Connection"

	DataConversion SubCategory = "DataConversion"
)

const (
	AppName         ExtraKey = "AppName"
	LoggerName      ExtraKey = "LoggerName"
	ClientIp        ExtraKey = "ClientIp"
	HostIp          ExtraKey = "HostIp"
	Method          ExtraKey = "Method"
	StatusCode      ExtraKey = "StatusCode"
	BodySize        ExtraKey = "BodySize"
	Path            ExtraKey = "Path"
	Latency         ExtraKey = "Latency"
	Body            ExtraKey = "Body"
	ErrorMessages   ExtraKey = "ErrorMessages"
	Headers         ExtraKey = "Headers"
	RequestBody     ExtraKey = "RequestBody"
	ResponseBody    ExtraKey = "ResponseBody"
	ErrorMessage    ExtraKey = "ErrorMessage"
	ExtraKeyRequest ExtraKey = "Request"
	ExtraKeyPhone   ExtraKey = "Phone"
	ExtraKeyText    ExtraKey = "Text"
	Migrations      ExtraKey = "Migrations"
)
