ALTER TABLE book_categories
ALTER COLUMN title UNIQUE;

ALTER TABLE book_categories 
DROP COLUMN updated_time;