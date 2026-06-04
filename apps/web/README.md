# TikFood Web

Next.js App Router frontend skeleton for TikFood discovery.

Implemented:

- Discovery landing page
- Interactive discovery filters
- Collapsible left discovery panel
- Dark glass discovery UI inspired by the homepage prompt
- MapLibre restaurant map with clickable dark-mode markers
- Right-side venue detail panel with hero media and trend/video sections
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
