package config

const EnvServerPort = "PORT"

const EnvDatabaseURL = "DATABASE_URL"
const EnvReplicaDatabaseURL = "REPLICA_DATABASE_URL"

const EnvMigrate = "DO_MIGRATION"
const EnvMigrationPath = "MIGRATION_PATH"

const EnvNotificationServiceURI = "NOTIFICATION_SERVICE_URI"

const EnvOauth2ServiceURI = "OAUTH2_SERVICE_URI"
const EnvOauth2ServiceClientSecret = "OAUTH2_SERVICE_CLIENT_SECRET"
const EnvOauth2ServiceAudience = "OAUTH2_SERVICE_AUDIENCE"

const EnvContactEncryptionKey = "CONTACT_ENCRYPTION_KEY"
const EnvContactEncryptionSalt = "CONTACT_ENCRYPTION_SALT"

const EnvOauth2JwtVerifyAudience = "OAUTH2_JWT_VERIFY_AUDIENCE"
const EnvOauth2JwtVerifyIssuer = "OAUTH2_JWT_VERIFY_ISSUER"

const EnvQueueVerification = ""
const QueueVerificationName = ""

const LengthOfVerificationPin = 5
const LengthOfVerificationLinkHash = 70
const VerificationPinExpiryTimeInSec = 24 * 60 * 60

const MessageTemplateContactVerification = "template.papi.contact.verification"
