# TikFood Web

Next.js App Router frontend skeleton for TikFood discovery.

Implemented:

- Discovery landing page
- Interactive discovery filters
- Collapsible left discovery panel
- MapLibre restaurant map with clickable markers
- Right-side venue detail panel
- Venue cards with AI summaries and trending dishes
- API integration with `GET /api/v1/discovery/venues`
- Fallback data when the API is not running

This app does not implement delivery, cart, order, payment, booking, reservations, chat, monetization, or livestream logic.

Run locally:

```bash
npm run dev
```

Optional API URL:

```bash
NEXT_PUBLIC_API_URL=http://localhost:18081
```
