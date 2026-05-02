begin;

truncate users, hosts, assets, products, video_groups cascade;

insert into assets (id, object_key, original_filename) values
  ('0199c2f2-528b-7e88-96e3-5e5088333a8b', 'cm7uqj3q500mglz8z2dqy8sdz.webp', 'cm7uqj3q500mglz8z2dqy8sdz.webp'),
  ('019b0c7c-c3c4-71c3-a630-7b33a847ca2a', '019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg', '019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg'),
  ('019beef9-ad4c-736f-9bb0-965b59ca21ae', '019beef9-ad4c-736f-9bb0-965b59ca21ae.png', 'drasher.png');

insert into hosts (id, salutation, given_name, family_name, profile_picture_id, country) values
  ('0199c2f2-528b-7e88-96e3-5e5088333a8a', 'common.hosts.salutation.dr', 'Sanjay', 'Modi',   '0199c2f2-528b-7e88-96e3-5e5088333a8b', 'IN'),
  ('019beef9-4287-714f-982b-2524fdef7063', 'common.hosts.salutation.dr', 'Asher',  'Shaikh', '019beef9-ad4c-736f-9bb0-965b59ca21ae', 'IN'),
  ('019b0c71-fde2-76b7-8c71-21c2e9ea23a5', 'common.hosts.salutation.dr', 'Herman', 'Jeggels','019b0c7c-c3c4-71c3-a630-7b33a847ca2a', 'ZA');

insert into products (id, product_type, title_en, title_pl, base_price_amount, base_price_currency) values
  ('019de49f-d17f-7435-b166-b8e9b3e4430c', 'event', 'Dr Asher Shaikh seminar',              'Seminarium z drem Asherem Shaikh',           560.00, 'PLN'),
  ('019de49f-d3d5-727b-915d-2cd3671cb72f', 'event', 'To Perfect the Art of Homeopathy',     'Udoskonalić kunszt homeopatyczny',           580.00, 'PLN'),
  ('019de49f-d612-7339-962b-800168a04bbb', 'event', 'To Perfect the Art of Homeopathy 2',   'Udoskonalić kunszt homeopatyczny 2',         640.00, 'PLN');

-- Dr Asher Shaikh seminar (paid, in-person)
insert into events (id, product_id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual,
  venue_name_en, venue_name_pl, venue_city_en, venue_city_pl, venue_country_code, venue_street, venue_postal_code,
  description_en, description_pl)
