# TODO — week of 2026-06-21

Goal: land **M1 (free events usable end-to-end)** and start **M2 (paid wiring)**. See `docs/SPEC.md` for full context.

## M1 — free events
- [ ] **1. Convert event register GET → POST.** `GET /events/:event_id/register` mutates state — crawlers/prefetch can register users. Move to POST + form/CSRF, update links. Small, security fix, do first (unblocks 2/3/5).
- [ ] **2. Build `/events` list page.** New public route. Cards: date, title, type badge, Online/venue, registrant count, host, price (Free/amount). All/Upcoming/Past tabs + name search. Order by `starts_at`.
- [ ] **3. Registrant count + "you're registered" state.** `COUNT(*)` per event on list + detail. On detail, flip primary button to registered state if current user already registered. (No join link yet — Zoom deferred.)
- [ ] **5. Free-RSVP confirmation email + `.ics`.** River job (or extend `SendUserEmailWorker`): event details + venue (in-person) + `.ics`. Localized to recipient `preferred_locale` (PL fallback). Wire into POST register handler.
- [ ] **6. Verify M1 end-to-end.** Browse `/events` → detail → register free event → confirm row + email + `.ics` + registered state + count. Closes M1.

## M2 — paid wiring (start)
- [ ] **7. Wire paid-event webhook → registration.** On `checkout.session.completed` for event-type product on an *upcoming* event, also create `event_registration` (plus existing `user_product_access` grant). Don't couple access to registration for passed/VOD products. Add buy-button → cart → existing checkout.
- [ ] **8. Set real `product_id` on one video group.** Activates existing access view + `VideoController` gating. Confirm non-owner blocked, owner can watch. Proves premium path. (Independent — good filler if blocked on Stripe test keys.)

## Backlog (not this week)
- [ ] **9. Zoom API join links.** Per-registrant join URLs via Zoom API (Meeting Registration API or Server-to-Server OAuth). Deliver in confirmation + reminder emails, reveal to registrants. Replaces dropped static `join_url`. Spec §4.9.

---
Sequencing: 1 → (2, 3) → 5 → 6 → 7. 8 anytime. Admin CMS (M3) deliberately deferred until free+paid flows proven.
