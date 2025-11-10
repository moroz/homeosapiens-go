-- db/seeds.sql
begin;

truncate events, hosts, assets, venues, events_hosts, event_prices cascade;

-- Insert assets
INSERT INTO assets (id, object_key, original_filename)
VALUES ('0199c2f2-528b-7e88-96e3-5e5088333a8b', 'cm7uqj3q500mglz8z2dqy8sdz.webp', 'cm7uqj3q500mglz8z2dqy8sdz.webp');

-- Insert hosts
INSERT INTO hosts (id, salutation, given_name, family_name, profile_picture_id)
VALUES ('0199c2f2-528b-7e88-96e3-5e5088333a8a', 'common.hosts.salutation.dr', 'Sanjay', 'Modi', '0199c2f2-528b-7e88-96e3-5e5088333a8b');

-- Insert venues
INSERT INTO venues (id, name_en, name_pl, city_en, city_pl, country_code, street, postal_code)
VALUES
    ('0199c2f2-528b-7e88-96e3-5e5088333a8d', 'Vienna House Easy By Wyndham Cracow', NULL, 'Cracow', 'Kraków', 'PL', 'ul. Przy Rondzie 2', NULL),
    ('0199c2f2-528b-7e88-96e3-5e5088333a8e', 'IOR Hotel', 'Hotel IOR', 'Poznań', NULL, 'PL', 'ul. Węgorka 20', '60-318');

-- Insert events
INSERT INTO events (id, title_en, title_pl, slug, starts_at, ends_at, is_virtual, description_en, venue_id, base_price_amount, base_price_currency)
VALUES
    ('0199c2f2-528b-7e88-96e3-5e5088333a8c',
     'To Perfect the Art of Homeopathy',
     'Udoskonalić kunszt homeopatyczny',
     'to-perfect-the-art-of-homeopathy',
     '2025-05-30 14:00:00+00',
     '2025-05-31 08:00:00+00',
     true,
     'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.',
     '0199c2f2-528b-7e88-96e3-5e5088333a8e',
     580.00000000,
     'PLN'),

    ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58',
     'To Perfect the Art of Homeopathy 2',
     'Udoskonalić kunszt homeopatyczny 2',
     'to-perfect-the-art-of-homeopathy-2',
     '2025-10-24 14:00:00+00',
     '2025-10-26 11:30:00+00',
     true,
     'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.

October 24-25 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków, Poland.

Online mode will also available (through Zoom). The lectures will be held in English with consecutive translation to Polish.',
     '0199c2f2-528b-7e88-96e3-5e5088333a8d',
     640.00000000,
     'PLN');

-- Insert event hosts
INSERT INTO events_hosts (event_id, host_id, position)
VALUES
    ('0199c2f2-528b-7e88-96e3-5e5088333a8c', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0),
    ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0);

-- Insert event prices
INSERT INTO event_prices (event_id, price_amount, price_currency, rule_type, valid_until, priority)
VALUES
    ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', 560.00000000, 'PLN', 'EarlyBird', '2025-09-20 21:59:59+00', 10);

INSERT INTO event_prices (event_id, price_amount, price_currency, rule_type, discount_code, priority)
VALUES
    ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', 500.00000000, 'PLN', 'DiscountCode', 'wshlif', 20);
commit;
