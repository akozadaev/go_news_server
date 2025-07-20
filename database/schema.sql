-- Создание таблицы новостей
CREATE TABLE IF NOT EXISTS "News" (
    "Id" BIGSERIAL PRIMARY KEY,
    "Title" VARCHAR(255) NOT NULL,
    "Content" TEXT NOT NULL
);

-- Создание таблицы категорий
CREATE TABLE IF NOT EXISTS "Categories" (
    "Id" BIGSERIAL PRIMARY KEY,
    "Name" VARCHAR(100) NOT NULL UNIQUE,
    "Description" TEXT,
    "CreatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы связи новостей и категорий
CREATE TABLE IF NOT EXISTS "NewsCategories" (
    "NewsId" BIGINT NOT NULL,
    "CategoryId" BIGINT NOT NULL,
    PRIMARY KEY ("NewsId", "CategoryId"),
    FOREIGN KEY ("NewsId") REFERENCES "News"("Id") ON DELETE CASCADE,
    FOREIGN KEY ("CategoryId") REFERENCES "Categories"("Id") ON DELETE CASCADE
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_news_title ON "News"("Title");
CREATE INDEX IF NOT EXISTS idx_categories_name ON "Categories"("Name");
CREATE INDEX IF NOT EXISTS idx_news_categories_news_id ON "NewsCategories"("NewsId");
CREATE INDEX IF NOT EXISTS idx_news_categories_category_id ON "NewsCategories"("CategoryId");

-- Вставка тестовых данных для категорий
INSERT INTO "Categories" ("Name", "Description") VALUES
    ('Technology', 'Technology related news'),
    ('Sports', 'Sports news and updates'),
    ('Politics', 'Political news and analysis'),
    ('Entertainment', 'Entertainment and celebrity news'),
    ('Business', 'Business and economic news')
ON CONFLICT ("Name") DO NOTHING; 