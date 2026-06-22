# Homeo Sapiens — Project Specification (v1: replace EzyCourse)

> Status: planning. Owner: developer. Client: site owner (homeosapiens.eu).
> Stack: Go (Echo v5) · gomponents/templ · pgx + sqlc · River jobs · PostgreSQL · Stripe · SMTP · go-i18n (PL/EN).

## 1. Goal & definition of done

**v1 ships when EzyCourse can be cancelled.** Everything the client currently does on the SaaS must work here:

- Public site: Home, About Us, Watch (free video landing), Events (list + detail + register), login/signup.
- Logged-in: "My products" video gallery, watch purchased/granted recordings.
- Paid events sold online (Stripe), free events via one-click RSVP.
- **The client manages all content herself** (events, prices, recordings) through an admin UI — no SQL, no redeploy.
- Registration & payment emails go out automatically (confirmation, venue details, receipt, reminder). Zoom join links deferred to a later Zoom-API integration (§4.9).

Success metric: one full event lifecycle runs end-to-end without the developer touching the database — client creates a paid webinar, a user buys it, attends, and later watches the recording.

## 2. Core design decision (locked 2026-06-22, after stakeholder discussion)

**Recordings and events are separate products. A video group is its own product; it can be bought self-serve (fully automated), and is also granted manually to people who attended the corresponding seminar.**

```
video_groups.product_id ──► products(title_*, base_price_*)
                                  │
                                  ├──► product_prices (optional)
                                  │
   ┌──────────────────────────────┼───────────────────────────────┐
   │ A. SELF-SERVE PURCHASE        │ B. MANUAL GRANT                │
   │   (fully automated)           │   (admin, after a seminar /    │
   │   buy → cart → Stripe →       │    off-platform payment)       │
   │   webhook                     │   admin inserts directly       │
   └──────────────────────────────┼───────────────────────────────┘
                                   ▼
                  user_product_access(user_id, product_id)
                                   ▼
              user_video_group_access view → recording unlocked
```

Locked principles:

1. **Video groups ≠ events.** A `video_group` carries **its own** `product_id`, decoupled from any event. Do **not** reuse an event's `product_id` for its recording, and do **not** auto-unlock recordings from event payment.
2. **Two ways to gain recording access, both ending in a `user_product_access` row:**
   - **A — Self-serve purchase (fully automated):** anyone can buy a recording's product through the existing cart → Stripe → webhook path. No human in the loop.
   - **B — Manual grant:** the general rule "attending a seminar gives you the recording" is applied **manually by the client after the event** (and also covers off-platform/cash payers). There is **no automatic** event→recording wiring.
