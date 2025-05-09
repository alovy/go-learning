CREATE TABLE IF NOT EXISTS url_responses (
  id        SERIAL PRIMARY KEY,
  url       TEXT NOT NULL,
  response  TEXT NOT NULL
);
