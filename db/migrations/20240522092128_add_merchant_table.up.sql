DROP TYPE IF EXISTS "MerchantCategory";
CREATE TYPE "MerchantCategory" AS ENUM (
	'SmallRestaurant',
	'MediumRestaurant',
	'LargeRestaurant',
	'MerchandiseRestaurant',
	'BoothKiosk',
	'ConvenienceStore'
);

CREATE TABLE IF NOT EXISTS "merchant" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "name" varchar(30) NOT NULL,
    "merchantCategory" "MerchantCategory" NOT NULL,
    "imageUrl" varchar(255) NOT NULL,
    "lat" REAL NOT NULL,
    "lon" REAL NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW()),
    "updatedAt" timestamp NOT NULL DEFAULT(NOW())
);
DROP TYPE IF EXISTS "ProductCategory";
CREATE TYPE "ProductCategory" AS ENUM (
    'Beverage',
    'Food',
    'Snack',
    'Condiments',
    'Additions'
);

CREATE TABLE IF NOT EXISTS "item" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "name" varchar(30) NOT NULL,
    "productCategory" "ProductCategory" NOT NULL,
    "imageUrl" varchar(255) NOT NULL,
    "price" int NOT NULL,
    "merchantId" uuid NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT(NOW()),
    "updatedAt" timestamp NOT NULL DEFAULT(NOW())
);

ALTER TABLE "item" ADD FOREIGN KEY ("merchantId") REFERENCES "merchant" ("id");