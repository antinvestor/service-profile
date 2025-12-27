package models

import (
	"testing"

	profilev1 "buf.build/gen/go/antinvestor/profile/protocolbuffers/go/profile/v1"
	"github.com/antinvestor/service-profile/apps/default/config"
	"github.com/pitabwire/frame/data"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
)

func TestProfileTypeIDToEnum(t *testing.T) {
	tests := []struct {
		name     string
		typeID   uint
		expected profilev1.ProfileType
	}{
		{"Person type", ProfileTypePersonID, profilev1.ProfileType_PERSON},
		{"Bot type", ProfileTypeBotID, profilev1.ProfileType_BOT},
		{"Institution type", ProfileTypeInstitutionID, profilev1.ProfileType_INSTITUTION},
		{"Unknown type defaults to Person", 999, profilev1.ProfileType_PERSON},
		{"Zero type defaults to Person", 0, profilev1.ProfileType_PERSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProfileTypeIDToEnum(tt.typeID)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRelationshipTypeIDToEnum(t *testing.T) {
	tests := []struct {
		name     string
		typeID   uint
		expected profilev1.RelationshipType
	}{
		{"Member type", RelationshipTypeMemberID, profilev1.RelationshipType_MEMBER},
		{"Affiliated type", RelationshipTypeAffiliatedID, profilev1.RelationshipType_AFFILIATED},
		{"Blacklisted type", RelationshipTypeBlackListedID, profilev1.RelationshipType_BLACK_LISTED},
		{"Unknown type defaults to Member", 999, profilev1.RelationshipType_MEMBER},
		{"Zero type defaults to Member", 0, profilev1.RelationshipType_MEMBER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RelationshipTypeIDToEnum(tt.typeID)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestContact_DecryptDetail(t *testing.T) {
	// Create a test encryption key (32 bytes for AES-256)
	key := []byte("12345678901234567890123456789012")
	keyID := "test-key-id"

	// Encrypt a test value
	testDetail := "test@example.com"
	encryptedDetail, err := util.EncryptValue(key, []byte(testDetail))
	require.NoError(t, err)

	tests := []struct {
		name          string
		contact       *Contact
		decryptKeyID  string
		decryptKey    []byte
		expectedValue string
		wantErr       bool
	}{
		{
			name: "Successful decryption",
			contact: &Contact{
				EncryptedDetail: encryptedDetail,
				EncryptionKeyID: keyID,
			},
			decryptKeyID:  keyID,
			decryptKey:    key,
			expectedValue: testDetail,
			wantErr:       false,
		},
		{
			name: "Wrong key ID",
			contact: &Contact{
				EncryptedDetail: encryptedDetail,
				EncryptionKeyID: keyID,
			},
			decryptKeyID:  "wrong-key-id",
			decryptKey:    key,
			expectedValue: "",
			wantErr:       true,
		},
		{
			name: "Wrong decryption key",
			contact: &Contact{
				EncryptedDetail: encryptedDetail,
				EncryptionKeyID: keyID,
			},
			decryptKeyID:  keyID,
			decryptKey:    []byte("wrong-key-wrong-key-wrong-key123"),
			expectedValue: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.contact.DecryptDetail(tt.decryptKeyID, tt.decryptKey)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, result)
		})
	}
}

func TestContact_ToAPI(t *testing.T) {
	// Create a test encryption key (32 bytes for AES-256)
	key := []byte("12345678901234567890123456789012")
	keyID := "test-key-id"

	// Encrypt a test value
	testDetail := "test@example.com"
	encryptedDetail, err := util.EncryptValue(key, []byte(testDetail))
	require.NoError(t, err)

	dek := &config.DEK{
		KeyID:     keyID,
		Key:       key,
		LookUpKey: []byte("lookup-key"),
	}

	tests := []struct {
		name              string
		contact           *Contact
		partial           bool
		expectedDetail    string
		expectedType      profilev1.ContactType
		expectedVerified  bool
		wantErr           bool
	}{
		{
			name: "Full contact with verification",
			contact: &Contact{
				BaseModel:          data.BaseModel{ID: "contact-1"},
				EncryptedDetail:    encryptedDetail,
				EncryptionKeyID:    keyID,
				ContactType:        "EMAIL",
				CommunicationLevel: "ALL",
				VerificationID:     "verification-1",
			},
			partial:          false,
			expectedDetail:   testDetail,
			expectedType:     profilev1.ContactType_EMAIL,
			expectedVerified: true,
			wantErr:          false,
		},
		{
			name: "Partial contact hides verification",
			contact: &Contact{
				BaseModel:          data.BaseModel{ID: "contact-2"},
				EncryptedDetail:    encryptedDetail,
				EncryptionKeyID:    keyID,
				ContactType:        "MSISDN",
				CommunicationLevel: "MARKETING",
				VerificationID:     "verification-2",
			},
			partial:          true,
			expectedDetail:   testDetail,
			expectedType:     profilev1.ContactType_MSISDN,
			expectedVerified: false,
			wantErr:          false,
		},
		{
			name: "Unknown contact type defaults to EMAIL",
			contact: &Contact{
				BaseModel:          data.BaseModel{ID: "contact-3"},
				EncryptedDetail:    encryptedDetail,
				EncryptionKeyID:    keyID,
				ContactType:        "UNKNOWN",
				CommunicationLevel: "INVALID",
			},
			partial:          false,
			expectedDetail:   testDetail,
			expectedType:     profilev1.ContactType_EMAIL,
			expectedVerified: false,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.contact.ToAPI(dek, tt.partial)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, tt.expectedDetail, result.GetDetail())
			require.Equal(t, tt.expectedType, result.GetType())
			require.Equal(t, tt.expectedVerified, result.GetVerified())
			require.Equal(t, tt.contact.GetID(), result.GetId())
		})
	}
}

func TestRoster_ToAPI(t *testing.T) {
	// Create a test encryption key (32 bytes for AES-256)
	key := []byte("12345678901234567890123456789012")
	keyID := "test-key-id"

	// Encrypt a test value
	testDetail := "+256757546244"
	encryptedDetail, err := util.EncryptValue(key, []byte(testDetail))
	require.NoError(t, err)

	dek := &config.DEK{
		KeyID:     keyID,
		Key:       key,
		LookUpKey: []byte("lookup-key"),
	}

	contact := &Contact{
		BaseModel:          data.BaseModel{ID: "contact-1"},
		EncryptedDetail:    encryptedDetail,
		EncryptionKeyID:    keyID,
		ContactType:        "MSISDN",
		CommunicationLevel: "ALL",
	}

	roster := &Roster{
		BaseModel:  data.BaseModel{ID: "roster-1"},
		ProfileID:  "profile-123",
		ContactID:  contact.ID,
		Contact:    contact,
		Properties: data.JSONMap{"name": "Test User"},
	}

	result, err := roster.ToAPI(dek)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "roster-1", result.GetId())
	require.Equal(t, "profile-123", result.GetProfileId())
	require.NotNil(t, result.GetContact())
	require.Equal(t, testDetail, result.GetContact().GetDetail())
	require.NotNil(t, result.GetExtra())
}

func TestRelationship_ToAPI(t *testing.T) {
	relationshipType := &RelationshipType{
		BaseModel:   data.BaseModel{ID: "type-1"},
		UID:         RelationshipTypeMemberID,
		Name:        "member",
		Description: "Member relationship",
	}

	relationship := &Relationship{
		BaseModel:          data.BaseModel{ID: "relationship-1"},
		ParentObject:       "Profile",
		ParentObjectID:     "parent-profile-id",
		ChildObject:        "Profile",
		ChildObjectID:      "child-profile-id",
		RelationshipTypeID: relationshipType.ID,
		RelationshipType:   relationshipType,
		Properties:         data.JSONMap{"role": "admin"},
	}

	result := relationship.ToAPI()
	require.NotNil(t, result)
	require.Equal(t, "relationship-1", result.GetId())
	require.Equal(t, profilev1.RelationshipType(RelationshipTypeMemberID), result.GetType())
	require.NotNil(t, result.GetChildEntry())
	require.Equal(t, "Profile", result.GetChildEntry().GetObjectName())
	require.Equal(t, "child-profile-id", result.GetChildEntry().GetObjectId())
	require.NotNil(t, result.GetParentEntry())
	require.Equal(t, "Profile", result.GetParentEntry().GetObjectName())
	require.Equal(t, "parent-profile-id", result.GetParentEntry().GetObjectId())
	require.NotNil(t, result.GetProperties())
}

func TestProfileTypeIDMap(t *testing.T) {
	// Verify the map contains expected entries
	require.Equal(t, ProfileTypePersonID, ProfileTypeIDMap[profilev1.ProfileType_PERSON])
	require.Equal(t, ProfileTypeBotID, ProfileTypeIDMap[profilev1.ProfileType_BOT])
	require.Equal(t, ProfileTypeInstitutionID, ProfileTypeIDMap[profilev1.ProfileType_INSTITUTION])
}

func TestRelationshipTypeIDMap(t *testing.T) {
	// Verify the map contains expected entries
	require.Equal(t, RelationshipTypeMemberID, RelationshipTypeIDMap[profilev1.RelationshipType_MEMBER])
	require.Equal(t, RelationshipTypeAffiliatedID, RelationshipTypeIDMap[profilev1.RelationshipType_AFFILIATED])
	require.Equal(t, RelationshipTypeBlackListedID, RelationshipTypeIDMap[profilev1.RelationshipType_BLACK_LISTED])
}
