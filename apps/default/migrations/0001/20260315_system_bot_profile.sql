-- Create the default system bot profile with a fixed ID.
-- This profile is used by all internal service accounts.
-- Contact (system.bot@stawi.org) is added at application startup since
-- contacts require encryption that can't be done in SQL.
INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES (
    'system_bot_profile_01',
    NOW(),
    NOW(),
    1,
    'bjt4h376abi8cg3kgr80',  -- bot profile type
    '{"name": "System Bot", "description": "Default bot profile for internal service accounts"}'
)
ON CONFLICT (id) DO NOTHING;
