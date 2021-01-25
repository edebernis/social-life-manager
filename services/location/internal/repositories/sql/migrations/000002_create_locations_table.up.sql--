BEGIN;

CREATE TABLE "locations"
(
 "id"          uuid NOT NULL,
 "name"        varchar(50) NOT NULL,
 "address"     text NOT NULL,
 "category_id" uuid NOT NULL,
 "user_id"     uuid NOT NULL,
 CONSTRAINT "PK_locations" PRIMARY KEY ( "id" ),
 CONSTRAINT "FK_locations_category" FOREIGN KEY ( "category_id" ) REFERENCES "categories" ( "id" )
);

CREATE INDEX "fkIdx_locations_category" ON "locations"
(
    "category_id"
);

CREATE UNIQUE INDEX "idx_locations_name" ON "locations"
(
    "name", "user_id"
);

COMMIT;
