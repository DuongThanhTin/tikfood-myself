-- Authentication (ADR-0007), milestone M1. Rotating refresh tokens; the raw value is
-- never stored, only its hash. A token is usable when revoked_at is null and expires_at > now().
create table if not exists refresh_tokens (
  id          uuid primary key default gen_random_uuid(),
  user_id     uuid not null references users(id) on delete cascade,
  token_hash  text not null unique,
  expires_at  timestamptz not null,
  revoked_at  timestamptz,
  created_at  timestamptz not null default now(),
  user_agent  text not null default '',      -- captured for audit only
  ip          inet
);

create index if not exists refresh_tokens_user_id_idx on refresh_tokens (user_id);
