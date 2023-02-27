CREATE TABLE IF NOT EXISTS "branches"(
    "id" SERIAL PRIMARY KEY ,
    "office_id" INTEGER REFERENCES offices(id),
    "name" VARCHAR(30)
);