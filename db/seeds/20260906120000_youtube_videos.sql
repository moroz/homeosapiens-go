-- Videos published on the Homeo sapiens YouTube channel that were missing from
-- the initial seed data, together with the hosts featured in them.
--
-- Standalone and idempotent: identifiers are fixed, every statement is
-- ON CONFLICT DO NOTHING, so re-running the file changes nothing.
--
--   psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f 20260906120000_youtube_videos.sql

begin;

-- Hosts

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7a25-9831-08c44b13dcf2', 'common.hosts.salutation.dr', 'Bruno', 'Galeazzi', 'IT')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7e8b-9b76-ff548969e3ee', 'common.hosts.salutation.dr', 'Vangelis', 'Zafeiriou', 'GR')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-75eb-8846-d8a867303094', NULL, 'Murilo', 'Fontoura', 'BR')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7f6c-ada5-0645fc379542', NULL, 'Kim', 'Elia', 'US')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-707e-98f6-88c43ea68c0b', 'common.hosts.salutation.dr', 'Helmut', 'Roniger', NULL)
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7b0f-82c9-2cac1c245a6b', 'common.hosts.salutation.dr', 'Mini', 'Mehta', 'IN')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7d1b-9081-81a257a45753', 'common.hosts.salutation.dr', 'Alex', 'Bekker', 'US')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7637-9b37-16c3511479aa', 'common.hosts.salutation.dr', 'Rainer', 'Schäferkordt', 'DE')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7c70-b0a8-be926bd9cd0d', 'common.hosts.salutation.prof', 'Aaron', 'To', 'HK')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-7922-b2fc-b257db6032dc', 'common.hosts.salutation.dr', 'Shepherd Roee', 'Singer', 'IL')
on conflict do nothing;

insert into hosts (id, salutation, given_name, family_name, country) values
    ('01a078bc-02d5-76a1-8d0b-2ed2d129ba3f', 'common.hosts.salutation.dr', 'Kishore', 'Mehta', 'IN')
on conflict do nothing;

-- Prof. Michael Frass was seeded with the doctor salutation before a
-- professor salutation existed.
update hosts set salutation = 'common.hosts.salutation.prof', updated_at = now()
where id = '019f75f3-7540-7971-b2e4-a9fca75ab7ab' and salutation is distinct from 'common.hosts.salutation.prof';


