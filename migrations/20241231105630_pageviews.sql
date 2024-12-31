-- +goose Up
CREATE TABLE "pageviews" (
    "id" INTEGER PRIMARY KEY,
    "slug" TEXT NOT NULL,
    "ip_address" TEXT NOT NULL,
    "user_agent" TEXT NOT NULL,
    "referrer" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE "pageviews";
