air: mise exec -- env PORT=5000 air
assets: cd assets && pnpm run dev
caddy: caddy run
stripe: stripe listen --forward-to localhost:5000/webhooks/stripe --events $STRIPE_LISTEN_TO_EVENTS