-- Videos

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7cd3-8c21-89f49425975e', 'youtube', true,
     'Psychosomatics forever — Dr Bruno Galeazzi on hypertension, alexithymia and the language of symptoms',
     'Psychosomatyka na zawsze — dr Bruno Galeazzi o nadciśnieniu, aleksytymii i języku objawów',
     'bruno-galeazzi-psychosomatics-forever',
     1191, '2026-09-04', '01a078bc-02d5-7a25-9831-08c44b13dcf2', 'vMLgkxqP62o',
     'Italy''s Dr Bruno Galeazzi returns to Homeo sapiens to talk with Dr Barbara Liberowicz about psychosomatics seen through — but not limited to — a homeopathic lens: hypertension, alexithymia (the inability to put one''s emotions into words) and the body''s own "language" of symptoms.',
     'Włoski lekarz dr Bruno Galeazzi ponownie gości w Homeo sapiens i rozmawia z dr Barbarą Liberowicz o psychosomatyce widzianej przez pryzmat homeopatii (choć nie tylko): o nadciśnieniu, aleksytymii, czyli niezdolności do wyrażania emocji słowami, oraz o „języku” objawów ciała.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7ead-9326-20c2e24c6903', 'youtube', true,
     'Psychotherapy or homeopathy? An interview with Dr Vangelis Zafeiriou',
     'Psychoterapia czy homeopatia? Rozmowa z drem Vangelisem Zafeiriou',
     'vangelis-zafeiriou-psychotherapy-or-homeopathy',
     908, '2026-07-25', '01a078bc-02d5-7e8b-9b76-ff548969e3ee', 'A_Q5kHO3Nu4',
     'Greek psychiatrist and homeopath Dr Vangelis Zafeiriou talks about his inquiry into the personalities of Samuel Hahnemann and Sigmund Freud, how the personality of a method''s creator shapes the method itself, and what he has learned from working in both the psychotherapeutic and the homeopathic world.',
     'Grecki psychiatra i homeopata dr Vangelis Zafeiriou opowiada o swoich badaniach nad osobowościami Samuela Hahnemanna i Zygmunta Freuda, o tym, jak osobowość twórcy metody wpływa na samą metodę, oraz o doświadczeniach z obu światów — psychoterapii i homeopatii.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7679-9180-5588223d5a4e', 'youtube', true,
     'Dr Asher Shaikh webinar, part 2: homeopathy and German New Medicine for musculoskeletal problems',
     'Webinarium z drem Asherem Shaikh, część 2: homeopatia i Germańska Nowa Medycyna w dolegliwościach narządu ruchu',
     'asher-shaikh-webinar-musculoskeletal-problems-pt-2',
     949, '2026-05-29', '019beef9-4287-714f-982b-2524fdef7063', 'htuFnEAuP08',
     '"What restricts my movements and my activity? Combining homeopathy and German New Medicine for musculoskeletal problems" — part 2 of Dr Asher Shaikh''s Homeo sapiens webinar of 10 May 2026.',
     '„Co ogranicza moje ruchy i moją aktywność? Połączenie homeopatii i Germańskiej Nowej Medycyny w dolegliwościach narządu ruchu” — część 2 webinarium dra Ashera Shaikha w Homeo sapiens z 10 maja 2026 r.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7f8f-9484-fc9b4de844a4', 'youtube', true,
     'Opioid-addicted rats and ultramolecular dilutions of esketamine — Murilo Fontoura',
     'Szczury uzależnione od opioidów i ultramolekularne rozcieńczenia esketaminy — Murilo Fontoura',
     'murilo-fontoura-esketamine-ultramolecular-dilutions',
     713, '2026-05-21', '01a078bc-02d5-75eb-8846-d8a867303094', 'FWpXMB0hx0k',
     'Murilo Fontoura presents his research on allopathic and ultramolecular doses of esketamine in opioid-dependent rats, published as _Beneficial effects of Esketamine on Morphine preference reacquisition in male rats_.',
     'Murilo Fontoura przedstawia swoje badania nad allopatycznymi i ultramolekularnymi dawkami esketaminy u szczurów uzależnionych od opioidów, opublikowane jako _Beneficial effects of Esketamine on Morphine preference reacquisition in male rats_.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-74d7-a227-544d0fb5fa89', 'youtube', true,
     'Nanostructures, C.G. Jung and homeopathy — Dr Bruno Galeazzi',
     'Nanostruktury, C.G. Jung i homeopatia — dr Bruno Galeazzi',
     'bruno-galeazzi-nanostructures-jung-and-homeopathy',
     906, '2026-05-05', '01a078bc-02d5-7a25-9831-08c44b13dcf2', 'w3FDbvP05eY',
     'Dr Bruno Galeazzi, past president of FIAMO (the Italian federation of homeopathic associations and homeopaths), talks about his research on homeopathy from the perspective of quantum physics, the influence of Carl Gustav Jung''s psychology on his thinking, and his own homeopathic journey.',
     'Dr Bruno Galeazzi, były prezes FIAMO (włoskiej federacji stowarzyszeń homeopatycznych i homeopatów), opowiada o swoich badaniach nad homeopatią z perspektywy fizyki kwantowej, o wpływie psychologii Carla Gustava Junga na jego myślenie oraz o własnej drodze do homeopatii.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7246-af71-1dd0a4b3c8d5', 'youtube', true,
     'Hahnemann goes Hollywood? Kim Elia on his film _Introducing Homeopathy_',
     'Hahnemann w Hollywood? Kim Elia o swoim filmie _Introducing Homeopathy_',
     'kim-elia-introducing-homeopathy-movie',
     742, '2026-04-16', '01a078bc-02d5-7f6c-ada5-0645fc379542', 'y98-pfhiDpI',
     'Barbara Liberowicz talks to homeopath Kim Elia, creator of _Introducing Homeopathy_ — a feature-length documentary exploring the science, the history and the untold potential of one of the most misunderstood healing modalities of our time.',
     'Barbara Liberowicz rozmawia z homeopatą Kimem Elią, twórcą pełnometrażowego dokumentu _Introducing Homeopathy_, opowiadającego o nauce, historii i nieopowiedzianym potencjale jednej z najbardziej niezrozumianych metod leczenia naszych czasów.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7d9e-9f69-6386317ef58b', 'youtube', true,
     'Experts'' point: the connection between German New Medicine and homeopathy clarified',
     'Zdaniem ekspertów: związek Germańskiej Nowej Medycyny z homeopatią wyjaśniony',
     'gnm-connection-to-homeopathy-clarified',
     881, '2026-03-24', '019beef9-4287-714f-982b-2524fdef7063', 'WxiFHfU3CuA',
     'Dr Asher Shaikh and Dr Robert Zawiślak discuss the new diagnostic possibilities that German New Medicine opens up in homeopathic case taking.',
     'Dr Asher Shaikh i dr Robert Zawiślak rozmawiają o nowych możliwościach diagnostycznych, jakie Germańska Nowa Medycyna otwiera w homeopatycznym wywiadzie.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7fe5-a978-9b0d1d1994c7', 'youtube', true,
     'A rather unusual career of a holistic doctor — Dr Helmut Roniger',
     'Nietypowa kariera lekarza holistycznego — dr Helmut Roniger',
     'helmut-roniger-unusual-career-of-a-holistic-doctor',
     1201, '2026-03-17', '01a078bc-02d5-707e-98f6-88c43ea68c0b', 'taIF5g6tGGM',
     'Dr Helmut Roniger worked at the Royal London Homeopathic Hospital for 23 years. He talks about his path through conventional medicine and homeopathy, his approach to inflammatory bowel disease, cancer, lifestyle medicine and other chronic illnesses — and about probably the largest homeoprophylaxis trial ever conducted, in which he played a part.',
     'Dr Helmut Roniger przez 23 lata pracował w Royal London Homeopathic Hospital. Opowiada o swojej drodze przez medycynę konwencjonalną i homeopatię, o podejściu do nieswoistych zapaleń jelit, chorób nowotworowych, medycyny stylu życia i innych chorób przewlekłych, a także o prawdopodobnie największym badaniu homeoprofilaktycznym w historii, w którym brał udział.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7215-992a-65ab92a3b638', 'youtube', true,
     'Dr Asher Shaikh webinar, part 1: psychodynamics and homeopathy for musculoskeletal problems',
     'Webinarium z drem Asherem Shaikh, część 1: psychodynamika i homeopatia w dolegliwościach narządu ruchu',
     'asher-shaikh-webinar-psychodynamics-musculoskeletal-pt-1',
     900, '2026-02-18', '019beef9-4287-714f-982b-2524fdef7063', 'QaVpeaJ9h6I',
     'Part 1 of the Homeo sapiens webinar with Dr Asher Shaikh (8 February 2026), in which he explains how he uses German New Medicine as a tool that supports homeopathic diagnosis and the choice of remedy.',
     'Część 1 webinarium Homeo sapiens z drem Asherem Shaikh (8 lutego 2026 r.), w którym wyjaśnia, jak wykorzystuje Germańską Nową Medycynę jako narzędzie wspierające homeopatyczną diagnozę i dobór leku.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7f90-8c18-c1add5aeeb70', 'youtube', true,
     'Dr Asher Shaikh webinar, part 2: psychodynamics and homeopathy for musculoskeletal problems',
     'Webinarium z drem Asherem Shaikh, część 2: psychodynamika i homeopatia w dolegliwościach narządu ruchu',
     'asher-shaikh-webinar-psychodynamics-musculoskeletal-pt-2',
     1715, '2026-02-18', '019beef9-4287-714f-982b-2524fdef7063', 'Y2kPONvmB-w',
     'Part 2 of the Homeo sapiens webinar with Dr Asher Shaikh (8 February 2026), in which he explains how he uses German New Medicine as a tool that supports homeopathic diagnosis and the choice of remedy.',
     'Część 2 webinarium Homeo sapiens z drem Asherem Shaikh (8 lutego 2026 r.), w którym wyjaśnia, jak wykorzystuje Germańską Nową Medycynę jako narzędzie wspierające homeopatyczną diagnozę i dobór leku.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-752b-af53-c1f0883573fb', 'youtube', true,
     'The study on homeopathy in cancer that was… — an interview with Prof. Michael Frass',
     'Badanie nad homeopatią w chorobie nowotworowej, które jednak było… — rozmowa z prof. Michaelem Frassem',
     'michael-frass-the-study-on-homeopathy-in-cancer',
     661, '2026-01-27', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 'Yirh6RhhnsM',
     'Prof. Michael Frass is the author of a double-blind, placebo-controlled study comparing quality of life and survival in patients with advanced non-small cell lung cancer treated with add-on homeopathy.',
     'Prof. Michael Frass jest autorem badania z podwójnie ślepą próbą i kontrolą placebo, porównującego jakość życia i przeżycie pacjentów z zaawansowanym niedrobnokomórkowym rakiem płuca leczonych dodatkowo homeopatycznie.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-717e-8f5d-53d424af1c63', 'youtube', true,
     'Best of 2025',
     'Najlepsze momenty 2025 roku',
     'best-of-2025',
     611, '2025-12-23', NULL, 'rLGTTX4X9rI',
     'A selection of Homeo sapiens moments from 2025.',
     'Wybór momentów Homeo sapiens z 2025 roku.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-761c-88cc-59337d2cc1d5', 'youtube', true,
     'Dr Mini Mehta on palliative homeopathic cancer care, rejected cases and German New Medicine',
     'Dr Mini Mehta o homeopatycznej opiece paliatywnej w onkologii, przypadkach odrzuconych i Germańskiej Nowej Medycynie',
     'mini-mehta-palliative-cancer-care-and-gnm',
     1106, '2025-12-01', '01a078bc-02d5-7b0f-82c9-2cac1c245a6b', 'g7Iq4cKc7JA',
     'Dr Mini Mehta talks about her path as a homeopathic doctor, her experience of working in an allopathic hospital, and how she uses German New Medicine to help choose the homeopathic remedy. She has practised in New Delhi since 2009 and works in palliative oncology at the Rajiv Gandhi Cancer Hospital.',
     'Dr Mini Mehta opowiada o swojej drodze jako lekarki homeopatki, o pracy w szpitalu allopatycznym i o tym, jak wykorzystuje Germańską Nową Medycynę przy doborze leku homeopatycznego. Praktykuje w Nowym Delhi od 2009 roku i zajmuje się onkologią paliatywną w Rajiv Gandhi Cancer Hospital.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7e80-bb19-f91f2c5d0418', 'youtube', true,
     'Dr V on _Chelidonium_, black bile and swallows',
     'Dr V o leku _Chelidonium_, czarnej żółci i jaskółkach',
     'dr-v-on-chelidonium-black-bile-and-swallows',
     1175, '2025-11-22', '019f75f3-7540-7aac-99de-4ad4bc37b9b4', 'I8c6uF4pwv8',
     'Dr V (Dr Shailendra Vaishampayan) discusses a remarkable but little-known homeopathic remedy that deserves far more attention: _Chelidonium_.',
     'Dr V (dr Shailendra Vaishampayan) omawia niezwykły, choć mało znany lek homeopatyczny, który zasługuje na znacznie więcej uwagi — _Chelidonium_.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-76bf-839f-fb4a0c978a57', 'youtube', true,
     'Desire for light? Dr Modi and Dr Zawiślak talk in Kraków',
     'Pragnienie światła? Dr Modi i dr Zawiślak rozmawiają w Krakowie',
     'modi-and-zawislak-desire-for-light-krakow',
     389, '2025-11-17', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'lqpsJo097fk',
     'Dr Sanjay Modi and Dr Robert Zawiślak talk about the magnesium salts and much more, just before Dr Modi''s Kraków seminar.',
     'Dr Sanjay Modi i dr Robert Zawiślak rozmawiają o solach magnezu i nie tylko, tuż przed krakowskim seminarium dra Modiego.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-79b3-a127-e81f9776c905', 'youtube', true,
     'Dr Sanjay Modi webinar, part 3: absent-mindedness',
     'Webinarium z drem Sanjayem Modim, część 3: roztargnienie',
     'sanjay-modi-webinar-absent-mindedness-pt-3',
     793, '2025-10-20', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'LlEnLoqdTOM',
     'Part 3 of Dr Modi''s webinar _Boldness and daring vs terror and restlessness — qualities viewed from the point of view of homeopathic remedy differentiation_, recorded on 13 September 2025.',
     'Część 3 webinarium dra Modiego _Śmiałość i odwaga a przerażenie i niepokój — cechy widziane z perspektywy różnicowania leków homeopatycznych_, nagranego 13 września 2025 r.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7de1-9cfd-39801cb77c95', 'youtube', true,
     'Dr Sanjay Modi on the magnesium salts and Kraków',
     'Dr Sanjay Modi o solach magnezu i Krakowie',
     'sanjay-modi-magnesium-salts-and-krakow',
     585, '2025-10-10', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'w_CPwlko4bo',
     'Dr Sanjay Modi talks about an underrated group of homeopathic remedies — the magnesium salts — the subject of his upcoming seminar.',
     'Dr Sanjay Modi opowiada o niedocenianej grupie leków homeopatycznych — solach magnezu — którym poświęcone będzie jego nadchodzące seminarium.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-77da-b919-a72d00537fe5', 'youtube', true,
     'Divali and magnesium — an interview with Dr Sanjay Modi',
     'Divali i magnez — rozmowa z drem Sanjayem Modim',
     'sanjay-modi-divali-and-magnesium',
     387, '2025-10-06', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'sIkcp80WIx4',
     'Dr Sanjay Modi talks about Divali and about the magnesium remedies ahead of his October seminar.',
     'Dr Sanjay Modi opowiada o święcie Divali i o lekach magnezowych przed swoim październikowym seminarium.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-704e-a7b2-df7e38fd3b40', 'youtube', true,
     'Dr Sanjay Modi webinar, part 2: terror',
     'Webinarium z drem Sanjayem Modim, część 2: przerażenie',
     'sanjay-modi-webinar-terror-pt-2',
     1054, '2025-10-02', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'QpR-wPOB8K4',
     'Part 2 of Dr Modi''s webinar _Boldness and daring vs terror and restlessness — qualities viewed from the point of view of homeopathic remedy differentiation_, recorded on 13 September 2025.',
     'Część 2 webinarium dra Modiego _Śmiałość i odwaga a przerażenie i niepokój — cechy widziane z perspektywy różnicowania leków homeopatycznych_, nagranego 13 września 2025 r.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7d0d-a38a-6710dece1740', 'youtube', true,
     '_Homeopathy — Groundbreaking Science and Global Health 2025_ conference — Dr Alex Bekker',
     'Konferencja _Homeopathy — Groundbreaking Science and Global Health 2025_ — dr Alex Bekker',
     'alex-bekker-aih-conference-2025',
     295, '2025-09-28', '01a078bc-02d5-7d1b-9081-81a257a45753', 'IbHXIQ848uQ',
     'Dr Alex Bekker talks about the American Institute of Homeopathy conference _Groundbreaking Science and Global Health 2025_, held on 17–19 October 2025 at UConn Health in Farmington, Connecticut, and online.',
     'Dr Alex Bekker opowiada o konferencji American Institute of Homeopathy _Groundbreaking Science and Global Health 2025_, która odbyła się 17–19 października 2025 r. w UConn Health w Farmington w stanie Connecticut oraz online.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7e6e-bbb4-1295085f0d0c', 'youtube', true,
     'Dr Sanjay Modi webinar, part 1: boldness',
     'Webinarium z drem Sanjayem Modim, część 1: śmiałość',
     'sanjay-modi-webinar-boldness-pt-1',
     1114, '2025-09-24', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'T2UJ20y9uwk',
     'Part 1 of Dr Modi''s webinar _Boldness and daring vs terror and restlessness — qualities viewed from the point of view of homeopathic remedy differentiation_, recorded on 13 September 2025.',
     'Część 1 webinarium dra Modiego _Śmiałość i odwaga a przerażenie i niepokój — cechy widziane z perspektywy różnicowania leków homeopatycznych_, nagranego 13 września 2025 r.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7587-8850-d12335a05b69', 'youtube', true,
     'J.T. Kent — the most misunderstood homeopath of all time?',
     'J.T. Kent — najbardziej niezrozumiany homeopata wszech czasów?',
     'sanjay-modi-on-j-t-kent',
     512, '2025-09-22', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'sIdAXJF1MrE',
     'James Tyler Kent remains one of the most influential yet most misunderstood figures of classical homeopathy. Dr Sanjay Modi discusses the difference between Kent''s clinical practice and his lectures, why remedy pictures from the _Materia Medica_ should not be mistaken for patient personalities, and how Kent''s work connects back to Hahnemann and Boger.',
     'James Tyler Kent pozostaje jedną z najbardziej wpływowych, a zarazem najczęściej źle rozumianych postaci homeopatii klasycznej. Dr Sanjay Modi omawia różnicę między praktyką kliniczną Kenta a jego wykładami, wyjaśnia, dlaczego obrazów leków z _Materia Medica_ nie należy mylić z osobowością pacjenta, i pokazuje, jak dzieło Kenta łączy się z Hahnemannem i Bogerem.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7bde-a355-f38f5916ee99', 'youtube', true,
     'Dr Modi''s seminar in Kraków is coming up',
     'Dr Modi w Krakowie — zapowiedź seminarium',
     'upcoming-sanjay-modi-seminar-krakow',
     262, '2025-09-18', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'ZZ1G0yGU-zQ',
     'An invitation to the seminar _To perfect the art of homeopathy 2_ with Dr Sanjay Modi, held on 24–25 October 2025 in Kraków and online, organised together with the Polish Homeopathic Society and the Polish Association of Homeopathic Physicians and Pharmacists.',
     'Zaproszenie na seminarium „Udoskonalić kunszt homeopatyczny 2” z drem Sanjayem Modim, 24–25 października 2025 r. w Krakowie i online, organizowane wspólnie z Polskim Towarzystwem Homeopatycznym oraz Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7aba-80b6-88bb555d4fd0', 'youtube', true,
     'Objective symptoms in homeopathy: lessons from Boger and Kevin the cat',
     'Objawy obiektywne w homeopatii: lekcje Bogera i kota Kevina',
     'alex-bekker-objective-symptoms-boger',
     1120, '2025-09-01', '01a078bc-02d5-7d1b-9081-81a257a45753', 'TQRT1QRn33w',
     'Dr Alex Bekker discusses the clinical approach of Bönninghausen and his disciple Boger, showing how objective symptoms and simple modalities guide effective homeopathic treatment — illustrated with real cases, Kevin the cat among them.',
     'Dr Alex Bekker omawia podejście kliniczne Bönninghausena i jego ucznia Bogera, pokazując, jak objawy obiektywne i proste modalności prowadzą do skutecznego leczenia homeopatycznego — na przykładach rzeczywistych przypadków, w tym kota Kevina.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7192-b0c9-37078b126f91', 'youtube', true,
     'From psychotherapy to the homeopathic proving of _Plutonium_ — an interview with Didier Lustig',
     'Od psychoterapii do homeopatycznego badania leku _Plutonium_ — rozmowa z Didierem Lustigiem',
     'didier-lustig-proving-of-plutonium',
     1051, '2025-08-18', '019f75f3-7540-7ad3-8068-2f060beaf060', 'HFHzVZDP6ZM',
     'Barbara Liberowicz asks homeopath and astrologer Didier Lustig about _Plutonium_: how he was inspired to make a remedy from it, how he obtained the radioactive material for potentisation, and how the remedy is indicated in life-threatening illness, deep depression and anger turned against others or against oneself.',
     'Barbara Liberowicz rozmawia z homeopatą i astrologiem Didierem Lustigiem o leku _Plutonium_: skąd wziął się pomysł jego sporządzenia, jak zdobył materiał promieniotwórczy do potencjonowania i w jakich stanach lek bywa wskazany — w chorobach zagrażających życiu, głębokiej depresji oraz gniewie kierowanym przeciw innym lub sobie.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7953-b235-b6f2e15f8c03', 'youtube', true,
     'A homeopathic trip to Dr V''s clinic in Mumbai',
     'Homeopatyczna wyprawa do gabinetu dra V w Bombaju',
     'homeopathic-trip-to-dr-v-clinic-in-mumbai',
     618, '2025-08-08', '019f75f3-7540-7aac-99de-4ad4bc37b9b4', '5OURyFQM7Kk',
     'Barbara Liberowicz interviews Dr Shailendra Vaishampayan (Dr V) while driving through the busy streets of Mumbai — a journey that ends with a visit to his clinic.',
     'Barbara Liberowicz rozmawia z drem Shailendrą Vaishampayanem (drem V) podczas jazdy zatłoczonymi ulicami Bombaju — podróż kończy się wizytą w jego gabinecie.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-73f7-9064-41bd9e031644', 'youtube', true,
     'Homeopathy, data mining and the landscape of repertories — an interview with Dr Rainer Schäferkordt',
     'Homeopatia, eksploracja danych i krajobraz repertoriów — rozmowa z drem Rainerem Schäferkordtem',
     'rainer-schaeferkordt-data-mining-and-repertories',
     375, '2025-08-01', '01a078bc-02d5-7637-9b37-16c3511479aa', 'Jl21A7vW24I',
     'Repertories are pivotal tools of homeopathic practice, yet their quality and scientific standards are rarely examined. Dr Rainer Schäferkordt, who has spent fifteen years analysing homeopathic data, built a new repertory from scratch, based solely on the _Materia Medica_ — _Phenomena_, made not of rubrics but of single phenomena that can be combined into any rubric.',
     'Repertoria to kluczowe narzędzia praktyki homeopatycznej, a mimo to ich jakość i standardy naukowe rzadko bywają przedmiotem analizy. Dr Rainer Schäferkordt, który od piętnastu lat analizuje dane homeopatyczne, zbudował od podstaw nowe repertorium oparte wyłącznie na _Materia Medica_ — _Phenomena_, złożone nie z rubryk, lecz z pojedynczych zjawisk, które można łączyć w dowolną rubrykę.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7daa-98c6-d98393c23e4e', 'youtube', true,
     'Homeopathy for sleep disorders — Dr Jawahar Shah on sleeplessness, the dreams of the magnesiums and Cinderella',
     'Homeopatia w zaburzeniach snu — dr Jawahar Shah o bezsenności, snach leków magnezowych i Kopciuszku',
     'jawahar-shah-homeopathy-for-sleep-disorders',
     946, '2025-07-25', '019f75d0-f770-70b0-b82f-0a4fea4cfdb8', 'Hu5y6_qGbvk',
     'Dr Jawahar Shah joins Barbara Liberowicz to explore one of the most overlooked aspects of health — sleep. Drawing on decades of practice, he explains how homeopathy treats insomnia, sleep apnoea and restless legs syndrome, and how miasmatic phases point to remedies such as the magnesium salts, _Nux vomica_ and _Ignatia_.',
     'Dr Jawahar Shah rozmawia z Barbarą Liberowicz o jednym z najbardziej pomijanych aspektów zdrowia — o śnie. Opierając się na dziesięcioleciach praktyki, wyjaśnia, jak homeopatia leczy bezsenność, bezdech senny i zespół niespokojnych nóg oraz jak fazy miazmatyczne wskazują na leki takie jak sole magnezu, _Nux vomica_ czy _Ignatia_.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7552-8b02-223994992f44', 'youtube', true,
     'Dr Sanjay Modi seminar: introverted and extraverted remedies',
     'Seminarium z drem Sanjayem Modim: leki introwertyczne i ekstrawertyczne',
     'sanjay-modi-seminar-introverted-and-extraverted-remedies',
     536, '2025-07-22', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'rSZ-I5S-U6M',
     'A fragment of Dr Modi''s seminar _To perfect the art of homeopathy_, held on 30–31 May 2025 in Poznań and online.',
     'Fragment seminarium dra Modiego „Udoskonalić kunszt homeopatyczny”, które odbyło się 30–31 maja 2025 r. w Poznaniu i online.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7b5a-9d2f-a6a3b85a790e', 'youtube', true,
     'Dr Asher Shaikh on German New Medicine, infertility and _Spongia_',
     'Dr Asher Shaikh o Germańskiej Nowej Medycynie, niepłodności i leku _Spongia_',
     'asher-shaikh-gnm-infertility-and-spongia',
     1386, '2025-07-14', '019beef9-4287-714f-982b-2524fdef7063', 'J4NRGCOdme8',
     'Fragments of the Homeo sapiens webinar with Dr Asher Shaikh: _Polycystic ovary and infertility cases — combining German New Medicine and homeopathy_.',
     'Fragmenty webinarium Homeo sapiens z drem Asherem Shaikh: _Zespół policystycznych jajników i przypadki niepłodności — połączenie Germańskiej Nowej Medycyny i homeopatii_.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-77a5-a46d-e4999678cbf8', 'youtube', true,
     'Dr Mini Mehta — endometrial hyperplasia treated with homeopathic remedies',
     'Dr Mini Mehta — rozrost endometrium leczony lekami homeopatycznymi',
     'mini-mehta-endometrial-hyperplasia',
     893, '2025-07-05', '01a078bc-02d5-7b0f-82c9-2cac1c245a6b', 'u950sOLsRwY',
     'A fragment of the Homeo sapiens webinar _Best of Utrecht_. At the LMHI congress in Utrecht, Dr Mini Mehta presented a case of endometrial hyperplasia in which the well-chosen remedy led to evacuation of the mass — "a surgery without a knife".',
     'Fragment webinarium Homeo sapiens „Najlepsze z Utrechtu”. Na kongresie LMHI w Utrechcie dr Mini Mehta przedstawiła przypadek rozrostu endometrium, w którym dobrze dobrany lek doprowadził do wydalenia zmiany — „operacji bez noża”.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-70f4-b6e5-6f94f187247f', 'youtube', true,
     'Dr Sanjay Modi seminar: _Kalium arsenicosum_',
     'Seminarium z drem Sanjayem Modim: _Kalium arsenicosum_',
     'sanjay-modi-seminar-kalium-arsenicosum',
     609, '2025-06-28', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '_niRrUwMzsM',
     'A fragment of Dr Modi''s seminar _To perfect the art of homeopathy_, held on 30–31 May 2025 in Poznań and online.',
     'Fragment seminarium dra Modiego „Udoskonalić kunszt homeopatyczny”, które odbyło się 30–31 maja 2025 r. w Poznaniu i online.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7eb1-9d32-bb209c2578a0', 'youtube', true,
     'Prof. Aaron To on AI in homeopathy: _Rhus toxicodendron_, sciatica and the beautiful liar',
     'Prof. Aaron To o sztucznej inteligencji w homeopatii: _Rhus toxicodendron_, rwa kulszowa i piękny kłamca',
     'aaron-to-on-ai-in-homeopathy',
     1235, '2025-06-24', '01a078bc-02d5-7c70-b0a8-be926bd9cd0d', 'hMBbGv8PabI',
     'Fragments of the webinar _Best of Utrecht_. Prof. Aaron To discusses data collection in homeopathy, large language models and the dangers of AI hallucination.',
     'Fragmenty webinarium „Najlepsze z Utrechtu”. Prof. Aaron To omawia gromadzenie danych w homeopatii, duże modele językowe i zagrożenia związane z halucynacjami sztucznej inteligencji.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-768c-ab05-37f52261fc83', 'youtube', true,
     'Is homeopathy effective for children? Surprising results of a clinical trial',
     'Czy homeopatia działa u dzieci? Zaskakujące wyniki badania klinicznego',
     'roee-singer-homeopathy-for-children-trial',
     716, '2025-05-23', '01a078bc-02d5-7922-b2fc-b257db6032dc', 'uDjvHRcgJKo',
     'Dr Shepherd Roee Singer discusses the results of a randomised controlled trial comparing homeopathy and conventional primary care in children during the first 24 months of life.',
     'Dr Shepherd Roee Singer omawia wyniki randomizowanego badania kontrolowanego, porównującego homeopatię z konwencjonalną podstawową opieką zdrowotną u dzieci w pierwszych 24 miesiącach życia.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7138-b118-5b2e827e2176', 'youtube', true,
     'Dr Sanjay Modi: introverts and extraverts, the kali salts, reaching precision',
     'Dr Sanjay Modi: introwertycy i ekstrawertycy, sole potasu i dążenie do precyzji',
     'sanjay-modi-introverts-extraverts-kali-salts',
     605, '2025-05-11', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'mPDkAmKHPyI',
     'Ahead of his Poznań seminar, Dr Sanjay Modi shares the small differentiating features of lesser-known remedies and explains why details from the provings matter for telling remedies apart and for precision in everyday practice.',
     'Przed seminarium w Poznaniu dr Sanjay Modi pokazuje drobne cechy różnicujące mniej znanych leków i wyjaśnia, dlaczego szczegóły z badań lekowych są kluczowe dla ich odróżniania i dla precyzji w codziennej praktyce.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-787d-9e50-836837486da4', 'youtube', true,
     'Can homeopathy help lung cancer patients? Exploring the Michael Frass study',
     'Czy homeopatia może pomóc chorym na raka płuca? O badaniu Michaela Frassa',
     'michael-frass-lung-cancer-study',
     853, '2025-05-03', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 'vGfYeojvWw8',
     'Barbara Liberowicz talks to Prof. Michael Frass about his study _Homeopathic treatment as an add-on therapy may improve quality of life and prolong survival in patients with non-small cell lung cancer_ — a prospective, randomised, placebo-controlled, double-blind, three-arm, multicentre trial.',
     'Barbara Liberowicz rozmawia z prof. Michaelem Frassem o jego badaniu _Homeopathic treatment as an add-on therapy may improve quality of life and prolong survival in patients with non-small cell lung cancer_ — prospektywnym, randomizowanym, kontrolowanym placebo, podwójnie zaślepionym, trójramiennym badaniu wieloośrodkowym.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7346-ab32-8da8a93a6925', 'youtube', true,
     'Actinide homeopathy explained: Didier Lustig on radioactive remedies',
     'Homeopatia aktynowców: Didier Lustig o lekach promieniotwórczych',
     'didier-lustig-actinide-remedies',
     843, '2025-04-24', '019f75f3-7540-7ad3-8068-2f060beaf060', 'b-WWYcSO_Yo',
     'Didier Lustig — French astrologer, homeopath, maker of the actinide remedies and author of _The Big Book of Actinides_ — talks about working with substances such as _Plutonium_ and _Neptunium_, how they were sourced and potentised, and how their themes answer moments of symbolic or literal disintegration.',
     'Didier Lustig — francuski astrolog, homeopata, twórca leków z grupy aktynowców i autor _The Big Book of Actinides_ — opowiada o pracy z substancjami takimi jak _Plutonium_ i _Neptunium_, o tym, skąd je pozyskał i jak je potencjonował, oraz o tym, jak ich tematy odpowiadają momentom symbolicznego lub dosłownego rozpadu.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7492-bc31-ee9cf58d534c', 'youtube', true,
     '_Organon_ decoded: rapidity or permanence?',
     '_Organon_ rozszyfrowany: szybkość czy trwałość?',
     'organon-decoded-rapidity-or-permanence',
     1454, '2025-04-11', '0199c2f2-528b-7e88-96e3-5e5088333a8a', '76gQITyMLHc',
     'A fragment of the Homeo sapiens webinar _Best of Sevilla 2024_, in which Dr Sanjay Modi reads the _Organon_ on what the cure should be measured by.',
     'Fragment webinarium Homeo sapiens „Najlepsze z Sewilli 2024”, w którym dr Sanjay Modi czyta _Organon_ w poszukiwaniu miary wyleczenia.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7e6d-8e80-1419458e3ddd', 'youtube', true,
     'Once you hit the bull''s-eye…',
     'Kiedy trafisz w dziesiątkę…',
     'once-you-hit-the-bulls-eye',
     858, '2025-04-06', '019beef9-4287-714f-982b-2524fdef7063', 'R9l7CSOMRe4',
     'Dr Asher Shaikh talks further about treating patients with both homeopathy and German New Medicine, and how the two approaches complement each other in his practice.',
     'Dr Asher Shaikh ponownie opowiada o swoim doświadczeniu w stosowaniu homeopatii w połączeniu z Germańską Nową Medycyną (GNM) i o tym, jak oba systemy uzupełniają się w jego praktyce.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7a25-bd3a-b9feeddd47cf', 'youtube', true,
     'Dutiful people',
     'Ludzie obowiązkowi',
     'dutiful-people',
     774, '2025-03-30', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'hr5ELtBoIy8',
     'Dr Sanjay Modi explains how to differentiate three duty-conscious remedies — _Aurum metallicum_, _Nux vomica_ and _Silicea_ — pointing to their common features of fear of failure, exhaustion from work and conscientiousness, and to the motivation behind each. Excerpts from his webinar of 22 March 2025.',
     'Dr Sanjay Modi wyjaśnia, jak różnicować trzy leki o dużym poczuciu obowiązku — _Aurum metallicum_, _Nux vomica_ i _Silicea_. Wskazuje na cechy wspólne: lęk przed niepowodzeniem, wyczerpanie pracą i sumienność, oraz na motywacje stojące za każdym z nich. Fragmenty webinarium z 22 marca 2025 r.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-72f5-8be5-cc713aa4cd79', 'youtube', true,
     'Match the patient to the remedy, or the remedy to the patient?',
     'Lek do pacjenta czy pacjent do leku?',
     'match-patient-to-remedy-or-remedy-to-patient',
     766, '2025-03-13', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'yvC1L6ibPVw',
     'Barbara Liberowicz talks to Dr Sanjay Modi about the subtle difference between matching a remedy to a patient and matching a patient to a remedy. Dr Modi underlines the need to remain an unprejudiced observer — which also means abandoning fixed remedy pictures instead of projecting them onto the patient.',
     'Barbara Liberowicz rozmawia z drem Sanjayem Modim o dyskretnej różnicy między dobraniem leku do pacjenta a dobraniem pacjenta do leku. Dr Modi zwraca uwagę na konieczność pozostania „nieuprzedzonym obserwatorem”, co oznacza także pozbycie się zafiksowanych obrazów leków i unikanie projektowania ich na pacjenta.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-771d-a0fd-7846aedc0cc1', 'youtube', true,
     'How German New Medicine and homeopathy work together — Dr Asher Shaikh''s holistic approach',
     'Jak Germańska Nowa Medycyna i homeopatia współgrają — holistyczne podejście dra Ashera Shaikha',
     'how-gnm-and-homeopathy-work-together',
     1176, '2025-03-09', '019beef9-4287-714f-982b-2524fdef7063', 'K8WJlg_zP38',
     'Dr Asher Shaikh on homeopathy and German New Medicine, and on how he brings the two together in daily practice.',
     'Dr Asher Shaikh o homeopatii i Germańskiej Nowej Medycynie oraz o tym, jak łączy oba podejścia w codziennej praktyce.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-799d-80cd-1b31cf6cba3b', 'youtube', true,
     'Prof. Michael Frass on intensive care, sepsis, COPD and homeopathy',
     'Prof. Michael Frass o intensywnej terapii, sepsie, POChP i homeopatii',
     'michael-frass-icu-sepsis-copd',
     2021, '2025-02-22', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 'bBQOs3QLGWU',
     'Prof. Michael Frass — internist, intensivist and retired professor of the Medical University of Vienna — shares his experience of homeopathic treatment in the ICU: septic states, respiratory failure, COPD, embolism and intoxications. He is the author of _Homeopathy in Intensive Care and Emergency Medicine_.',
     'Prof. Michael Frass — internista, specjalista intensywnej terapii i emerytowany profesor Uniwersytetu Medycznego w Wiedniu — dzieli się doświadczeniem leczenia homeopatycznego na OIOM-ie: w stanach septycznych, niewydolności oddechowej, POChP, zatorowości i zatruciach. Jest autorem książki _Homeopathy in Intensive Care and Emergency Medicine_.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7e88-9917-9f77ef0a52ad', 'youtube', true,
     'Why is the _Organon_ still so important? (Polish-subtitled edition)',
     'Dlaczego _Organon_ jest wciąż tak ważny?',
     'dlaczego-organon-jest-wciaz-tak-wazny',
     724, '2025-02-13', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 'H3eBg-ieT_E',
     'Barbara Liberowicz invites Dr Robert Zawiślak and Dr Sanjay Modi to discuss the significance of Samuel Hahnemann''s _Organon_ — the life''s work of the founder of homeopathy, still a point of reference and a source of principles for those who practise it.',
     'Barbara Liberowicz zaprosiła dra Roberta Zawiślaka i dra Sanjaya Modiego do rozmowy o znaczeniu _Organonu_ Samuela Hahnemanna. Dzieło życia twórcy homeopatii nie straciło na aktualności — pozostaje punktem odniesienia i źródłem zasad dla praktykujących homeopatów.')
