-- Bootstrap profiles for stawi-jobs internal service accounts.
-- Contacts are added at application startup because they require encryption.
--
-- Profile ID Reference (stawi-jobs services):
--   d75qclkpf2t1uum8ijhg  stawi-jobs-candidates   stawi-jobs-candidates.bot@stawi.org
--   d75qclkpf2t1uum8iji0  stawi-jobs-crawler       stawi-jobs-crawler.bot@stawi.org
--   d75qclkpf2t1uum8ijig  stawi-jobs-api           stawi-jobs-api.bot@stawi.org
--   d75qclkpf2t1uum8ijj0  stawi-jobs-scheduler     stawi-jobs-scheduler.bot@stawi.org

INSERT INTO profiles (id, created_at, modified_at, version, profile_type_id, properties)
VALUES
    ('d75qclkpf2t1uum8ijhg', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"stawi-jobs-candidates", "description": "Stawi Jobs candidate matching and delivery service"}'),
    ('d75qclkpf2t1uum8iji0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"stawi-jobs-crawler", "description": "Stawi Jobs crawling and pipeline service"}'),
    ('d75qclkpf2t1uum8ijig', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"stawi-jobs-api", "description": "Stawi Jobs public search API"}'),
    ('d75qclkpf2t1uum8ijj0', NOW(), NOW(), 1, 'bjt4h376abi8cg3kgr80',
     '{"au_name":"stawi-jobs-scheduler", "description": "Stawi Jobs scheduling service"}')
ON CONFLICT (id) DO NOTHING;
