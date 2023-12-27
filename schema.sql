DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    author_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

INSERT INTO authors (id, name)
VALUES
    (1, 'Макс'),
    (2, 'Джон'),
    (3, 'Алиса');

INSERT INTO posts (id, author_id, author_name, title, content)
VALUES
    ( 1, 1, 'Макс', 'Средневековье', 'Какая то история про средневековье'),
    ( 2, 2, 'Джон', 'Загадки космоса', 'Рассмотрение загадок космоса и нашего места в нем'),
    ( 3, 3, 'Алиса', 'Искусство программирования', 'Творческий подход к программированию и созданию кода'),
    ( 4, 1, 'Макс', 'Гастрономические приключения', 'Путешествие в мир разнообразных вкусов и блюд'),
    ( 5, 1, 'Макс', 'Природные красоты', 'Очарование природы и ее влияние на человека'),
    ( 6, 2, 'Джон', 'Робототехника и будущее', 'Перспективы развития робототехники в современном мире');
