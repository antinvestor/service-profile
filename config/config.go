package config

const EnvServerPort = "PORT"

const EnvDatabaseUrl = "DATABASE_URL"
const EnvReplicaDatabaseUrl = "REPLICA_DATABASE_URL"

const EnvMigrate = "DO_MIGRATION"
const EnvMigrationPath = "MIGRATION_PATH"

const EnvNotificationServiceUri = "NOTIFICATION_SERVICE_URI"

const EnvContactEncryptionKey = "CONTACT_ENCRYPTION_KEY"
const EnvContactEncryptionSalt = "CONTACT_ENCRYPTION_SALT"

const EnvQueueVerification = ""
const QueueVerificationName = ""


const LengthOfVerificationPin = 5
const LengthOfVerificationLinkHash = 70
const VerificationPinExpiryTimeInSec = 24 * 60 * 60

const MessageTemplateContactVerification = "template.papi.contact.verification"
