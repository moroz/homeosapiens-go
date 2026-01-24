package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
)

func MustParseUUID(u string) pgtype.UUID {
	parsed := uuid.MustParse(u)
	return pgtype.UUID{
		Valid: true,
		Bytes: parsed,
	}
}

func MustParseTimestamp(t string) pgtype.Timestamp {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Fatal(err)
	}
	return pgtype.Timestamp{
		Time:  parsed,
		Valid: true,
	}
}

func StringPtr(s string) *string {
	return &s
}

func MustParseDecimal(d string) decimal.Decimal {
	parsed, err := decimal.NewFromString(d)
	if err != nil {
		log.Fatal(err)
	}
	return parsed
}

func MaybeDecimal(d *string) *decimal.Decimal {
	if d == nil {
		return nil
	}
	val := MustParseDecimal(*d)
	return &val
}

func main() {
	fmt.Println(config.DatabaseUrl)
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(context.Background())

	log.Printf("Cleaning database...")
	_, err = db.Exec(context.Background(), "truncate events, hosts, assets, venues, events_hosts, event_prices, event_registrations, user_tokens, videos, video_sources")
	if err != nil {
		log.Fatal(err)
	}

	users := []*types.CreateUserParams{
		{
			Email:      "karol@moroz.dev",
			GivenName:  "Karol",
			FamilyName: "Moroz",
			Country:    "PL",
			Password:   "foobar",
			Role:       queries.UserRoleAdministrator,
		},
		{
			GivenName:  "Sanjay",
			FamilyName: "Modi",
			Email:      "sanjay.modi@example.com",
			Country:    "IN",
			Password:   "foobar",
			Role:       queries.UserRoleRegular,
		},
	}

	log.Printf("Creating users...")
	userService := services.NewUserService(db)
	for _, user := range users {
		if _, err := userService.CreateUser(context.Background(), user); err != nil {
			log.Fatal(err)
		}
	}

	assets := []*types.CreateAssetParams{
		{
			ID:               MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8b"),
			ObjectKey:        "cm7uqj3q500mglz8z2dqy8sdz.webp",
			OriginalFilename: "cm7uqj3q500mglz8z2dqy8sdz.webp",
		},
		{
			ID:               MustParseUUID("019b0c7c-c3c4-71c3-a630-7b33a847ca2a"),
			ObjectKey:        "019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg",
			OriginalFilename: "019b0c7c-c3c4-71c3-a630-7b33a847ca2a.jpg",
		},
		{
			ID:               MustParseUUID("019beef9-ad4c-736f-9bb0-965b59ca21ae"),
			ObjectKey:        "019beef9-ad4c-736f-9bb0-965b59ca21ae.png",
			OriginalFilename: "drasher.png",
		},
	}
	log.Printf("Creating assets...")
	for _, asset := range assets {
		params := &queries.UpsertAssetParams{
			ID:               asset.ID,
			ObjectKey:        asset.ObjectKey,
			OriginalFilename: &asset.OriginalFilename,
		}
		if _, err := queries.New(db).UpsertAsset(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	hosts := []*types.CreateHostParams{
		{
			ID:               MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8a"),
			Salutation:       "common.hosts.salutation.dr",
			GivenName:        "Sanjay",
			FamilyName:       "Modi",
			ProfilePictureId: MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8b"),
			Country:          "IN",
		},
		{
			ID:               MustParseUUID("019beef9-4287-714f-982b-2524fdef7063"),
			Salutation:       "common.hosts.salutation.dr",
			GivenName:        "Asher",
			FamilyName:       "Shaikh",
			ProfilePictureId: MustParseUUID("019beef9-ad4c-736f-9bb0-965b59ca21ae"),
			Country:          "IN",
		},
		{
			ID:               MustParseUUID("019b0c71-fde2-76b7-8c71-21c2e9ea23a5"),
			Salutation:       "common.hosts.salutation.dr",
			GivenName:        "Herman",
			FamilyName:       "Jeggels",
			ProfilePictureId: MustParseUUID("019b0c7c-c3c4-71c3-a630-7b33a847ca2a"),
			Country:          "ZA",
		},
	}
	log.Printf("Creating hosts...")
	for _, host := range hosts {
		params := &queries.UpsertHostParams{
			ID:               host.ID,
			Salutation:       &host.Salutation,
			GivenName:        host.GivenName,
			FamilyName:       host.FamilyName,
			ProfilePictureID: host.ProfilePictureId,
			Country:          &host.Country,
		}
		if _, err := queries.New(db).UpsertHost(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	venues := []*types.CreateVenueParams{
		{
			ID:          MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8d"),
			NameEn:      "Vienna House Easy By Wyndham Cracow",
			CityEn:      "Cracow",
			CityPl:      StringPtr("Kraków"),
			CountryCode: "PL",
			Street:      "ul. Przy Rondzie 2",
		},
		{
			ID:          MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8e"),
			NameEn:      "IOR Hotel",
			NamePl:      StringPtr("Hotel IOR"),
			CityEn:      "Poznań",
			CountryCode: "PL",
			Street:      "ul. Węgorka 20",
			PostalCode:  "60-318",
		},
	}
	log.Printf("Creating venues...")
	for _, venue := range venues {
		params := &queries.UpsertVenueParams{
			ID:          venue.ID,
			NameEn:      venue.NameEn,
			NamePl:      venue.NamePl,
			Street:      venue.Street,
			CityEn:      venue.CityEn,
			CityPl:      venue.CityPl,
			PostalCode:  &venue.PostalCode,
			CountryCode: venue.CountryCode,
		}
		if _, err := queries.New(db).UpsertVenue(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	events := []*types.CreateEventParams{
		{
			ID:            MustParseUUID("019b0c80-a410-7728-ab6b-c1eff529dfd1"),
			EventType:     queries.EventTypeWebinar,
			TitleEn:       "A Series of Critical Cardiac Cases",
			TitlePl:       "Seria krytycznych problemów kardiologicznych",
			Slug:          "a-series-of-critical-cardiac-cases",
			StartsAt:      MustParseTimestamp("2025-12-13T16:00:00Z"),
			EndsAt:        MustParseTimestamp("2025-12-13T17:30:00Z"),
			IsVirtual:     true,
			DescriptionEn: "Dear Homeopathic Friends, \n\nWe are happy to invite you to the next Homeo sapiens Academy webinar. Experienced clinician and homeopath Dr. Herman Jeggels from Cape Town, South Africa will discuss homeopathic treatment in advanced circulatory pathology. He will present documented cases of infective endocarditis, complete AV block and heart failure.\n\nThe webinar will be hosted on December 13th 10.00am CET (Poland) / 2.30pm IST (India) / 11:00am SAST (South Africa).\n\nThe webinar is free of charge. It will be held in English with consecutive translation to Polish.\n\nIt will be held on Zoom via our website (you need to register using email address and a password).",
			DescriptionPl: StringPtr(`Z przyjemnością zapraszamy na kolejny webinar organizowany przez Homeo sapiens. Naszym gościem będzie Dr Herman Jeggels z Cape Town (RPA), doświadczony klinicysta i homeopata, który przedstawi przypadki leczenia homeopatycznego w zaawansowanych chorobach układu krążenia. Omówione zostaną udokumentowane historie leczenia pacjentów z problemami kardiologicznymi, między innymi infekcyjne zapalenie wsierdzia, całkowity blok przedsionkowo-komorowy i niewydolność serca.

Webinar odbędzie się 13 grudnia 2025 o godzinie 10.00 czasu polskiego.

Wykład będzie tłumaczony konsekutywnie na język polski.

Webinar jest bezpłatny. Odbędzie się na platformie Zoom za pośrednictwem naszej strony internetowej. Wymagana jest rejestracja z użyciem adresu email i ustawienie hasła.`),
		},
		{
			ID:                MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8c"),
			EventType:         queries.EventTypeSeminar,
			TitleEn:           "To Perfect the Art of Homeopathy",
			TitlePl:           "Udoskonalić kunszt homeopatyczny",
			Slug:              "to-perfect-the-art-of-homeopathy",
			StartsAt:          MustParseTimestamp("2025-05-30T14:00:00Z"),
			EndsAt:            MustParseTimestamp("2025-05-31T08:00:00Z"),
			IsVirtual:         true,
			DescriptionEn:     "Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.",
			VenueID:           MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8e"),
			BasePriceAmount:   StringPtr("580.00000000"),
			BasePriceCurrency: StringPtr("PLN"),
			DescriptionPl: StringPtr(`Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

30-31 maja 2025, sala wykładowa B Instytutu Ochrony Roślin, ul. Władysława Węgorka 20, 60-318 Poznań.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.

Omówionych zostanie szereg praktycznych problemów klinicznych, różnicowanie leków z grupy Kalium, leki introwertyczne/ekstrawertyczne, prezentacja przypadków klinicznych.`),
		},
		{
			ID:            MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			EventType:     queries.EventTypeSeminar,
			TitleEn:       "To Perfect the Art of Homeopathy 2",
			TitlePl:       "Udoskonalić kunszt homeopatyczny 2",
			Slug:          "to-perfect-the-art-of-homeopathy-2",
			StartsAt:      MustParseTimestamp("2025-10-24T14:00:00Z"),
			EndsAt:        MustParseTimestamp("2025-10-26T11:30:00Z"),
			IsVirtual:     true,
			DescriptionEn: "Dr. Sanjay Modi, former professor of Mumbai Homeopathic College. The webinar is organised in honorary cooperation with the Polish Homeopathic Society and the Polish Society of Homeopathic Doctors and Pharmacists.\n\nOctober 24-25 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków, Poland.\n\nOnline mode will also available (through Zoom). The lectures will be held in English with consecutive translation to Polish.",
			DescriptionPl: StringPtr(`Wykładowca Dr. Sanjay Modi, wieloletni wykładowca Mumbai Homeopathic College.

Seminarium organizowane jest we współpracy z Polskim Towarzystwem Homeopatycznym i Polskim Stowarzyszeniem Homeopatów Lekarzy i Farmaceutów.

24-25 października 2025, Vienna House Easy By Wyndham Cracow ul. Przy Rondzie 2, Kraków.

Seminarium będzie również dostępne na żywo on-line na platformie Zoom za pośrednictwem naszej strony internetowej. Wykłady będą prowadzone w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dla osób, które nie będą mogły wziąć udziału w szkoleniu w podanym terminie przewidujemy opcję udostępnienia nagrania, ale tylko dla zarejestrowanych uczestników.`),
			VenueID:           MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8d"),
			BasePriceAmount:   StringPtr("640.00000000"),
			BasePriceCurrency: StringPtr("PLN"),
		},
		{
			ID:        MustParseUUID("019bef00-6ef2-7636-9a15-c8cd1e87b997"),
			EventType: queries.EventTypeWebinar,
			TitleEn:   "What prevents me from moving on? Combining German New Medicine and Homeopathy for musculoskeletal problems",
			TitlePl:   "What prevents me from moving on? Combining German New Medicine and Homeopathy for musculoskeletal problems",
			Slug:      "what-prevents-me-from-moving-on",
			IsVirtual: true,
			StartsAt:  MustParseTimestamp("2026-02-08T15:00:00Z"),
			EndsAt:    MustParseTimestamp("2026-02-08T16:30:00Z"),
			DescriptionEn: `[EN] We kindly invite you to another free Homeo sapiens webinar. Experienced Homeopath Dr Asher Shaikh will share how he uses German New Medicine to facilitate homeopathic case-taking and remedy choice. Several musculoskeletal problems will be discussed, both theory and case-studies. The webinar will be held in English with consecutive translation to Polish.

Dr Asher Shaikh (India) – a homeopathic doctor with over 25 years of clinical experience. He is the Director of Asher Clinics – a network of 12 clinics in Mumbai, Pune, Dubai, and Nasik – and a mentor in German New Medicine, which he has taught in Dubai, India, Austria, and Israel. He currently serves as a professor at the Homoeopathic Medical College in Nasik and as the Director of Viveda Resort – an innovative holistic health center. He is the former president of the Indian Institute of Homoeopathic Physicians. Dr. Shaikh specializes in reversing autoimmune disorders.`,
			DescriptionPl: StringPtr(`Zapraszamy na kolejny darmowy webinar Homeo sapiens. Doświadczony homeopata dr Asher Shaikh opowie o sposobie, w jaki zastosowanie Nowej Germańskiej Medycyny (GNM) wspomaga przy homeopatycznym doborze leków. Podstawą do dyskusji na ten temat będzie omówienie kilku problemów układu mięśniowo-szkieletowego, zarówno teoretycznie, jak i w oparciu o studia przypadków. Webinar będzie prowadzony w języku angielskim z konsekutywnym tłumaczeniem na polski.

Dr Asher Shaikh (Indie) - lekarz homeopata z ponad 25-letnim doświadczeniem klinicznym. Jest dyrektorem Asher Clinics - 12 klinik w Mumbaju, Pune, Dubaju i Nasiku oraz mentorem Nowej Germańskiej Medycyny, którą wykładał w Dubaju, Indiach, Austrii i Izraelu. Pełni funkcję profesora w Homoeopathic Medical College w Nasiku oraz dyrektora Viveda Resort – innowacyjnego ośrodka zdrowia holistycznego. Były przewodniczący Indian Institute of Homoeopathic Physicians. Specjalizuje się w odwracaniu chorób autoimmunologicznych.`),
		},
	}
	log.Printf("Creating events...")
	q := queries.New(db)
	for _, event := range events {
		params := &queries.UpsertEventParams{
			ID:                event.ID,
			EventType:         event.EventType,
			TitleEn:           event.TitleEn,
			TitlePl:           event.TitlePl,
			Slug:              event.Slug,
			StartsAt:          event.StartsAt,
			EndsAt:            event.EndsAt,
			IsVirtual:         event.IsVirtual,
			DescriptionEn:     event.DescriptionEn,
			DescriptionPl:     event.DescriptionPl,
			VenueID:           event.VenueID,
			BasePriceAmount:   MaybeDecimal(event.BasePriceAmount),
			BasePriceCurrency: event.BasePriceCurrency,
		}
		if _, err := q.UpsertEvent(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	eventHosts := []*types.CreateEventHostParams{
		{
			EventID:  MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8c"),
			HostID:   MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8a"),
			Position: 0,
		},
		{
			EventID:  MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			HostID:   MustParseUUID("0199c2f2-528b-7e88-96e3-5e5088333a8a"),
			Position: 0,
		},
		{
			EventID:  MustParseUUID("019b0c80-a410-7728-ab6b-c1eff529dfd1"),
			HostID:   MustParseUUID("019b0c71-fde2-76b7-8c71-21c2e9ea23a5"),
			Position: 0,
		},
		{
			EventID:  MustParseUUID("019bef00-6ef2-7636-9a15-c8cd1e87b997"),
			HostID:   MustParseUUID("019beef9-4287-714f-982b-2524fdef7063"),
			Position: 0,
		},
	}
	log.Printf("Creating event hosts...")
	for _, eventHost := range eventHosts {
		params := &queries.UpsertEventHostParams{
			EventID:  eventHost.EventID,
			HostID:   eventHost.HostID,
			Position: eventHost.Position,
		}
		if _, err := q.UpsertEventHost(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	eventPrices := []*types.CreateEventPriceParams{
		{
			EventID:       MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			PriceAmount:   "560.00000000",
			PriceCurrency: "PLN",
			RuleType:      queries.PriceRuleTypeEarlyBird,
			ValidUntil:    MustParseTimestamp("2025-09-20T21:59:59Z"),
			Priority:      10,
			IsActive:      true,
			PriceType:     queries.PriceTypeFixed,
		},
		{
			EventID:       MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			PriceAmount:   "500.00000000",
			PriceCurrency: "PLN",
			RuleType:      queries.PriceRuleTypeDiscountCode,
			DiscountCode:  StringPtr("wshlif"),
			Priority:      20,
			IsActive:      true,
			PriceType:     queries.PriceTypeFixed,
		},
	}
	log.Printf("Creating event prices...")
	for _, price := range eventPrices {
		params := &queries.UpsertEventPriceParams{
			EventID:       price.EventID,
			PriceType:     price.PriceType,
			RuleType:      price.RuleType,
			PriceAmount:   MustParseDecimal(price.PriceAmount),
			PriceCurrency: price.PriceCurrency,
			DiscountCode:  price.DiscountCode,
			Priority:      price.Priority,
			IsActive:      price.IsActive,
			ValidFrom:     price.ValidFrom,
			ValidUntil:    price.ValidUntil,
		}
		if _, err := q.UpsertEventPrice(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	videos := []*types.CreateVideoParams{
		{
			ID:       MustParseUUID("019a8668-bb4f-7c9c-b9b8-3f274de96566"),
			EventID:  MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			Provider: queries.VideoProviderCloudfront,
			TitleEn:  "Day 1, Part 1",
			TitlePl:  "Dzień 1, Część 1",
			Slug:     "day-1-part-1",
			IsPublic: false,
		},
		{
			ID:       MustParseUUID("019a8ba5-fe29-7af8-bf54-b8d96af38461"),
			EventID:  MustParseUUID("0199c2fa-7e9d-72f6-ada1-88b5d04d9a58"),
			Provider: queries.VideoProviderCloudfront,
			TitleEn:  "Day 1, Part 2",
			TitlePl:  "Dzień 1, Część 2",
			Slug:     "day-1-part-2",
			IsPublic: false,
		},
	}
	log.Printf("Creating videos...")
	for _, video := range videos {
		params := &queries.UpsertVideoParams{
			ID:       video.ID,
			EventID:  video.EventID,
			Provider: video.Provider,
			TitleEn:  video.TitleEn,
			TitlePl:  video.TitlePl,
			Slug:     video.Slug,
			IsPublic: video.IsPublic,
		}
		if _, err := q.UpsertVideo(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	videoSources := []*types.CreateVideoSourceParams{
		{
			ID:          MustParseUUID("019a8ba6-c5ae-7f6f-becb-94b6957a52b2"),
			VideoID:     MustParseUUID("019a8668-bb4f-7c9c-b9b8-3f274de96566"),
			ContentType: "video/mp4",
			Codec:       "hev1",
			ObjectKey:   "/videos/019a8668-bb4f-7c9c-b9b8-3f274de96566/hevc_1080.mp4",
		},
		{
			ID:          MustParseUUID("019a8ba7-d04b-77ec-92c6-f76b6ec0e7ea"),
			VideoID:     MustParseUUID("019a8668-bb4f-7c9c-b9b8-3f274de96566"),
			ContentType: "video/webm",
			Codec:       "vp9,opus",
			ObjectKey:   "/videos/019a8668-bb4f-7c9c-b9b8-3f274de96566/webm_1080.webm",
		},
		{
			ID:          MustParseUUID("019a8bab-135e-7321-9857-f74d2dcda427"),
			VideoID:     MustParseUUID("019a8ba5-fe29-7af8-bf54-b8d96af38461"),
			ContentType: "video/mp4",
			Codec:       "hev1",
			ObjectKey:   "/videos/019a8ba5-fe29-7af8-bf54-b8d96af38461/hevc_1080.mp4",
		},
		{
			ID:          MustParseUUID("019a8bab-bc67-76f9-bf80-902043c922e6"),
			VideoID:     MustParseUUID("019a8ba5-fe29-7af8-bf54-b8d96af38461"),
			ContentType: "video/webm",
			Codec:       "vp9,opus",
			ObjectKey:   "/videos/019a8ba5-fe29-7af8-bf54-b8d96af38461/webm_1080.webm",
		},
	}
	log.Printf("Creating video sources...")
	for _, source := range videoSources {
		params := &queries.UpsertVideoSourceParams{
			ID:          source.ID,
			VideoID:     source.VideoID,
			ContentType: source.ContentType,
			Codec:       StringPtr(source.Codec),
			ObjectKey:   source.ObjectKey,
		}
		if _, err := q.UpsertVideoSource(context.Background(), params); err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit(context.Background())
}
