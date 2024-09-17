-- Fill with dummy data
INSERT INTO users(fullname, username, password_hash) VALUES
	('Ram Poudel', 'ram', '$2a$10$vT5JwJLYO6O7Tg5sIMeca.6MIKQGpybEwiLwWvo4n5Ae.TUMeG73y'), -- ram (bcrypted with default cost)
	('Shyam Acharya', 'shyam', '$2a$10$kvRD/jewfMPNt3oMKMGPpe9Hz6pTbWOI1GTfEqT9KsZw160WR5TZO'), -- shyam
	('Hari Rai', 'hari', '$2a$10$U32d4KQn4nR.vBEcHKIwt.WS3DGU16SRlNUSg4pDO3iqAzX23aKfO'), -- hari
	('Mohan Devkota', 'mohan', '$2a$10$bonHqtMXx4V6q550CE4AwuhPx5kF2RPCf7QqJZLSdd9aguVfbn4Ta'), -- mohan
	('Rita Gurung', 'rita', '$2a$10$v49I8/eiUOf.u6jflGiFw.zSlIB0NLbJFB8yB7PZbAUllxPc737bC'); -- rita

INSERT INTO friends VALUES
	(1, 2),
	(1, 3),
	(4, 5);

INSERT INTO messages (sender_id, receiver_id, text, timestamp) VALUES
	(1, 2, 'Hey Shyam, want to go hiking this weekend?', '2024-05-25 14:23:00+00'),
	(2, 1, 'Sure Ram, sounds great!', '2024-05-25 14:25:00+00'),
	(1, 2, 'Awesome! I found a new trail we can explore.', '2024-05-25 14:30:00+00'),
	(2, 1, 'That sounds exciting! How long is the trail?', '2024-05-25 14:32:00+00'),
	(1, 2, 'It''s about 5 miles round trip, with some great views.', '2024-05-25 14:35:00+00'),
	(2, 1, 'Perfect, I can bring some snacks for us.', '2024-05-25 14:37:00+00'),
	(1, 2, 'Great! I''ll bring the water and my camera.', '2024-05-25 14:40:00+00'),
	(2, 1, 'Cool, what time should we meet?', '2024-05-25 14:42:00+00'),
	(1, 2, 'How about 8 AM at the trailhead?', '2024-05-25 14:45:00+00'),
	(2, 1, 'Sounds good to me. See you then!', '2024-05-25 14:47:00+00'),
	(1, 2, 'By the way, did you finish that project at work?', '2024-05-25 15:00:00+00'),
	(2, 1, 'Yes, finally! It took forever, but I''m glad it''s done.', '2024-05-25 15:05:00+00'),
	(1, 2, 'Congrats! I know you worked hard on it.', '2024-05-25 15:10:00+00'),
	(2, 1, 'Thanks, Ram. Your support means a lot.', '2024-05-25 15:15:00+00'),
	(1, 2, 'Anytime, Shyam. ', '2024-05-25 15:20:00+00'),
	(1, 2, 'Friends always have each other''s backs.',  '2024-05-25 15:23:00+00'),
	(1, 2, 'Right?',  '2024-05-25 15:24:00+00'),
	(2, 1, 'Absolutely.', '2024-05-25 15:25:00+00'),
	(2, 1, 'By the way, have you seen the new movie out?', '2024-05-25 15:27:00+00'),
	(1, 2, 'Not yet, but I heard it''s really good. Want to watch it together?', '2024-05-25 15:30:00+00'),
	(2, 1, 'Sure, let''s plan for next weekend.', '2024-05-25 15:35:00+00'),
	(1, 2, 'Perfect. I''ll check the showtimes and let you know.', '2024-05-25 15:40:00+00'),
	(2, 1, 'Sounds like a plan. Talk to you later!', '2024-05-25 15:45:00+00'),
	(1, 3, 'Hi Carol, do you have any new photo tips?', '2024-05-25 15:00:00+00'),
	(3, 1, 'Hey Alice, I do! Let''s chat more.', '2024-05-25 15:05:00+00'),
	(4, 5, 'Eve, have you heard the new album by XYZ?', '2024-05-25 16:10:00+00'),
	(5, 4, 'Yes, David!', '2024-05-25 16:15:00+00'),
	(5, 4, 'It''s amazing!', '2024-05-25 16:20:00+00');

INSERT INTO conversations (user1_id, user2_id, timestamp) VALUES
	(1, 2, '2024-05-25 15:45:00+00'),
	(1, 3, '2024-05-25 15:05:00+00'),
	(4, 5, '2024-05-25 16:15:00+00')
	ON CONFLICT (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id))
	DO UPDATE SET timestamp = excluded.timestamp;

INSERT INTO friend_requests(requestor_id, receiver_id) VALUES
	(4, 1);	-- mohan -> ram