values (
  '019c5c9a-c5a4-7518-8317-65ae90516726',
  '019de49f-d17f-7435-b166-b8e9b3e4430c',
  'seminar',
  'Dr Asher Shaikh seminar',
  'Seminarium z drem Asherem Shaikh',
  'dr-asher-shaikh-seminar',
  '2026-06-05 14:00:00+00', '2026-06-06 20:00:00+00',
  false,
  'Marina Club Hotel', 'Hotel Marina Club',
  'Gdańsk', 'Gdańsk',
  'PL', 'ul. Szafarnia 10', '80-753',
  $en$We kindly invite you to the Homeo sapiens seminar with Dr Asher Shaikh. During the seminar, Dr Shaikh will present the practical application of German New Medicine (GNM) in clinical homeopathic practice, with particular focus on case-taking, identification of biological conflicts, and remedy selection. The seminar will include both theoretical background and case studies. **More info coming soon!**

The seminar will take place on **June 5–6, 2026**.

Venue: **Hotel Marina Club**, ul. Szafarnia 10, 80-753 Gdańsk, Poland + **online** (Zoom).

Participation fee:

*   **EARLY BIRD 480 PLN / 114 EUR** – until Feb 9, 2026
*   **560 PLN / 135 EUR** – until March 31, 2026
*   **640 PLN / 152 EUR** – after March 31, 2026
*   **700 PLN / 170 EUR** – on the day of the seminar
Discounted accommodation is available for participants at **Hotel Marina Club** ([https://marinaclubhotel.pl](https://marinaclubhotel.pl)) and **Hostel Szafarnia** ([https://szafarnia10.pl](https://szafarnia10.pl)). Reservations must be made online. Discount code: **Homeopatia** (active after activation in the hotel booking system).

**Dr Asher Shaikh (India)** is a homeopathic doctor with over 30 years of clinical experience. He is the Director of Asher Clinics, a network of 12 clinics in Mumbai, Pune, Dubai, and Nasik. He is a mentor of German New Medicine, which he has taught in Dubai, India, Austria, and Israel. He currently serves as a professor at the Homoeopathic Medical College in Nasik and as the Director of Viveda Resort, an innovative holistic health center. He is the former president of the Indian Institute of Homoeopathic Physicians and specializes in reversing autoimmune disorders.

Website: [www.asherclinic.com](http://www.asherclinic.com)
Instagram: @asherhomoeopathy, @doctor.ashar

Check out videos on [our channel](https://www.youtube.com/@Homeosapiens-p7z) featuring Dr. Asher:

*   [_How German New Medicine and Homeopathy Work Together | Dr Asher Shaikh's Holistic Approach_](https://www.youtube.com/watch?v=K8WJlg_zP38)

*   [_Dr Asher Shaikh on German New Medicine, infertility and Spongia_](https://www.youtube.com/watch?v=J4NRGCOdme8)

*   [_Once you hit the bull's-eye… | Kiedy trafisz w dziesiątkę…_](https://www.youtube.com/watch?v=R9l7CSOMRe4)$en$,
  $pl$Zapraszamy na seminarium Homeo sapiens z udziałem dr Ashera Shaikha. Podczas spotkania dr Asher Shaikh przedstawi praktyczne zastosowanie Nowej Germańskiej Medycyny (GNM) w pracy klinicznej homeopaty, ze szczególnym uwzględnieniem procesu prowadzenia wywiadu, analizy konfliktów biologicznych oraz doboru leków homeopatycznych. Omawiane zagadnienia obejmą zarówno część teoretyczną, jak i studia przypadków. **Więcej informacji już wkrótce!**

Seminarium odbędzie się w dniach **5–6 czerwca 2026 r.**

Miejsce wydarzenia: **Hotel Marina Club,** ul. Szafarnia 10, 80-753 Gdańsk oraz **online** (Zoom).

Cena uczestnictwa:

*   **EARLY BIRD 480 PLN / 114 EUR** – do 09.02.2026

*   **560 PLN / 135 EUR** – do 31.03.2026

*   **640 PLN / 152 EUR** – po 31.03.2026

*   **700 PLN / 170 EUR** – w dniu seminarium


Dla uczestników dostępna jest zniżka na noclegi w **Hotelu Marina Club** ([https://marinaclubhotel.pl](https://marinaclubhotel.pl)) oraz **Hostelu Szafarnia** ([https://szafarnia10.pl](https://szafarnia10.pl)) . Rezerwacje realizowane są wyłącznie online. Kod rabatowy: **Homeopatia** (kod aktywny po uruchomieniu w systemie hotelu).

**Dr Asher Shaikh (Indie)** – lekarz homeopata z ponad 30-letnim doświadczeniem klinicznym. Dyrektor Asher Clinics – sieci 12 klinik w Mumbaju, Pune, Dubaju i Nasiku. Mentor Nowej Germańskiej Medycyny, którą wykładał w Dubaju, Indiach, Austrii i Izraelu. Profesor Homoeopathic Medical College w Nasiku. Dyrektor Viveda Resort – ośrodka zdrowia holistycznego. Były przewodniczący Indian Institute of Homoeopathic Physicians. Specjalizuje się w odwracaniu chorób autoimmunologicznych.

Strona internetowa: [www.asherclinic.com](http://www.asherclinic.com)
Instagram: @asherhomoeopathy, @doctor.ashar

Obejrzyj wywiady z Dr Asherem na [naszym kanale](https://www.youtube.com/@Homeosapiens-p7z):

*   [_How German New Medicine and Homeopathy Work Together | Dr Asher Shaikh's Holistic Approach_](https://www.youtube.com/watch?v=K8WJlg_zP38)

*   [_Dr Asher Shaikh on German New Medicine, infertility and Spongia_](https://www.youtube.com/watch?v=J4NRGCOdme8)

*   [_Once you hit the bull's-eye… | Kiedy trafisz w dziesiątkę…_](https://www.youtube.com/watch?v=R9l7CSOMRe4)$pl$
);

-- A Series of Critical Cardiac Cases (free webinar)
insert into events (id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual, description_en, description_pl)
values (
  '019b0c80-a410-7728-ab6b-c1eff529dfd1',
  'webinar',
  'A Series of Critical Cardiac Cases',
  'Seria krytycznych problemów kardiologicznych',
  'a-series-of-critical-cardiac-cases',
  '2025-12-13 16:00:00+00', '2025-12-13 17:30:00+00',
  true,
  $en$Dear Homeopathic Friends,

We are happy to invite you to the next Homeo sapiens Academy webinar. Experienced clinician and homeopath Dr. Herman Jeggels from Cape Town, South Africa will discuss homeopathic treatment in advanced circulatory pathology. He will present documented cases of infective endocarditis, complete AV block and heart failure.

The webinar will be hosted on December 13th 10.00am CET (Poland) / 2.30pm IST (India) / 11:00am SAST (South Africa).

The webinar is free of charge. It will be held in English with consecutive translation to Polish.

It will be held on Zoom via our website (you need to register using email address and a password).$en$,
  $pl$Z przyjemnością zapraszamy na kolejny webinar organizowany przez Homeo sapiens. Naszym gościem będzie Dr Herman Jeggels z Cape Town (RPA), doświadczony klinicysta i homeopata, który przedstawi przypadki leczenia homeopatycznego w zaawansowanych chorobach układu krążenia. Omówione zostaną udokumentowane historie leczenia pacjentów z problemami kardiologicznymi, między innymi infekcyjne zapalenie wsierdzia, całkowity blok przedsionkowo-komorowy i niewydolność serca.

Webinar odbędzie się 13 grudnia 2025 o godzinie 10.00 czasu polskiego.

Wykład będzie tłumaczony konsekutywnie na język polski.

Webinar jest bezpłatny. Odbędzie się na platformie Zoom za pośrednictwem naszej strony internetowej. Wymagana jest rejestracja z użyciem adresu email i ustawienie hasła.$pl$
);

-- To Perfect the Art of Homeopathy (paid, hybrid)
insert into events (id, product_id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual,
  venue_name_en, venue_name_pl, venue_city_en, venue_city_pl, venue_country_code, venue_street, venue_postal_code,
  description_en, description_pl)
values (
  '0199c2f2-528b-7e88-96e3-5e5088333a8c',
  '019de49f-d3d5-727b-915d-2cd3671cb72f',
  'seminar',
  'To Perfect the Art of Homeopathy',
  'Udoskonalić kunszt homeopatyczny',
  'to-perfect-the-art-of-homeopathy',
  '2025-05-30 14:00:00+00', '2025-05-31 08:00:00+00',
  true,
  'Vienna House Easy By Wyndham Cracow', 'Vienna House Easy By Wyndham Kraków',
  'Cracow', 'Kraków',
  'PL', 'ul. Przy Rondzie 2', '31-547',
  'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.',
  $pl$Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

30-31 maja 2025, sala wykładowa B Instytutu Ochrony Roślin, ul. Władysława Węgorka 20, 60-318 Poznań.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.

Omówionych zostanie szereg praktycznych problemów klinicznych, różnicowanie leków z grupy Kalium, leki introwertyczne/ekstrawertyczne, prezentacja przypadków klinicznych.$pl$
);

-- To Perfect the Art of Homeopathy 2 (paid, hybrid)
insert into events (id, product_id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual,
  venue_name_en, venue_name_pl, venue_city_en, venue_city_pl, venue_country_code, venue_street, venue_postal_code,
  description_en, description_pl)
values (
  '0199c2fa-7e9d-72f6-ada1-88b5d04d9a58',
  '019de49f-d612-7339-962b-800168a04bbb',
  'seminar',
  'To Perfect the Art of Homeopathy 2',
  'Udoskonalić kunszt homeopatyczny 2',
  'to-perfect-the-art-of-homeopathy-2',
  '2025-10-24 14:00:00+00', '2025-10-26 11:30:00+00',
  true,
  'IOR Hotel', 'Hotel IOR',
  'Poznań', 'Poznań',
  'PL', 'ul. Węgorka 20', '60-318',
  'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.

October 24-25 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków, Poland.

Online mode will also available (through Zoom). The lectures will be held in English with consecutive translation to Polish.',
  $pl$Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

24-25 października 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.$pl$
);

-- What prevents me from moving on? (free webinar)
insert into events (id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual, subtitle_en, subtitle_pl, description_en, description_pl)
values (
  '019bef00-6ef2-7636-9a15-c8cd1e87b997',
  'webinar',
  'What prevents me from moving on?',
  'What prevents me from moving on?',
  'what-prevents-me-from-moving-on',
  '2026-02-08 15:00:00+00', '2026-02-08 16:30:00+00',
  true,
  'Combining German New Medicine and Homeopathy for musculoskeletal problems',
  'Zastosowania Nowej Germańskiej Medycyny i homeopatii w dolegliwościach układu ruchu',
  $en$We kindly invite you to another free Homeo sapiens webinar. Experienced Homeopath Dr Asher Shaikh will share how he uses German New Medicine to facilitate homeopathic case-taking and remedy choice. Several musculoskeletal problems will be discussed, both theory and case-studies. The webinar will be held in English with consecutive translation to Polish.

Dr Asher Shaikh (India) – a homeopathic doctor with over 25 years of clinical experience. He is the Director of Asher Clinics – a network of 12 clinics in Mumbai, Pune, Dubai, and Nasik – and a mentor in German New Medicine, which he has taught in Dubai, India, Austria, and Israel. He currently serves as a professor at the Homoeopathic Medical College in Nasik and as the Director of Viveda Resort – an innovative holistic health center. He is the former president of the Indian Institute of Homoeopathic Physicians. Dr. Shaikh specializes in reversing autoimmune disorders.$en$,
  $pl$Zapraszamy na kolejny darmowy webinar Homeo sapiens. Doświadczony homeopata dr Asher Shaikh opowie o sposobie, w jaki zastosowanie Nowej Germańskiej Medycyny (GNM) wspomaga przy homeopatycznym doborze leków. Podstawą do dyskusji na ten temat będzie omówienie kilku problemów układu mięśniowo-szkieletowego, zarówno teoretycznie, jak i w oparciu o studia przypadków. Webinar będzie prowadzony w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dr Asher Shaikh (Indie) - lekarz homeopata z ponad 25-letnim doświadczeniem klinicznym. Jest dyrektorem Asher Clinics - 12 klinik w Mumbaju, Pune, Dubaju i Nasiku oraz mentorem Nowej Germańskiej Medycyny, którą wykładał w Dubaju, Indiach, Austrii i Izraelu. Pełni funkcję profesora w Homoeopathic Medical College w Nasiku oraz dyrektora Viveda Resort – innowacyjnego ośrodka zdrowia holistycznego. Były przewodniczący Indian Institute of Homoeopathic Physicians. Specjalizuje się w odwracaniu chorób autoimmunologicznych.$pl$
);

insert into events_hosts (event_id, host_id, position) values
  ('0199c2f2-528b-7e88-96e3-5e5088333a8c', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0),
  ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0),
  ('019b0c80-a410-7728-ab6b-c1eff529dfd1', '019b0c71-fde2-76b7-8c71-21c2e9ea23a5', 0),
  ('019bef00-6ef2-7636-9a15-c8cd1e87b997', '019beef9-4287-714f-982b-2524fdef7063', 0),
  ('019c5c9a-c5a4-7518-8317-65ae90516726', '019beef9-4287-714f-982b-2524fdef7063', 0);

insert into product_prices (product_id, price_type, rule_type, price_amount, price_currency, priority, is_active, valid_until) values
  ('019de49f-d612-7339-962b-800168a04bbb', 'fixed', 'early_bird',    560.00, 'PLN', 10, true, '2025-09-20 21:59:59+00');

insert into product_prices (product_id, price_type, rule_type, price_amount, price_currency, priority, is_active, discount_code) values
  ('019de49f-d612-7339-962b-800168a04bbb', 'fixed', 'discount_code', 500.00, 'PLN', 20, true, 'wshlif');

insert into video_groups (id, title_en, title_pl, slug) values
  ('019da123-449c-7038-aae3-303255746cc4', 'Dr Sanjay Modi: To Perfect the Art of Homeopathy',   'Udoskonalić kunszt homeopatyczny: Seminarium z drem Sanjayem Modim',   'dr-sanjay-modi-to-perfect-the-art-of-homeopathy'),
  ('019daf95-04e3-7615-99fc-ba808d1dd589', 'Dr Asher Shaikh Webinar',                             'Webinarium z drem Asherem Shaikh',                                     'dr-asher-shaikh-webinar'),
  ('019daf95-d855-748d-93a9-4c17d0536f2f', 'Dr Sanjay Modi Webinar',                              'Webinarium z drem Sanjayem Modim',                                     'dr-sanjay-modi-webinar'),
  ('019daf9b-7234-71bb-be93-f9f965d56ac6', 'Dr Sanjay Modi: To Perfect the Art of Homeopathy 2', 'Udoskonalić kunszt homeopatyczny 2: Seminarium z drem Sanjayem Modim', 'dr-sanjay-modi-to-perfect-the-art-of-homeopathy-2'),
  ('019dc005-f4a8-76fb-afdd-2e5caff8fb5a', 'Dr Herman Jeggels Webinar',                           'Webinarium z drem Hermanem Jeggelsem',                                 'dr-herman-jeggels-webinar');

insert into videos (id, provider, title_en, title_pl, slug, recorded_on, host_id) values
  ('019dbfeb-e6f2-7521-b990-119d82b8665f', 'cloudfront', 'To Perfect the Art of Homeopathy: Day 1',          'Udoskonalić kunszt homeopatyczny: Dzień 1',                         'to-perfect-the-art-of-homeopathy-day-1',          '2025-10-24', '0199c2f2-528b-7e88-96e3-5e5088333a8a'),
  ('019dbfeb-e5ec-73ae-881a-d76c8582644e', 'cloudfront', 'To Perfect the Art of Homeopathy: Day 2',          'Udoskonalić kunszt homeopatyczny: Dzień 2',                         'to-perfect-the-art-of-homeopathy-day-2',          '2025-10-25', '0199c2f2-528b-7e88-96e3-5e5088333a8a'),
  ('019a8668-bb4f-7c9c-b9b8-3f274de96566', 'cloudfront', 'To Perfect the Art of Homeopathy 2: Day 1, Part 1','Udoskonalić kunszt homeopatyczny 2: Dzień 1, Część 1',              'to-perfect-the-art-of-homeopathy-2-day-1-part-1', '2025-05-31', '0199c2f2-528b-7e88-96e3-5e5088333a8a'),
  ('019a8ba5-fe29-7af8-bf54-b8d96af38461', 'cloudfront', 'To Perfect the Art of Homeopathy 2: Day 1, Part 2','Udoskonalić kunszt homeopatyczny 2: Dzień 1, Część 2',              'to-perfect-the-art-of-homeopathy-2-day-1-part-2', '2025-05-30', '0199c2f2-528b-7e88-96e3-5e5088333a8a'),
  ('019dbfeb-e512-740f-80ea-d8c30a99fa5b', 'cloudfront', 'Dutiful Remedies: Differential Diagnosis',         'Sumienne leki: Diagnostyka różnicowa',                              'sanjay-modi-dutiful-remedies',                    '2025-03-22', '0199c2f2-528b-7e88-96e3-5e5088333a8a'),
  ('019dbfeb-e43a-7324-bb52-65457afc331b', 'cloudfront', 'What prevents me from moving on?',                 'What prevents me from moving on?',                                  'asher-shaikh-what-prevents-me-from-moving-on',   '2026-02-08', '019beef9-4287-714f-982b-2524fdef7063'),
  ('019dbfec-770a-702f-aa5c-e2431a930395', 'cloudfront', 'A Series of Critical Cardiac Cases',               'Seria krytycznych przypadków kardiologicznych',                     'jeggels-critical-cardiac-cases',                  '2025-12-13', '019b0c71-fde2-76b7-8c71-21c2e9ea23a5');

insert into video_groups_videos (video_id, video_group_id, position) values
  ('019dbfeb-e6f2-7521-b990-119d82b8665f', '019da123-449c-7038-aae3-303255746cc4', 0),
  ('019dbfeb-e5ec-73ae-881a-d76c8582644e', '019da123-449c-7038-aae3-303255746cc4', 1),
  ('019a8668-bb4f-7c9c-b9b8-3f274de96566', '019daf9b-7234-71bb-be93-f9f965d56ac6', 0),
  ('019a8ba5-fe29-7af8-bf54-b8d96af38461', '019daf9b-7234-71bb-be93-f9f965d56ac6', 1),
  ('019dbfeb-e512-740f-80ea-d8c30a99fa5b', '019daf95-d855-748d-93a9-4c17d0536f2f', 0),
  ('019dbfeb-e43a-7324-bb52-65457afc331b', '019daf95-04e3-7615-99fc-ba808d1dd589', 0),
  ('019dbfec-770a-702f-aa5c-e2431a930395', '019dc005-f4a8-76fb-afdd-2e5caff8fb5a', 0);

insert into video_sources (id, video_id, content_type, codec, object_key, priority) values
  -- To Perfect 2 Day 1 Part 1
  ('019dc00d-97b0-743c-9146-f214dacc65f7', '019a8668-bb4f-7c9c-b9b8-3f274de96566', 'application/vnd.apple.mpegurl', null,        '/videos/019a8668-bb4f-7c9c-b9b8-3f274de96566/hls/p1_hls.m3u8',                    0),
  ('019a8ba6-c5ae-7f6f-becb-94b6957a52b2', '019a8668-bb4f-7c9c-b9b8-3f274de96566', 'video/mp4',                    'hev1',       '/videos/019a8668-bb4f-7c9c-b9b8-3f274de96566/hevc_1080.mp4',                       1),
  ('019a8ba7-d04b-77ec-92c6-f76b6ec0e7ea', '019a8668-bb4f-7c9c-b9b8-3f274de96566', 'video/webm',                   'vp9,opus',   '/videos/019a8668-bb4f-7c9c-b9b8-3f274de96566/webm_1080.webm',                      2),
  -- To Perfect 2 Day 1 Part 2
  ('019a8bab-135e-7321-9857-f74d2dcda427', '019a8ba5-fe29-7af8-bf54-b8d96af38461', 'video/mp4',                    'hev1',       '/videos/019a8ba5-fe29-7af8-bf54-b8d96af38461/hevc_1080.mp4',                       0),
  ('019a8bab-bc67-76f9-bf80-902043c922e6', '019a8ba5-fe29-7af8-bf54-b8d96af38461', 'video/webm',                   'vp9,opus',   '/videos/019a8ba5-fe29-7af8-bf54-b8d96af38461/webm_1080.webm',                      1),
  -- To Perfect 1 Day 1
  ('019dc005-f553-7766-aaaa-b48b30707c22', '019dbfeb-e6f2-7521-b990-119d82b8665f', 'video/mp4',                    'avc1.640028,mp4a.40.2', '/videos/019dbfeb-e6f2-7521-b990-119d82b8665f/avc1_1080.mp4',             0),
  ('019dc00b-eb9d-749e-98bf-5ff4b4b76716', '019dbfeb-e6f2-7521-b990-119d82b8665f', 'application/vnd.apple.mpegurl', null,        '/videos/019dbfeb-e6f2-7521-b990-119d82b8665f/hls/index.m3u8',                     1),
  -- To Perfect 1 Day 2
  ('019dc005-f5f7-773f-8e4d-a73fbae12631', '019dbfeb-e5ec-73ae-881a-d76c8582644e', 'video/mp4',                    'avc1.640028,mp4a.40.2', '/videos/019dbfeb-e5ec-73ae-881a-d76c8582644e/avc1_1080.mp4',             0),
  ('019dc00b-ec50-7460-9633-c09016979e48', '019dbfeb-e5ec-73ae-881a-d76c8582644e', 'application/vnd.apple.mpegurl', null,        '/videos/019dbfeb-e5ec-73ae-881a-d76c8582644e/hls/index.m3u8',                     1),
  -- Dr Sanjay Modi webinar (720p)
  ('019dc005-f68e-73db-8700-ec654ffe675c', '019dbfeb-e512-740f-80ea-d8c30a99fa5b', 'video/mp4',                    'avc1.64001f,mp4a.40.2', '/videos/019dbfeb-e512-740f-80ea-d8c30a99fa5b/avc1_720.mp4',              0),
  ('019dc00b-ed09-759f-b756-3968e8e49531', '019dbfeb-e512-740f-80ea-d8c30a99fa5b', 'application/vnd.apple.mpegurl', null,        '/videos/019dbfeb-e512-740f-80ea-d8c30a99fa5b/hls/index.m3u8',                     1),
  -- Dr Asher Shaikh webinar
  ('019dc005-f735-765f-a735-f1c66e74e860', '019dbfeb-e43a-7324-bb52-65457afc331b', 'video/mp4',                    'avc1.640028,mp4a.40.2', '/videos/019dbfeb-e43a-7324-bb52-65457afc331b/avc1_1080.mp4',             0),
  ('019dc00b-edba-70ef-80b9-12b5fdf8f35d', '019dbfeb-e43a-7324-bb52-65457afc331b', 'application/vnd.apple.mpegurl', null,        '/videos/019dbfeb-e43a-7324-bb52-65457afc331b/hls/index.m3u8',                     1),
  -- Dr Herman Jeggels webinar
  ('019dc005-f7df-7661-8c0c-af1a4c4a8f33', '019dbfec-770a-702f-aa5c-e2431a930395', 'video/mp4',                    'avc1.640028,mp4a.40.2', '/videos/019dbfec-770a-702f-aa5c-e2431a930395/avc1_1080.mp4',             0),
  ('019dc00b-ee64-77ce-b15e-2e03b8c573b2', '019dbfec-770a-702f-aa5c-e2431a930395', 'application/vnd.apple.mpegurl', null,        '/videos/019dbfec-770a-702f-aa5c-e2431a930395/hls/index.m3u8',                     1);

commit;
