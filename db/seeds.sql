begin;

truncate events cascade;

INSERT INTO events (id, title_en, title_pl, starts_at, ends_at, is_virtual, description_en, event_type, base_price_amount, base_price_currency)
VALUES
    ('0199c2f2-528b-7e88-96e3-5e5088333a8c',
     'To Perfect the Art of Homeopathy',
     'Udoskonalić kunszt homeopatyczny',
     '2025-05-30 14:00:00',
     '2025-05-31 08:00:00',
     true,
     'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.',
     'webinar',
     580.00000000,
     'PLN'),

    ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58',
     'To Perfect the Art of Homeopathy 2',
     'Udoskonalić kunszt homeopatyczny 2',
     '2025-10-24 14:00:00',
     '2025-10-26 11:30:00',
     true,
     'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.

October 24-25 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków, Poland.

Online mode will also available (through Zoom). The lectures will be held in English with consecutive translation to Polish.',
     'webinar',
     640.00000000,
     'PLN');

commit;