CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULl PRIMARY KEY,
    user_name TEXT UNIQUE NOT NULL ,
    first_name TEXT NOT NULL ,
    last_name TEXT NOT NULL ,
    email TEXT UNIQUE NOT NULL ,
    password TEXT NOT NULL ,
    role VARCHAR(50),
    bio TEXT,
    website TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id UUID NOT NULl PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id UUID  NOT NULl ,
    category TEXT NOT NULL,
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS views(
    user_id UUID NOT NULL ,
    post_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS post_like (
    dislike TEXT ,
    post_id UUID  NOT NULl,
    user_id UUID  NOT NULl,
    status BOOL
);

CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID UNIQUE NOT NULl,
    post_id UUID  NOT NULl,
    user_id UUID  NOT NULl,
    content TEXT NOT NULL,
    likes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comment_like (
    comment_id UUID UNIQUE NOT NULl,
    user_id UUID UNIQUE NOT NULl,
    status BOOL
);