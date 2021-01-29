

INSERT INTO profile_types (id, uid, name, description) VALUES
('bjr98v76abi9n5c9p8a0', 0, 'person', 'Human beings using the system'),
('bjsml5v6abi8e2809or0', 1, 'institution', 'Organized legal entities'),
('bjt4h376abi8cg3kgr80', 2, 'bot', 'Robots of all types');

INSERT INTO contact_types (id, uid, name, description) VALUES
('bjr98v76abi9n5c9p8a0', 0, 'email', 'Email address of a profile'),
('bjsml5v6abi8e2809or0', 1, 'phone', 'Phone number associated to profile');

INSERT INTO communication_levels (id, uid, name, description) VALUES
('bjr98v76abi9n5c9p8a0', 0, 'all', 'The system can send any kind of communication to such a contact'),
('bjsml5v6abi8e2809or0', 1, 'system alerts', 'Only messages resulting from the users actions can be sent to them'),
('bjt4h376abi8cg3kgr80', 2, 'no contact', 'This one may be a terrorist as we can not even contact them' );
