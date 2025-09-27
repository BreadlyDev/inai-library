ALTER TABLE book_categories 
DROP COLUMN updated_time;

ALTER TABLE book_categories 
ALTER COLUMN title VARCHAR(255) NOT NULL;