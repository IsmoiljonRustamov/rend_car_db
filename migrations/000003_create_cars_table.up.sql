CREATE TABLE IF NOT EXISTS "cars"(
    "id" SERIAL PRIMARY KEY,
    "branch_id" INTEGER REFERENCES branches(id),
    "name" VARCHAR(30),
    "color" VARCHAR(30),
    "cost_day" NUMERIC(8,2),
    "amount" INTEGER
);