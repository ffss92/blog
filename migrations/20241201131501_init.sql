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

CREATE VIRTUAL TABLE "blog_posts_fts" USING fts5 (
    "slug",
    "title",
    "content",
    tokenize = 'porter'
);

-- +goose Down
DROP TABLE "authors";
DROP TABLE "blog_posts_fts";
