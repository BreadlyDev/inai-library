CREATE TABLE IF NOT EXISTS book_categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NULL 
);

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    pass_hash VARCHAR(255) NOT NULL,
    access_level INT NOT NULL DEFAULT 50,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books(
    id SERIAL PRIMARY KEY, 
    title VARCHAR(255) NOT NULL,
    description TEXT NULL, 
    image_path TEXT NULL, 
    file_path TEXT NULL,
    category_id INT NOT NULL,
    language VARCHAR(255) NOT NULL,
    edition_year SMALLINT NOT NULL, 
    added_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP NULL, 

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES book_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS book_copies(
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    inventory_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available', -- available, borrowed, reserved, lost, repair

    CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES books(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS reservations(
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,    
    book_id INT NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NULL,
    returned_date TIMESTAMP NULL, 

    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES books(id)
    ON DELETE SET NULL ON UPDATE CASCADE 
);

CREATE TABLE IF NOT EXISTS reviews(
    id SERIAL PRIMARY KEY, 
    author_id UUID NOT NULL,
    rating INT NULL CHECK(rating >= 0 AND rating <= 5),
    book_id INT NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (author_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES books(id)
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications(
    id SERIAL PRIMARY KEY, 
    title VARCHAR(255) NULL,
    message TEXT NOT NULL, 
    recipient_id UUID NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (recipient_id) REFERENCES users(id) 
    ON DELETE SET NULL ON UPDATE CASCADE
); 