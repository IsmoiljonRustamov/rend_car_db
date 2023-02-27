CREATE TABLE IF NOT EXISTS "address"(
    "id" SERIAL PRIMARY KEY,
    "branch_id" INTEGER REFERENCES branches(id),
    "street" VARCHAR(50),
    "city" VARCHAR(50)
);