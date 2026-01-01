DO $$
BEGIN
  IF to_regclass('public.instagram_tokens') IS NOT NULL THEN
    IF EXISTS (SELECT 1 FROM instagram_tokens LIMIT 1) THEN
      RAISE EXCEPTION 'instagram_tokens exists and contains data; aborting initialization.';
    END IF;
  ELSE
    EXECUTE 'CREATE TABLE instagram_tokens (
      id BOOLEAN PRIMARY KEY DEFAULT TRUE,
      access_token TEXT NOT NULL,
      expires_at TIMESTAMPTZ NOT NULL,
      updated_at TIMESTAMPTZ DEFAULT now()
    )';
  END IF;
END;
$$ LANGUAGE plpgsql;

INSERT INTO instagram_tokens (id, access_token, expires_at)
VALUES (
  TRUE,
  '<LONG_LIVED_ACCESS_TOKEN>',
  now() + interval '3 days'
);