on conflict (slug) do nothing;

insert into videos (id, provider, is_public, title_en, title_pl, slug,
                    duration_seconds, recorded_on, host_id, youtube_id,
                    description_en, description_pl) values
    ('01a078bc-02d5-7be1-b71f-fe85962e9d06', 'youtube', true,
     'Dr Kishore Mehta: homeopathy, Ayurveda and multiple sclerosis',
     'Dr Kishore Mehta: homeopatia, ajurweda i stwardnienie rozsiane',
     'kishore-mehta-homeopathy-ayurveda-multiple-sclerosis',
     1595, '2025-02-08', '01a078bc-02d5-76a1-8d0b-2ed2d129ba3f', 'uwLzCJqfd4w',
     'A conversation with Dr Kishore Mehta at his Mumbai clinic on the place of homeopathy and Ayurveda in modern medicine, the treatment of multiple sclerosis, and why homeopathy is so widely used in India.',
     'Rozmowa z drem Kishorem Mehtą w jego gabinecie w Bombaju o miejscu homeopatii i ajurwedy we współczesnej medycynie, o leczeniu stwardnienia rozsianego i o tym, dlaczego homeopatia jest w Indiach tak popularna.')
on conflict (slug) do nothing;


-- Video hosts (skipped for any video the insert above left out)

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7122-b676-131b8db65b5b', '01a078bc-02d5-7cd3-8c21-89f49425975e', '01a078bc-02d5-7a25-9831-08c44b13dcf2', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7cd3-8c21-89f49425975e')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7d20-bebf-d024971305d8', '01a078bc-02d5-7ead-9326-20c2e24c6903', '01a078bc-02d5-7e8b-9b76-ff548969e3ee', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7ead-9326-20c2e24c6903')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7665-8b7a-d8d21ebc4e07', '01a078bc-02d5-7679-9180-5588223d5a4e', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7679-9180-5588223d5a4e')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7450-9c47-705c35d3ba83', '01a078bc-02d5-7f8f-9484-fc9b4de844a4', '01a078bc-02d5-75eb-8846-d8a867303094', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7f8f-9484-fc9b4de844a4')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-763f-b312-6aae5abb12cb', '01a078bc-02d5-74d7-a227-544d0fb5fa89', '01a078bc-02d5-7a25-9831-08c44b13dcf2', 1
where exists (select 1 from videos where id = '01a078bc-02d5-74d7-a227-544d0fb5fa89')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7944-9a8b-4e29d129d11b', '01a078bc-02d5-7246-af71-1dd0a4b3c8d5', '01a078bc-02d5-7f6c-ada5-0645fc379542', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7246-af71-1dd0a4b3c8d5')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-753a-8aa8-28c326197f58', '01a078bc-02d5-7d9e-9f69-6386317ef58b', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7d9e-9f69-6386317ef58b')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7897-979f-2527a4576592', '01a078bc-02d5-7d9e-9f69-6386317ef58b', '019f75d0-f76f-7b99-abf8-e40499c35222', 2
where exists (select 1 from videos where id = '01a078bc-02d5-7d9e-9f69-6386317ef58b')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7f4a-8964-2f2ff4306128', '01a078bc-02d5-7fe5-a978-9b0d1d1994c7', '01a078bc-02d5-707e-98f6-88c43ea68c0b', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7fe5-a978-9b0d1d1994c7')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7156-a89e-65206f2c02ce', '01a078bc-02d5-7215-992a-65ab92a3b638', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7215-992a-65ab92a3b638')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7f60-88c8-a0832d2fb656', '01a078bc-02d5-7f90-8c18-c1add5aeeb70', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7f90-8c18-c1add5aeeb70')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-745b-9a61-74335617d1de', '01a078bc-02d5-752b-af53-c1f0883573fb', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 1
where exists (select 1 from videos where id = '01a078bc-02d5-752b-af53-c1f0883573fb')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7c9a-8591-e48fbae9801a', '01a078bc-02d5-761c-88cc-59337d2cc1d5', '01a078bc-02d5-7b0f-82c9-2cac1c245a6b', 1
where exists (select 1 from videos where id = '01a078bc-02d5-761c-88cc-59337d2cc1d5')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-71d3-960e-112bfd62e7a4', '01a078bc-02d5-7e80-bb19-f91f2c5d0418', '019f75f3-7540-7aac-99de-4ad4bc37b9b4', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7e80-bb19-f91f2c5d0418')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-72d0-8710-1ecacf6b63f1', '01a078bc-02d5-76bf-839f-fb4a0c978a57', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-76bf-839f-fb4a0c978a57')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7b62-99d8-0dcb7f1b2240', '01a078bc-02d5-76bf-839f-fb4a0c978a57', '019f75d0-f76f-7b99-abf8-e40499c35222', 2
where exists (select 1 from videos where id = '01a078bc-02d5-76bf-839f-fb4a0c978a57')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7471-a095-f95d738676d6', '01a078bc-02d5-79b3-a127-e81f9776c905', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-79b3-a127-e81f9776c905')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-77de-882e-92fd2775aa04', '01a078bc-02d5-7de1-9cfd-39801cb77c95', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7de1-9cfd-39801cb77c95')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7ad1-8f3f-8cf1ae931e2a', '01a078bc-02d5-77da-b919-a72d00537fe5', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-77da-b919-a72d00537fe5')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-762f-9b37-90fd640de383', '01a078bc-02d5-704e-a7b2-df7e38fd3b40', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-704e-a7b2-df7e38fd3b40')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7bce-9728-4f6584facbb0', '01a078bc-02d5-7d0d-a38a-6710dece1740', '01a078bc-02d5-7d1b-9081-81a257a45753', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7d0d-a38a-6710dece1740')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7695-8b0c-85aab7c4963a', '01a078bc-02d5-7e6e-bbb4-1295085f0d0c', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7e6e-bbb4-1295085f0d0c')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7a36-86e9-024d169b09fd', '01a078bc-02d5-7587-8850-d12335a05b69', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7587-8850-d12335a05b69')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7d35-85ed-9629878f6143', '01a078bc-02d5-7bde-a355-f38f5916ee99', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7bde-a355-f38f5916ee99')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7268-8c2f-36102d9943fa', '01a078bc-02d5-7aba-80b6-88bb555d4fd0', '01a078bc-02d5-7d1b-9081-81a257a45753', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7aba-80b6-88bb555d4fd0')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7e6a-a9ac-3151dbfaf750', '01a078bc-02d5-7192-b0c9-37078b126f91', '019f75f3-7540-7ad3-8068-2f060beaf060', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7192-b0c9-37078b126f91')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7f65-8480-971b12a6a075', '01a078bc-02d5-7953-b235-b6f2e15f8c03', '019f75f3-7540-7aac-99de-4ad4bc37b9b4', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7953-b235-b6f2e15f8c03')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7754-a20e-7eb942c43967', '01a078bc-02d5-73f7-9064-41bd9e031644', '01a078bc-02d5-7637-9b37-16c3511479aa', 1
where exists (select 1 from videos where id = '01a078bc-02d5-73f7-9064-41bd9e031644')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7de8-9311-0a830510d2d0', '01a078bc-02d5-7daa-98c6-d98393c23e4e', '019f75d0-f770-70b0-b82f-0a4fea4cfdb8', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7daa-98c6-d98393c23e4e')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7d60-8230-afa7436bca20', '01a078bc-02d5-7552-8b02-223994992f44', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7552-8b02-223994992f44')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7b85-ba3e-7e2dc1cad2ef', '01a078bc-02d5-7b5a-9d2f-a6a3b85a790e', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7b5a-9d2f-a6a3b85a790e')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-73ed-99be-2b8b2183342f', '01a078bc-02d5-77a5-a46d-e4999678cbf8', '01a078bc-02d5-7b0f-82c9-2cac1c245a6b', 1
where exists (select 1 from videos where id = '01a078bc-02d5-77a5-a46d-e4999678cbf8')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7017-836a-05fb54102622', '01a078bc-02d5-70f4-b6e5-6f94f187247f', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-70f4-b6e5-6f94f187247f')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7f49-80d9-fb54312844e3', '01a078bc-02d5-7eb1-9d32-bb209c2578a0', '01a078bc-02d5-7c70-b0a8-be926bd9cd0d', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7eb1-9d32-bb209c2578a0')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7657-97eb-71a2bd9798c7', '01a078bc-02d5-768c-ab05-37f52261fc83', '01a078bc-02d5-7922-b2fc-b257db6032dc', 1
where exists (select 1 from videos where id = '01a078bc-02d5-768c-ab05-37f52261fc83')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-72d5-8f3f-5dacd1260de2', '01a078bc-02d5-7138-b118-5b2e827e2176', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7138-b118-5b2e827e2176')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7757-bb3e-709b356ddb37', '01a078bc-02d5-787d-9e50-836837486da4', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 1
where exists (select 1 from videos where id = '01a078bc-02d5-787d-9e50-836837486da4')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7085-b715-a65bfd8c751f', '01a078bc-02d5-7346-ab32-8da8a93a6925', '019f75f3-7540-7ad3-8068-2f060beaf060', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7346-ab32-8da8a93a6925')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7371-a1e6-879edd0e9f0e', '01a078bc-02d5-7492-bc31-ee9cf58d534c', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7492-bc31-ee9cf58d534c')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7b04-a88b-1be5fd4fa73e', '01a078bc-02d5-7e6d-8e80-1419458e3ddd', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7e6d-8e80-1419458e3ddd')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7139-9e37-af2e7ab2cb6b', '01a078bc-02d5-7a25-bd3a-b9feeddd47cf', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7a25-bd3a-b9feeddd47cf')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7bda-9f3b-3be2ca0c5ef7', '01a078bc-02d5-72f5-8be5-cc713aa4cd79', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-72f5-8be5-cc713aa4cd79')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7cfb-b37e-61c74c1de6bb', '01a078bc-02d5-771d-a0fd-7846aedc0cc1', '019beef9-4287-714f-982b-2524fdef7063', 1
where exists (select 1 from videos where id = '01a078bc-02d5-771d-a0fd-7846aedc0cc1')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7b3d-9a20-679eca041ff4', '01a078bc-02d5-799d-80cd-1b31cf6cba3b', '019f75f3-7540-7971-b2e4-a9fca75ab7ab', 1
where exists (select 1 from videos where id = '01a078bc-02d5-799d-80cd-1b31cf6cba3b')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7085-b2a5-f4c6f388399f', '01a078bc-02d5-7e88-9917-9f77ef0a52ad', '0199c2f2-528b-7e88-96e3-5e5088333a8a', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7e88-9917-9f77ef0a52ad')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7fb5-9870-3ce81beda1be', '01a078bc-02d5-7e88-9917-9f77ef0a52ad', '019f75d0-f76f-7b99-abf8-e40499c35222', 2
where exists (select 1 from videos where id = '01a078bc-02d5-7e88-9917-9f77ef0a52ad')
on conflict do nothing;

insert into video_hosts (id, video_id, host_id, position)
select '01a078bc-02d5-7a88-8074-528219ec2850', '01a078bc-02d5-7be1-b71f-fe85962e9d06', '01a078bc-02d5-76a1-8d0b-2ed2d129ba3f', 1
where exists (select 1 from videos where id = '01a078bc-02d5-7be1-b71f-fe85962e9d06')
on conflict do nothing;


commit;
