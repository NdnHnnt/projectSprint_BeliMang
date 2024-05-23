CREATE TABLE IF NOT EXISTS "user" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "username" varchar(30) UNIQUE NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "password" varchar(30) NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW()),
    "updatedAt" timestamp NOT NULL DEFAULT(NOW())
);
