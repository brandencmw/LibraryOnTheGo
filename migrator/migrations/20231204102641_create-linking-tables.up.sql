CREATE TABLE IF NOT EXISTS book_authors (
    book_id INTEGER REFERENCES books(id),
    author_id INTEGER REFERENCES authors(id),
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE IF NOT EXISTS book_categories (
    book_id INTEGER REFERENCES books(id),
    category_id INTEGER REFERENCES categories(id),
    PRIMARY KEY (book_id, category_id)
);