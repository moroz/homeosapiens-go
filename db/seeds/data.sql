--
-- PostgreSQL database dump
--

\restrict dkaEVfQ7aJPVCLz83Uc5hWX8hOUx5ZbRjVRN1r0kyITh5j1L3bHltva1fC7C6UF

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg12+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg12+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: assets; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.assets VALUES ('0199c2f2-528b-7e88-96e3-5e5088333a8b', 'cm7uqj3q500mglz8z2dqy8sdz.webp', 'cm7uqj3q500mglz8z2dqy8sdz.webp', '2026-05-02 11:47:25', '2026-05-02 11:47:25', false);
INSERT INTO public.assets VALUES ('019b0c7c-c3c4-71c3-a630-7b33a847ca2a', '019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg', '019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg', '2026-05-02 11:47:25', '2026-05-02 11:47:25', false);
INSERT INTO public.assets VALUES ('019beef9-ad4c-736f-9bb0-965b59ca21ae', '019beef9-ad4c-736f-9bb0-965b59ca21ae.png', 'drasher.png', '2026-05-02 11:47:25', '2026-05-02 11:47:25', false);
INSERT INTO public.assets VALUES ('019de856-e5a2-7edb-9bba-215d32de9250', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e4bd-799d-a689-8d4d5cf7370b', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e6ef-73f3-88f6-9b1772fdf073', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e650-7d54-a6ba-ff70c8e812fb', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e865-79a5-829a-45cad21e4e34', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e799-70b2-85a6-0c36770b45e3', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e903-7939-ad2b-8a3f0abc43b1', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-e9bb-7ad6-81ea-04c12e24768b', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-ea73-7a5b-9798-a0b68193ebe5', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-eb1b-779d-85a1-95e040bfc733', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-ec6e-76f6-a5cc-10508960adc8', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-ebc5-79ee-ab25-a7482adebcb0', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-ed0c-7bb5-931c-4061cf50e25d', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);
INSERT INTO public.assets VALUES ('019de856-edd2-7494-9d7e-0f69f3629301', NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', true);


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.products VALUES ('019de49f-d17f-7435-b166-b8e9b3e4430c', 'event', 'Seminarium z drem Asherem Shaikh', 'Dr Asher Shaikh seminar', 560.00, 'PLN', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.products VALUES ('019de49f-d3d5-727b-915d-2cd3671cb72f', 'event', 'Udoskonalić kunszt homeopatyczny', 'To Perfect the Art of Homeopathy', 580.00, 'PLN', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.products VALUES ('019de49f-d612-7339-962b-800168a04bbb', 'event', 'Udoskonalić kunszt homeopatyczny 2', 'To Perfect the Art of Homeopathy 2', 640.00, 'PLN', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.events VALUES ('019c5c9a-c5a4-7518-8317-65ae90516726', 'Dr Asher Shaikh seminar', 'Seminarium z drem Asherem Shaikh', '2026-06-05 14:00:00', '2026-06-06 20:00:00', false, 'We kindly invite you to the Homeo sapiens seminar with Dr Asher Shaikh. During the seminar, Dr Shaikh will present the practical application of German New Medicine (GNM) in clinical homeopathic practice, with particular focus on case-taking, identification of biological conflicts, and remedy selection. The seminar will include both theoretical background and case studies. **More info coming soon!**

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

*   [_How German New Medicine and Homeopathy Work Together | Dr Asher Shaikh''s Holistic Approach_](https://www.youtube.com/watch?v=K8WJlg_zP38)

*   [_Dr Asher Shaikh on German New Medicine, infertility and Spongia_](https://www.youtube.com/watch?v=J4NRGCOdme8)

*   [_Once you hit the bull''s-eye… | Kiedy trafisz w dziesiątkę…_](https://www.youtube.com/watch?v=R9l7CSOMRe4)', 'Zapraszamy na seminarium Homeo sapiens z udziałem dr Ashera Shaikha. Podczas spotkania dr Asher Shaikh przedstawi praktyczne zastosowanie Nowej Germańskiej Medycyny (GNM) w pracy klinicznej homeopaty, ze szczególnym uwzględnieniem procesu prowadzenia wywiadu, analizy konfliktów biologicznych oraz doboru leków homeopatycznych. Omawiane zagadnienia obejmą zarówno część teoretyczną, jak i studia przypadków. **Więcej informacji już wkrótce!**

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

*   [_How German New Medicine and Homeopathy Work Together | Dr Asher Shaikh''s Holistic Approach_](https://www.youtube.com/watch?v=K8WJlg_zP38)

*   [_Dr Asher Shaikh on German New Medicine, infertility and Spongia_](https://www.youtube.com/watch?v=J4NRGCOdme8)

*   [_Once you hit the bull''s-eye… | Kiedy trafisz w dziesiątkę…_](https://www.youtube.com/watch?v=R9l7CSOMRe4)', 'seminar', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'dr-asher-shaikh-seminar', NULL, NULL, 'Marina Club Hotel', 'Hotel Marina Club', 'ul. Szafarnia 10', 'Gdańsk', 'Gdańsk', '80-753', 'PL', '019de49f-d17f-7435-b166-b8e9b3e4430c');
INSERT INTO public.events VALUES ('019b0c80-a410-7728-ab6b-c1eff529dfd1', 'A Series of Critical Cardiac Cases', 'Seria krytycznych problemów kardiologicznych', '2025-12-13 16:00:00', '2025-12-13 17:30:00', true, 'Dear Homeopathic Friends,

We are happy to invite you to the next Homeo sapiens Academy webinar. Experienced clinician and homeopath Dr. Herman Jeggels from Cape Town, South Africa will discuss homeopathic treatment in advanced circulatory pathology. He will present documented cases of infective endocarditis, complete AV block and heart failure.

The webinar will be hosted on December 13th 10.00am CET (Poland) / 2.30pm IST (India) / 11:00am SAST (South Africa).

The webinar is free of charge. It will be held in English with consecutive translation to Polish.

It will be held on Zoom via our website (you need to register using email address and a password).', 'Z przyjemnością zapraszamy na kolejny webinar organizowany przez Homeo sapiens. Naszym gościem będzie Dr Herman Jeggels z Cape Town (RPA), doświadczony klinicysta i homeopata, który przedstawi przypadki leczenia homeopatycznego w zaawansowanych chorobach układu krążenia. Omówione zostaną udokumentowane historie leczenia pacjentów z problemami kardiologicznymi, między innymi infekcyjne zapalenie wsierdzia, całkowity blok przedsionkowo-komorowy i niewydolność serca.

Webinar odbędzie się 13 grudnia 2025 o godzinie 10.00 czasu polskiego.

Wykład będzie tłumaczony konsekutywnie na język polski.

Webinar jest bezpłatny. Odbędzie się na platformie Zoom za pośrednictwem naszej strony internetowej. Wymagana jest rejestracja z użyciem adresu email i ustawienie hasła.', 'webinar', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'a-series-of-critical-cardiac-cases', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);
INSERT INTO public.events VALUES ('0199c2f2-528b-7e88-96e3-5e5088333a8c', 'To Perfect the Art of Homeopathy', 'Udoskonalić kunszt homeopatyczny', '2025-05-30 14:00:00', '2025-05-31 08:00:00', true, 'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.', 'Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

30-31 maja 2025, sala wykładowa B Instytutu Ochrony Roślin, ul. Władysława Węgorka 20, 60-318 Poznań.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.

Omówionych zostanie szereg praktycznych problemów klinicznych, różnicowanie leków z grupy Kalium, leki introwertyczne/ekstrawertyczne, prezentacja przypadków klinicznych.', 'seminar', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'to-perfect-the-art-of-homeopathy', NULL, NULL, 'Vienna House Easy By Wyndham Cracow', 'Vienna House Easy By Wyndham Kraków', 'ul. Przy Rondzie 2', 'Cracow', 'Kraków', '31-547', 'PL', '019de49f-d3d5-727b-915d-2cd3671cb72f');
INSERT INTO public.events VALUES ('0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', 'To Perfect the Art of Homeopathy 2', 'Udoskonalić kunszt homeopatyczny 2', '2025-10-24 14:00:00', '2025-10-26 11:30:00', true, 'Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.

October 24-25 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków, Poland.

Online mode will also available (through Zoom). The lectures will be held in English with consecutive translation to Polish.', 'Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

24-25 października 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.', 'seminar', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'to-perfect-the-art-of-homeopathy-2', NULL, NULL, 'IOR Hotel', 'Hotel IOR', 'ul. Węgorka 20', 'Poznań', 'Poznań', '60-318', 'PL', '019de49f-d612-7339-962b-800168a04bbb');
INSERT INTO public.events VALUES ('019bef00-6ef2-7636-9a15-c8cd1e87b997', 'What prevents me from moving on?', 'What prevents me from moving on?', '2026-02-08 15:00:00', '2026-02-08 16:30:00', true, 'We kindly invite you to another free Homeo sapiens webinar. Experienced Homeopath Dr Asher Shaikh will share how he uses German New Medicine to facilitate homeopathic case-taking and remedy choice. Several musculoskeletal problems will be discussed, both theory and case-studies. The webinar will be held in English with consecutive translation to Polish.

Dr Asher Shaikh (India) – a homeopathic doctor with over 25 years of clinical experience. He is the Director of Asher Clinics – a network of 12 clinics in Mumbai, Pune, Dubai, and Nasik – and a mentor in German New Medicine, which he has taught in Dubai, India, Austria, and Israel. He currently serves as a professor at the Homoeopathic Medical College in Nasik and as the Director of Viveda Resort – an innovative holistic health center. He is the former president of the Indian Institute of Homoeopathic Physicians. Dr. Shaikh specializes in reversing autoimmune disorders.', 'Zapraszamy na kolejny darmowy webinar Homeo sapiens. Doświadczony homeopata dr Asher Shaikh opowie o sposobie, w jaki zastosowanie Nowej Germańskiej Medycyny (GNM) wspomaga przy homeopatycznym doborze leków. Podstawą do dyskusji na ten temat będzie omówienie kilku problemów układu mięśniowo-szkieletowego, zarówno teoretycznie, jak i w oparciu o studia przypadków. Webinar będzie prowadzony w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dr Asher Shaikh (Indie) - lekarz homeopata z ponad 25-letnim doświadczeniem klinicznym. Jest dyrektorem Asher Clinics - 12 klinik w Mumbaju, Pune, Dubaju i Nasiku oraz mentorem Nowej Germańskiej Medycyny, którą wykładał w Dubaju, Indiach, Austrii i Izraelu. Pełni funkcję profesora w Homoeopathic Medical College w Nasiku oraz dyrektora Viveda Resort – innowacyjnego ośrodka zdrowia holistycznego. Były przewodniczący Indian Institute of Homoeopathic Physicians. Specjalizuje się w odwracaniu chorób autoimmunologicznych.', 'webinar', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'what-prevents-me-from-moving-on', 'Combining German New Medicine and Homeopathy for musculoskeletal problems', 'Zastosowania Nowej Germańskiej Medycyny i homeopatii w dolegliwościach układu ruchu', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL);


--
-- Data for Name: hosts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.hosts VALUES ('0199c2f2-528b-7e88-96e3-5e5088333a8a', 'common.hosts.salutation.dr', 'Sanjay', 'Modi', '0199c2f2-528b-7e88-96e3-5e5088333a8b', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'IN');
INSERT INTO public.hosts VALUES ('019beef9-4287-714f-982b-2524fdef7063', 'common.hosts.salutation.dr', 'Asher', 'Shaikh', '019beef9-ad4c-736f-9bb0-965b59ca21ae', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'IN');
INSERT INTO public.hosts VALUES ('019b0c71-fde2-76b7-8c71-21c2e9ea23a5', 'common.hosts.salutation.dr', 'Herman', 'Jeggels', '019b0c7c-c3c4-71c3-a630-7b33a847ca2a', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 'ZA');


--
-- Data for Name: events_hosts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.events_hosts VALUES ('019de883-cd45-7364-9729-32ce9cf2b54b', '0199c2f2-528b-7e88-96e3-5e5088333a8c', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0, '2026-05-02 11:47:25', '2026-05-02 11:47:25');
INSERT INTO public.events_hosts VALUES ('019de883-cd45-753b-824d-815911cb636b', '0199c2fa-7e9d-72f6-ada1-88b5d04d9a58', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 0, '2026-05-02 11:47:25', '2026-05-02 11:47:25');
INSERT INTO public.events_hosts VALUES ('019de883-cd45-7558-86a2-698af86afece', '019b0c80-a410-7728-ab6b-c1eff529dfd1', '019b0c71-fde2-76b7-8c71-21c2e9ea23a5', 0, '2026-05-02 11:47:25', '2026-05-02 11:47:25');
INSERT INTO public.events_hosts VALUES ('019de883-cd45-7567-83d3-893966ec5df9', '019bef00-6ef2-7636-9a15-c8cd1e87b997', '019beef9-4287-714f-982b-2524fdef7063', 0, '2026-05-02 11:47:25', '2026-05-02 11:47:25');
INSERT INTO public.events_hosts VALUES ('019de883-cd45-7574-97da-c8086bb3e95d', '019c5c9a-c5a4-7518-8317-65ae90516726', '019beef9-4287-714f-982b-2524fdef7063', 0, '2026-05-02 11:47:25', '2026-05-02 11:47:25');


--
-- Data for Name: product_prices; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_prices VALUES ('019de883-cd45-7ac8-b061-58a2acd390d8', 'Fixed', 'EarlyBird', 560.00000000, 'PLN', NULL, 10, true, NULL, '2025-09-20 21:59:59', '2026-05-02 11:47:25', '2026-05-02 11:47:25', '019de49f-d612-7339-962b-800168a04bbb');
INSERT INTO public.product_prices VALUES ('019de883-cd45-7e67-9698-ba121d093365', 'Fixed', 'DiscountCode', 500.00000000, 'PLN', 'wshlif', 20, true, NULL, NULL, '2026-05-02 11:47:25', '2026-05-02 11:47:25', '019de49f-d612-7339-962b-800168a04bbb');


--
-- Data for Name: video_groups; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.video_groups VALUES ('019da123-449c-7038-aae3-303255746cc4', 'Dr Sanjay Modi: To Perfect the Art of Homeopathy', 'Udoskonalić kunszt homeopatyczny: Seminarium z drem Sanjayem Modim', 'dr-sanjay-modi-to-perfect-the-art-of-homeopathy', NULL, '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups VALUES ('019daf95-04e3-7615-99fc-ba808d1dd589', 'Dr Asher Shaikh Webinar', 'Webinarium z drem Asherem Shaikh', 'dr-asher-shaikh-webinar', NULL, '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups VALUES ('019daf95-d855-748d-93a9-4c17d0536f2f', 'Dr Sanjay Modi Webinar', 'Webinarium z drem Sanjayem Modim', 'dr-sanjay-modi-webinar', NULL, '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups VALUES ('019daf9b-7234-71bb-be93-f9f965d56ac6', 'Dr Sanjay Modi: To Perfect the Art of Homeopathy 2', 'Udoskonalić kunszt homeopatyczny 2: Seminarium z drem Sanjayem Modim', 'dr-sanjay-modi-to-perfect-the-art-of-homeopathy-2', NULL, '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups VALUES ('019dc005-f4a8-76fb-afdd-2e5caff8fb5a', 'Dr Herman Jeggels Webinar', 'Webinarium z drem Hermanem Jeggelsem', 'dr-herman-jeggels-webinar', NULL, '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');


--
-- Data for Name: videos; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.videos VALUES ('019dbfeb-e6f2-7521-b990-119d82b8665f', 'cloudfront', false, 'To Perfect the Art of Homeopathy: Day 1', 'Udoskonalić kunszt homeopatyczny: Dzień 1', 'to-perfect-the-art-of-homeopathy-day-1', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 18107, '2025-10-24', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '019de856-ec6e-76f6-a5cc-10508960adc8', '019de856-ebc5-79ee-ab25-a7482adebcb0');
INSERT INTO public.videos VALUES ('019dbfeb-e5ec-73ae-881a-d76c8582644e', 'cloudfront', false, 'To Perfect the Art of Homeopathy: Day 2', 'Udoskonalić kunszt homeopatyczny: Dzień 2', 'to-perfect-the-art-of-homeopathy-day-2', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 18470, '2025-10-25', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '019de856-ea73-7a5b-9798-a0b68193ebe5', '019de856-eb1b-779d-85a1-95e040bfc733');
INSERT INTO public.videos VALUES ('019a8668-bb4f-7c9c-b9b8-3f274de96566', 'cloudfront', false, 'To Perfect the Art of Homeopathy 2: Day 1, Part 1', 'Udoskonalić kunszt homeopatyczny 2: Dzień 1, Część 1', 'to-perfect-the-art-of-homeopathy-2-day-1-part-1', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 8344, '2025-05-31', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '019de856-e5a2-7edb-9bba-215d32de9250', '019de856-e4bd-799d-a689-8d4d5cf7370b');
INSERT INTO public.videos VALUES ('019a8ba5-fe29-7af8-bf54-b8d96af38461', 'cloudfront', false, 'To Perfect the Art of Homeopathy 2: Day 1, Part 2', 'Udoskonalić kunszt homeopatyczny 2: Dzień 1, Część 2', 'to-perfect-the-art-of-homeopathy-2-day-1-part-2', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 7097, '2025-05-30', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '019de856-e6ef-73f3-88f6-9b1772fdf073', '019de856-e650-7d54-a6ba-ff70c8e812fb');
INSERT INTO public.videos VALUES ('019dbfeb-e512-740f-80ea-d8c30a99fa5b', 'cloudfront', false, 'Dutiful Remedies: Differential Diagnosis', 'Sumienne leki: Diagnostyka różnicowa', 'sanjay-modi-dutiful-remedies', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 4943, '2025-03-22', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '019de856-e903-7939-ad2b-8a3f0abc43b1', '019de856-e9bb-7ad6-81ea-04c12e24768b');
INSERT INTO public.videos VALUES ('019dbfeb-e43a-7324-bb52-65457afc331b', 'cloudfront', false, 'What prevents me from moving on?', 'What prevents me from moving on?', 'asher-shaikh-what-prevents-me-from-moving-on', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 8389, '2026-02-08', '019beef9-4287-714f-982b-2524fdef7063', '019de856-e865-79a5-829a-45cad21e4e34', '019de856-e799-70b2-85a6-0c36770b45e3');
INSERT INTO public.videos VALUES ('019dbfec-770a-702f-aa5c-e2431a930395', 'cloudfront', false, 'A Series of Critical Cardiac Cases', 'Seria krytycznych przypadków kardiologicznych', 'jeggels-critical-cardiac-cases', '2026-05-02 11:47:25', '2026-05-02 11:47:25', 5523, '2025-12-13', '019b0c71-fde2-76b7-8c71-21c2e9ea23a5', '019de856-ed0c-7bb5-931c-4061cf50e25d', '019de856-edd2-7494-9d7e-0f69f3629301');


--
-- Data for Name: video_groups_videos; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7da6-a99b-9ed5f8dfe413', 0, '019dbfeb-e6f2-7521-b990-119d82b8665f', '019da123-449c-7038-aae3-303255746cc4', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7f97-a00d-1d4b29a911eb', 1, '019dbfeb-e5ec-73ae-881a-d76c8582644e', '019da123-449c-7038-aae3-303255746cc4', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7fb7-b519-956e9979eb5a', 0, '019a8668-bb4f-7c9c-b9b8-3f274de96566', '019daf9b-7234-71bb-be93-f9f965d56ac6', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7fc8-a9a4-7711774a9b95', 1, '019a8ba5-fe29-7af8-bf54-b8d96af38461', '019daf9b-7234-71bb-be93-f9f965d56ac6', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7fe3-b5bd-56a742f5a2fc', 0, '019dbfeb-e512-740f-80ea-d8c30a99fa5b', '019daf95-d855-748d-93a9-4c17d0536f2f', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd46-7ff3-ad87-df0854df1c4e', 0, '019dbfeb-e43a-7324-bb52-65457afc331b', '019daf95-04e3-7615-99fc-ba808d1dd589', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');
INSERT INTO public.video_groups_videos VALUES ('019de883-cd47-7002-9e69-c6ec4ee52c91', 0, '019dbfec-770a-702f-aa5c-e2431a930395', '019dc005-f4a8-76fb-afdd-2e5caff8fb5a', '2026-05-02 11:47:25.116074', '2026-05-02 11:47:25.116074');


--
-- PostgreSQL database dump complete
--

\unrestrict dkaEVfQ7aJPVCLz83Uc5hWX8hOUx5ZbRjVRN1r0kyITh5j1L3bHltva1fC7C6UF