3. **Access source of truth = the `user_product_access` entitlement**, never a live Stripe payment. Recordings are gatekept (locked unless the user has the entitlement, is admin, or the group has no product).
4. **Events: some paid via Stripe, some free.** Free event ⇒ direct `event_registration` upsert. Paid event ⇒ cart/Stripe ticket. Event payment grants the **event** product and registration — it does **not** auto-grant any recording (that's path B, manual).
5. **Admin self-serve UI is in scope** (the biggest new build); the manual-grant screen is a high-priority piece of it.
6. **Full EzyCourse replacement** is the bar.

## 3. What already exists (do not rebuild)

| Area | State | Notes |
|---|---|---|
| Signup / login (password) | ✅ done | Argon2id, email-verification token flow, `email_confirmed_at` gate |
| Google OIDC | ✅ done | `oauth/google/redirect` + `/callback` |
| Sessions | ✅ done | SecureCookie, HKDF-derived key |
| Cart → Stripe checkout | ✅ done | `carts`, `cart_line_items`, `orders`, hosted checkout, signature-verified webhook |
| `user_product_access` + access view | ✅ done | Grants on paid webhook; admin = free access |
| Video model + HLS/CloudFront serving | ✅ done | `videos`, `video_sources`, `video_groups`, gating in `VideoController.Show` |
| Video access gating | ✅ wired | But **no group has a `product_id` set yet** → nothing is actually gated/sold |
| i18n PL/EN | ✅ done | `preferred_locale`, Accept-Language, bilingual `*_en`/`*_pl` columns |
| River jobs | ✅ done | order email, user email, token vacuum |
| Read-only admin | ⚠️ partial | `/admin` events list, `/admin/events/:id`, `/admin/users` — no create/edit |

## 4. Feature areas — current → target → tasks

### 4.1 Public marketing pages
- **Home** — hero + featured event banner + footer (matches legacy screenshot 1). Currently `PageController.Index`. Target: pull the featured/next upcoming event dynamically instead of hardcoding.
- **About Us** — static bilingual copy + image (screenshot 2). Low effort; ensure content editable or at least in template, not DB-critical.
- **Watch** — free/public video landing (screenshot 3). Lists `videos.is_public = true` with YouTube embeds + descriptions. Target: render from DB, ordered.
- Tasks: confirm these render from data, add the featured-event slot on Home.

### 4.2 Events (list, detail, registration)
**Current:** `/events/:slug` shows an event. Two routes hit the same `EventRegistrationController.Create` — `GET /events/:event_id/register` (router.go:128) and `POST /event_registrations/:event_id` (router.go:129); `DELETE /event_registrations/:event_id` unregisters. The GET exists so the post-login redirect (a GET) auto-registers an anonymous user: they click Register → `RequireAuthenticatedUser` redirects to `/sign-in?ref=<register url>` → after sign-in the GET bounce hits `Create` and signs them up. All behind auth. Registration is a bare upsert — no count, no email, no Zoom delivery, no paid path.

**Target:**
- **List page** (`/events`) — matches legacy: cards with date, title, type badge, "Online"/venue, registrant count, host, price ("Free" or amount), All/Upcoming/Past tabs, name search. *(New route — list page not in current routes.)*
- **Detail page** — full description, host(s), date/time in user TZ, venue or "Online", registrant count, and a single primary action button.
- **Registration flows:**
  - *Free event:* button "Register" → immediate `event_registration` row → confirmation email with venue (in-person) + `.ics` calendar attachment → page flips to "You're registered". Zoom join link deferred (see §4.9).
  - *Paid event:* button "Register" → add event product to cart → Stripe checkout → on paid webhook, create `event_registration` **and** grant `user_product_access` → confirmation email.
- **Registrant count** shown publicly (legacy shows "125 Registered"). Cache or `COUNT(*)`.
- **No capacity limits** (client decision — Zoom quota never exceeded). Don't build seat-blocking.
- **Past events** keep their recording link if a `video_group` is attached. For a passed *paid* event, the detail page's primary action becomes **"Buy recording"** instead of "Register" (see §4.3).

**Tasks:**
1. Add `GET /events` list page (filters, search, tabs).
2. Branch free vs paid in the registration handler (the GET/POST → `Create` dual-route stays — it's the intentional post-login auto-register pattern).
3. Webhook: create `event_registration` for paid event-type products (registration only — recordings are granted separately, §4.3).
4. ✅ Registration confirmation email + `.ics` + venue details (see §4.6). Zoom link deferred (§4.9).
5. Show registrant count + "you're registered" state.

### 4.3 Premium video / "My products"
**Current:** `/videos`, `/videos/:group`, `/videos/:group/:video` gated by the access view; "My products" gallery (legacy screenshot 5) lists groups with a crown on premium ones + "Watch Later".

**Target:**
- "My products" = the groups where `user_video_group_access.has_access = true` for this user (free groups + admin + purchased/granted).
- Premium groups the user does **not** own still appear (with crown / locked) and link to the buy path. Legacy shows owned + locked side by side.
- **Each `video_group` is its own product, decoupled from events** (locked §2). Two ways a user gets access, both ending in `user_product_access`:
  - **Self-serve purchase (automated):** locked group → "Buy" → existing cart → Stripe → webhook grants `user_product_access` → group unlocks. No event involved.
  - **Manual grant (admin):** the client grants the recording's product to seminar attendees after the event, or to off-platform payers (§4.4). No automatic event→recording wiring.
- A recording with no live event (pure VOD) is just a group with its own product — same two paths.
- "Watch Later" — a per-user saved list. *(New: `user_video_saves` table or reuse a lightweight join.)* Low priority; can defer.

**Tasks:**
1. Build "My products" page from the access view.
2. Show locked/premium groups with a "Buy" CTA → existing checkout.
3. (Defer) Watch Later.

### 4.4 Admin CMS — **the big build**
The client must operate without the developer. Scope it in layers so it ships incrementally.

**Layer A — Events (must-have for v1):**
- Event create/edit form: bilingual title/subtitle/description, `event_type`, `starts_at`/`ends_at`, `is_virtual`, slug (auto from title, editable), venue fields (denormalized on event), host picker.
- "Paid?" toggle → if paid, create/edit the linked `product` + a Base `product_price`. If free, no product.
- **Discount codes** (client uses them): manage `product_prices` rows with `rule_type='DiscountCode'` + `discount_code` — code, amount/percent, validity window. EarlyBird supported by schema but optional in admin v1.
- Registrants list per event (extend existing `/admin/events/:id`): names, count, paid/free/offline, export CSV.
- **Manual enrollment / offline payment** (important — many customers pay cash or bank transfer, bypassing Stripe): admin picks a user (or enters email) + product → grants `user_product_access` and, for upcoming events, the `event_registration`, **without a Stripe session**. Optionally records an order with `payment_method='cash'|'bank_transfer'|'comp'` and a manual `paid_at` so registrant lists and revenue views stay uniform. Access is granted by the entitlement, independent of payment channel.
- **Manual refund handling** (no auto-refunds): a button to mark an order refunded → revoke `user_product_access` and cancel the registration. Actual money refund done by hand in Stripe / bank.

**Manual grant / revoke recording access (HIGH PRIORITY — primary recording flow, see §2 path B):**
- Admin picks a user (or enters email) + a recording `product`/`video_group` → inserts `user_product_access(user, product)` with **no Stripe/order**. This is how the client opens locked recordings for seminar attendees and off-platform payers. Plus a revoke (delete the row).
- Optionally log an order with `payment_method` for record-keeping, but the bare grant is the must-have.

**Layer B — Recordings (high value):**
- After an event, admin attaches a recording: create/select a `video_group` with **its own** `product_id` (not the event's), add videos to it (ordered). Then grant access via the manual-grant flow above.
- **Video ingestion is heavy** (CloudFront HLS, `video_sources`, transcoding, thumbnails). Recommendation: admin UI captures **metadata + ordering**, but the actual transcode/HLS variant generation stays an **eager batch script** (consistent with the assets approach — predictable cost, dumb CDN, always-fast for a low-bandwidth cohort). Admin uploads source → script produces variants → group is publishable. Don't build on-the-fly transcoding for v1.

**Layer C — Hosts, public videos, About copy (nice-to-have):**
- Host CRUD (name, salutation, photo, country).
- Public "Watch" video CRUD (YouTube URL + bilingual title/description).
- Editable About copy.

**Tasks:** event CRUD forms → product/price sub-form → registrants view+export → recording attach UI → host CRUD → public video CRUD. (Zoom link management arrives with §4.9, not in Layer A.)

### 4.5 Asset / video ingestion pipeline
Carry over the confirmed asset model (see memory `asset-variants-decision`):
- Eager, code-defined variant sizes generated by a batch script; CloudFront stays dumb.
- For images: keep `object_key` nullable as an override (NULL ⇒ derive `{id}.{ext}`), store `extension`, and replace the `scaled` boolean with `variants_generated_at timestamptz` so re-runs are idempotent and the **original is never discarded**.
- For video: source upload → River/script transcodes to HLS renditions → populate `video_sources` (object_key, codec, priority) → generate bilingual thumbnails → mark group publishable.

### 4.6 Email & notifications
- **Registration confirmation** (free + paid): event details, venue (in-person), `.ics` attachment. New River job (or extend `SendUserEmailWorker`). Zoom join link added once §4.9 lands. ✅ Free-event confirmation done: `SendEventRegistrationEmailWorker` + `EventMailer`, bilingual PL/EN, `.ics` attached, event time shown in recipient's preferred timezone (UTC offset).
- **Payment receipt** — `SendOrderEmailWorker` exists; ensure it fires and includes the granted recording/access note.
- **Event reminder** — new River **cron** job (like `VacuumUserTokensWorker`): T-24h and/or T-1h before `starts_at`, email registrants. Carries the Zoom join link once §4.9 lands.
- **Email language follows the recipient's `preferred_locale`** (client decision), with PL as fallback. For buyers use `orders.preferred_locale`. Plain language, mobile-friendly (cohort is non-technical, on phones).

### 4.7 i18n & timezones
Already solid. Ensure every new admin/event/email string is in the bundle and that event times render in the visitor's resolved timezone (middleware exists). ✅ Browser timezone is now persisted: `users.preferred_timezone_encrypted` (bytea, application-level encryption) stores the IANA timezone name. The JS `Intl.DateTimeFormat` timezone is sent to `POST /api/v1/prefs/timezone` on first load; the session cookie is the runtime store and the DB is the persistent store. On sign-in, `maybeUpdatePreferredTimezone` syncs the session value to the DB if it differs. On subsequent sessions the DB value bootstraps the timezone before any JS runs.

**Timezone override (Profile preference):** add `preferred_timezone_locked boolean not null default false` to `users`. When `true`, auto-detection from the browser (`POST /api/v1/prefs/timezone` and `maybeUpdatePreferredTimezone`) must not overwrite the stored value — guard with `WHERE NOT preferred_timezone_locked` in `UpdateUserPreferredTimezone`. The Profile view exposes a timezone selector; saving it calls a new `SetUserTimezoneOverride` query that sets both `preferred_timezone_encrypted` and `preferred_timezone_locked = true`. Emails always use the stored value (`preferred_timezone_encrypted`), so a locked override takes full effect there.

### 4.8 Payments, tax & geography (client decisions)
- **Stripe stays as direct PSP, not a merchant-of-record.** Client explicitly does *not* want a MoR (Paddle/LemonSqueezy). MoR cost would exceed event revenue for the low-value Indian segment. Stripe-direct means **she remains the seller of record** and handles her own Polish/EU tax.
- **Stripe is not the only payment channel.** Many customers pay by **cash or bank transfer** outside the system. The app must therefore treat the `user_product_access` entitlement — not a Stripe payment — as the source of truth for access. Offline sales are entered via manual enrollment (§4.4). Bookkeeping/invoicing of offline sales is the client's responsibility and outside the app's scope; the app just records `payment_method` if she chooses to log the order.
- **Two audiences:** Polish + EU (paying, PLN) and Indian (mostly non-paying, attend free). Paid checkout effectively targets PL/EU.
  - Don't build any India-specific tax/MoR handling. Indians use free RSVP.
  - *Option:* restrict paid-checkout `billing_country` to PL/EU to avoid accidentally taking on foreign tax obligations. Low priority — flag, don't block.
- **Polish invoices (faktura):** as a Polish sole proprietorship the client likely must issue a faktura for PL buyers; a Stripe receipt alone may not satisfy PL tax. v1: keep capturing billing data + `billing_tax_id`, issue **manually** (or later integrate fakturownia.pl / wFirma). Confirm with the client — see §8.
- Currency: PLN (matches `products.base_price_currency` / `orders.currency`).

### 4.9 Zoom integration (deferred — not v1-week)
Join links should be **generated programmatically via the Zoom API**, not stored as a static shared URL.
- Approach options: Zoom **Meeting Registration API** (per-registrant unique join URL) or **Server-to-Server OAuth** app generating join tokens/ZAK. Decide based on whether per-attendee tracking is wanted.
- Flow: on registration (free or paid), call Zoom to register the attendee → store the returned per-registrant join URL against the `event_registration` → deliver in confirmation + reminder emails and reveal on event detail to registrants only.
- Replaces the dropped static `events.join_url`. Until this lands, in-person events carry venue details; online events simply omit the link.

## 5. Data model gaps to add
- ~~`events.capacity`~~ — **dropped** (no seat limits).
- ~~`events.join_url`~~ — **dropped for now**. Zoom join links to be generated programmatically via Zoom API later (§4.9), not stored as a static URL.
- `assets`: `extension`, `variants_generated_at` (replace `scaled`); keep original.
- Refund tracking: `orders.refunded_at` (revokes access on mark-refunded). Already have `cancelled_at` — decide which.
- `orders.payment_method` enum (`stripe`, `cash`, `bank_transfer`, `comp`); make `stripe_checkout_session_id` truly optional (already nullable) and allow manual `paid_at` for offline orders. Manual-grant path may create `user_product_access` with no order at all.
- (Defer) `user_video_saves` for Watch Later.
- **Every `video_group` is its own product** — give each a `product_id`, decoupled from events. Do not share an event's product. (Optionally enforce with `NOT NULL` once all groups are products — see §2 / the access-structure note.)
- `videos`: add `description_en` / `description_pl` (Watch page shows descriptions) and a nullable `youtube_id` for `provider='youtube'`, so YouTube ids stop being stuffed into `video_sources.object_key` (which is for CloudFront HLS renditions). One unified `videos` table, location stored per provider.

## 6. EzyCourse data migration
Before cancelling the SaaS, export and import:
- **Users** — in progress (encrypted seeds already added per recent commits). Reconcile against EzyCourse export; preserve who already paid → seed `user_product_access`.
- **Videos** — done (ripped to seeds).
- **Events** — historical + upcoming, with hosts and recordings.
- **Existing entitlements** — anyone who bought on EzyCourse must keep access here. Map their purchases to `user_product_access` rows.

## 7. Suggested milestones (cut line)

**M1 — Events usable (no money):** events list page, POST register, free-RSVP confirmation email + `.ics`, registrant count, "you're registered" state. *Ships the most-used legacy feature.*
- ✅ Free-RSVP confirmation email + `.ics` + timezone-aware event time (2026-06-22)

**M2 — Recordings + manual grant (primary monetization):** give each `video_group` a product, gate it, build the **admin manual grant/revoke** flow (§4.4) + "My products" locked/owned states + self-serve "Buy" → existing Stripe checkout for a recording.

**M3 — Paid events (if self-serve kept):** event→product→price wiring, buy ticket → existing Stripe flow, webhook creates registration (recording access stays manual), receipt email.

**M4 — Admin CMS Layer A + ingestion:** event/product/price CRUD + registrants view + recording-attach UI + transcode script + reminder cron. *Removes the developer as bottleneck.*

**M5 — Migration + cutover:** import entitlements, parity check on About/Watch, point homeosapiens.eu at the Go app, cancel EzyCourse.

Layer C admin (hosts, public-video, About editing) can trail after cutover.

## 8. Client decisions (resolved 2026-06-21)
1. **Capacity** — none. Unlimited seats (Zoom quota never hit). Don't build seat limits.
2. **Refunds** — no automatic refunds; handled manually. Build a manual mark-refunded → revoke-access admin action; refund the money by hand in Stripe.
3. **Standalone recordings** — yes, sell access to past recordings. Product grants recording independent of registration (§4.3).
4. **MoR / tax** — Stripe stays direct PSP, **no merchant-of-record**. PL/EU pay in PLN; Indians attend free (low value, MoR not worth it). She handles her own PL/EU tax (§4.8).
5. **Discount codes** — yes, support them in admin (EarlyBird optional).
6. **Email language** — follow recipient's `preferred_locale` (PL fallback).

7. **Recordings vs events decoupled (2026-06-22).** Video groups are standalone products. Recording access comes via (A) fully-automated self-serve purchase, or (B) manual admin grant for seminar attendees / off-platform payers. No automatic event→recording wiring. Manual grant is the primary recording flow.

### Still open
- **Faktura:** does she need automated Polish invoice generation, or is manual / Stripe receipt acceptable for v1? (§4.8)
- **Restrict paid checkout to PL/EU billing countries** to avoid foreign tax exposure — wanted, or leave open?
- **Self-serve paid *events*** — keep Stripe ticket checkout for upcoming events in v1, or are events also pay-her-directly + manual? (Decides whether the event Stripe path ships now.)
