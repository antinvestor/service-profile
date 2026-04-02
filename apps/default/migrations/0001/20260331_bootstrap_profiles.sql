-- Bootstrap profiles for the platform admin and all internal service accounts.
-- Contacts are added at application startup because they require encryption.
--
-- IDs are xid-generated and referenced by the tenancy service's
-- service_accounts.profile_id column.

-- Platform admin profile (person type)
INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES (
    'd75qclkpf2t1uum8ij3g',
    NOW(), NOW(), 1,
    'bjr98v76ad79n5c9p8a0',  -- person profile type
    '{"name": "Platform Admin", "description": "Bootstrap admin account"}'
)
ON CONFLICT (id) DO NOTHING;

-- Service account profiles (bot type = bjt4h376abi8cg3kgr80)
INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES
    ('d75qclkpf2t1uum8ij40', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-authentication", "description": "Authentication service"}'),
    ('d75qclkpf2t1uum8ij4g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-profile", "description": "Profile service"}'),
    ('d75qclkpf2t1uum8ij50', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-tenancy", "description": "Tenancy service"}'),
    ('d75qclkpf2t1uum8ij5g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-notification", "description": "Notification service"}'),
    ('d75qclkpf2t1uum8ij60', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-device", "description": "Devices service"}'),
    ('d75qclkpf2t1uum8ij6g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-setting", "description": "Settings service"}'),
    ('d75qclkpf2t1uum8ij70', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-payment", "description": "Payment service"}'),
    ('d75qclkpf2t1uum8ij7g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-payment-jenga", "description": "Payment Jenga integration"}'),
    ('d75qclkpf2t1uum8ij80', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-ledger", "description": "Ledger service"}'),
    ('d75qclkpf2t1uum8ij8g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-billing", "description": "Billing service"}'),
    ('d75qclkpf2t1uum8ij90', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-file", "description": "File service"}'),
    ('d75qclkpf2t1uum8ij9g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-chat-drone", "description": "Chat drone service"}'),
    ('d75qclkpf2t1uum8ija0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-chat-gateway", "description": "Chat gateway service"}'),
    ('d75qclkpf2t1uum8ijag', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "foundry", "description": "Foundry service"}'),
    ('d75qclkpf2t1uum8ijb0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "gitvault", "description": "Gitvault service"}'),
    ('d75qclkpf2t1uum8ijbg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "trustage", "description": "Trustage service"}'),
    ('d75qclkpf2t1uum8ijc0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-notification-africastalking", "description": "Africastalking notification integration"}'),
    ('d75qclkpf2t1uum8ijcg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-notification-emailsmtp", "description": "Email SMTP notification integration"}'),
    ('d75qclkpf2t1uum8ijd0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"name": "service-lender", "description": "Lender service"}')
ON CONFLICT (id) DO NOTHING;
