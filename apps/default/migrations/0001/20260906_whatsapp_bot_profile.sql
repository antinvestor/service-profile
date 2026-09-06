-- Bootstrap profile for the WhatsApp notification integration service account.
-- The contact is added at application startup because it requires encryption.
--
-- Profile ID Reference:
--   daenc7kpf2t8pa1q04hg  service-notification-whatsapp  notification-whatsapp.bot@stawi.org

INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES
    ('daenc7kpf2t8pa1q04hg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-notification-whatsapp", "description": "WhatsApp notification integration"}')
ON CONFLICT (id) DO NOTHING;
