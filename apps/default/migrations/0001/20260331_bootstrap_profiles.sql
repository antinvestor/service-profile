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
    '{"au_name": "Peter Bwire", "description": "System administrator"}'
)
ON CONFLICT (id) DO NOTHING;

-- Service account profiles (bot type = bjt4h376abi8cg3kgr80)
INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES
    ('d75qclkpf2t1uum8ij40', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-authentication", "description": "Authentication service"}'),
    ('d75qclkpf2t1uum8ij4g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-profile", "description": "Profile service"}'),
    ('d75qclkpf2t1uum8ij50', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-tenancy", "description": "Tenancy service"}'),
    ('d75qclkpf2t1uum8ij5g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-notification", "description": "Notification service"}'),
    ('d75qclkpf2t1uum8ij60', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-device", "description": "Devices service"}'),
    ('d75qclkpf2t1uum8ij6g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-setting", "description": "Settings service"}'),
    ('d75qclkpf2t1uum8ij70', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-payment", "description": "Payment service"}'),
    ('d75qclkpf2t1uum8ij7g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-payment-jenga", "description": "Payment Jenga integration"}'),
    ('d75qclkpf2t1uum8ij80', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-ledger", "description": "Ledger service"}'),
    ('d75qclkpf2t1uum8ij8g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-billing", "description": "Billing service"}'),
    ('d75qclkpf2t1uum8ij90', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-file", "description": "File service"}'),
    ('d75qclkpf2t1uum8ij9g', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-chat-drone", "description": "Chat drone service"}'),
    ('d75qclkpf2t1uum8ija0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-chat-gateway", "description": "Chat gateway service"}'),
    ('d75qclkpf2t1uum8ijag', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"foundry", "description": "Foundry service"}'),
    ('d75qclkpf2t1uum8ijb0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"gitvault", "description": "Gitvault service"}'),
    ('d75qclkpf2t1uum8ijbg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"trustage", "description": "Trustage service"}'),
    ('d75qclkpf2t1uum8ijc0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-notification-africastalking", "description": "Africastalking notification integration"}'),
    ('d75qclkpf2t1uum8ijcg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-notification-emailsmtp", "description": "Email SMTP notification integration"}'),
    ('d75qclkpf2t1uum8ijd0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-lender", "description": "Lender service"}'),
    ('d75qclkpf2t1uum8ijdg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-identity", "description": "Identity and KYC verification service"}'),
    ('d75qclkpf2t1uum8ije0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-loans", "description": "Loan management service"}'),
    ('d75qclkpf2t1uum8ijeg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-origination", "description": "Loan origination service"}'),
    ('d75qclkpf2t1uum8ijf0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-funding", "description": "Loan funding and disbursement service"}'),
    ('d75qclkpf2t1uum8ijfg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-savings", "description": "Savings account management service"}'),
    ('d75qclkpf2t1uum8ijg0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-operations", "description": "Operations and transfer service"}'),
    ('d75qclkpf2t1uum8ijgg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-seed", "description": "Seed direct lending service"}'),
    ('d75qclkpf2t1uum8ijh0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"service-stawi", "description": "Stawi workflow orchestration service"}')
ON CONFLICT (id) DO NOTHING;
