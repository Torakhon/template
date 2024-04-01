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

INSERT INTO users (id, user_name, first_name, last_name, email, password, role, bio, website)
VALUES
    ('eecdc55a-8e19-4151-8dc5-5a8e19115189', 'Axrorbek', 'Axrorbek', 'Uzb', 'axrorbek@gmail.com', '$2a$14$PiFOuEnA4QBzyQUuCX2ReuXQxlqjAQDtIgCPM6AHck2oSV..E.HKK', 'suAdmin', 'I am a go developer', 'https://janedoe.com'),
    ('2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'jane_doe', 'Jane', 'Doe', 'jane.doe@example.com', '202cb962ac59075b964b07152d234b70', 'user', 'I am a graphic designer', 'https://janedoe.com'),
    ('3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'alice_smith', 'Alice', 'Smith', 'alice.smith@example.com', 'c0e84a40713b95ccfebec4873d3d36a1', 'admin', 'I am an administrator', 'https://alicesmith.com'),
    ('4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'bob_jones', 'Bob', 'Jones', 'bob.jones@example.com', '098f6bcd4621d373cade4e832627b4f6', 'user', 'I am a writer', 'https://bobjones.com'),
    ('5f44e820-71c1-4a34-bd3a-1d4c34b5a6a5', 'emily_brown', 'Emily', 'Brown', 'emily.brown@example.com', '5d41402abc4b2a76b9719d911017c592', 'user', 'I am a teacher', 'https://emilybrown.com'),
    ('6f44e820-71c1-4a34-bd3a-1d4c34b5a6a6', 'michael_wilson', 'Michael', 'Wilson', 'michael.wilson@example.com', 'd8578edf8458ce06fbc5bb76a58c5ca4', 'user', 'I am a doctor', 'https://michaelwilson.com'),
    ('7f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', 'sarah_johnson', 'Sarah', 'Johnson', 'sarah.johnson@example.com', 'a5bfc9e07964f8dddeb95fc584cd965d', 'user', 'I am a musician', 'https://sarahjohnson.com'),
    ('8f44e820-71c1-4a34-bd3a-1d4c34b5a6a8', 'david_martin', 'David', 'Martin', 'david.martin@example.com', 'd1fe173d5f31199baa62c78a67df98d1', 'user', 'I am an engineer', 'https://davidmartin.com'),
    ('9f44e820-71c1-4a34-bd3a-1d4c34b5a6a9', 'laura_garcia', 'Laura', 'Garcia', 'laura.garcia@example.com', 'd1f41c45bb260d867f091c25e313d3f9', 'user', 'I am an artist', 'https://lauragarcia.com'),
    ('af44e820-71c1-4a34-bd3a-1d4c34b5a6aa', 'james_anderson', 'James', 'Anderson', 'james.anderson@example.com', 'd5e14d37355a7b90d02aa4e64e42aee7', 'user', 'I am a photographer', 'https://jamesanderson.com'),
    ('bf44e820-71c1-4a34-bd3a-1d4c34b5a6ab', 'linda_taylor', 'Linda', 'Taylor', 'linda.taylor@example.com', 'e99a18c428cb38d5f260853678922e03', 'user', 'I am a chef', 'https://lindataylor.com');


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

INSERT INTO posts (id, title, content, user_id, category)
VALUES
--      Jane
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Post Title 1', 'Post Content 1', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Category 1'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Post Title 2', 'Post Content 2', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Category 2'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Post Title 3', 'Post Content 3', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Category 1'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a5', 'Post Title 4', 'Post Content 4', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Category 2'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a6', 'Post Title 5', 'Post Content 5', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Category 1'),

--     Alice
    ('2f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', 'Post Title 6', 'Post Content 6', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Category 1'),
    ('2f44e820-71c1-4a34-bd3a-1d4c34b5a6a8', 'Post Title 7', 'Post Content 7', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Category 2'),
    ('2f44e820-71c1-4a34-bd3a-1d4c34b5a6a9', 'Post Title 8', 'Post Content 8', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Category 1'),
    ('300f8f45-896d-4266-8f8f-45896d626639', 'Post Title 9', 'Post Content 9', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Category 2'),
    ('fbd62e58-1fa7-461d-962e-581fa7f61de2', 'Post Title 10', 'Post Content 10', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Category 1'),

--     Bob
    ('4f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', 'Post Title 11', 'Post Content 11', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Category 1'),
    ('4f44e820-71c1-4a34-bd3a-1d4c34b5a6a8', 'Post Title 12', 'Post Content 12', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Category 2'),
    ('4f44e820-71c1-4a34-bd3a-1d4c34b5a6a9', 'Post Title 13', 'Post Content 13', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Category 1'),
    ('c9bb883e-118a-4e0a-bb88-3e118a3e0a8f', 'Post Title 14', 'Post Content 14', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Category 2'),
    ('0f225560-bf5e-4fb9-a255-60bf5e2fb962', 'Post Title 15', 'Post Content 15', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Category 1');



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

INSERT INTO comments (comment_id, post_id, user_id, content)
VALUES
--      Jane
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a1', '1f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Bu birinchi mock izoh'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', '1f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Bu ikkinchi mock izoh'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', '1f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a2', 'Bu uchinchi mock izoh'),

--      Alice
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Bu birinchi mock izoh'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a8', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Bu ikkinchi mock izoh'),
    ('1f44e820-71c1-4a34-bd3a-1d4c34b5a6a9', '2f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '3f44e820-71c1-4a34-bd3a-1d4c34b5a6a3', 'Bu uchinchi mock izoh'),

--     Bob
    ('5ea898b4-da02-450f-a898-b4da02450f20', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Bu birinchi mock izoh'),
    ('56fc7bb4-3a8b-4640-bc7b-b43a8b4640e4', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Bu ikkinchi mock izoh'),
    ('00dd2492-4820-4ede-9d24-9248205ede88', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a7', '4f44e820-71c1-4a34-bd3a-1d4c34b5a6a4', 'Bu uchinchi mock izoh');


CREATE TABLE IF NOT EXISTS comment_like (
    comment_id UUID UNIQUE NOT NULl,
    user_id UUID UNIQUE NOT NULl,
    status BOOL
);