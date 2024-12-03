-- +goose Up
CREATE TABLE "authors" (
    "id" INTEGER PRIMARY KEY,
    "handle" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "bio" TEXT NOT NULL,
    "birth" DATE NOT NULL,
    "image_url" TEXT NOT NULL,
    "github_url" TEXT NOT NULL
);
CREATE UNIQUE INDEX "idx_authors_handle" ON "authors"("handle");

INSERT INTO "authors" ("handle", "name", "bio", "birth", "image_url", "github_url") VALUES
    ('ffss', 'Felipe dos Santos', 'I like coding in Go', '1992-04-27', 'https://avatars.githubusercontent.com/u/64739815', 'https://github.com/ffss92');


CREATE VIRTUAL TABLE "blog_posts_fts" USING fts5 (
    "slug",
    "title",
    "subtitle",
    "content",
    tokenize = 'porter'
);

-- +goose Down
DROP TABLE "authors";
DROP TABLE "blog_posts_fts";
