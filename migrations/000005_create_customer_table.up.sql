CREATE TABLE IF NOT EXISTS "customer"(
    "id" SERIAL PRIMARY KEY ,
    "car_id" INTEGER REFERENCES cars(id),
    "name" VARCHAR(30),
    "age" INTEGER,
    "phone_number" VARCHAR(30),
    "address"  VARCHAR(50)
);