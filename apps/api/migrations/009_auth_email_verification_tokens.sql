-- Authentication follow-up: email verification. Single-use tokens sent to a user's email
-- to prove ownership. Like refresh tokens, only the hash of the raw value is stored. A
-- token is usable when consumed_at is null and expires_at > now().
create table if not exists email_verification_tokens (
  id          uuid primary key default gen_random_uuid(),
  user_id     uuid not null references users(id) on delete cascade,
  token_hash  text not null unique,
  expires_at  timestamptz not null,
  consumed_at timestamptz,                      -- null = unused; set = already verified with it
  created_at  timestamptz not null default now()
);

create index if not exists email_verification_tokens_user_id_idx on email_verification_tokens (user_id);
