CREATE TABLE IF NOT EXISTS "Databases" (
	-- Meta
	"Id" UUID NOT NULL DEFAULT(gen_random_uuid()),
	"CreatedAt" timestamp NOT NULL,
	"ModifiedAt" timestamp NOT NULL,
	"IsArchived" boolean NOT NULL DEFAULT(FALSE),
	"IsDeleted" boolean NOT NULL DEFAULT(FALSE),

	-- Fields
	"Name" TEXT NOT NULL,
	"Type" INTEGER NOT NULL,
	"ConnectionString" TEXT NOT NULL,
	"Schema" JSONB,
	"LastSync" TIMESTAMP,

	PRIMARY KEY("Id")
);

CREATE TABLE IF NOT EXISTS "Charts" (
	-- Meta
	"Id" UUID NOT NULL DEFAULT(gen_random_uuid()),
	"CreatedAt" timestamp NOT NULL,
	"ModifiedAt" timestamp NOT NULL,
	"IsArchived" boolean NOT NULL DEFAULT(FALSE),
	"IsDeleted" boolean NOT NULL DEFAULT(FALSE),

	-- BaseChart
	"Title" TEXT NOT NULL,
	"Type" TEXT NOT NULL,
	"Query" TEXT NOT NULL,
	"DatabaseId" UUID NOT NULL,

	-- Bar/Linear Fields
	"CategoricalColumnName" TEXT,
	"ValueColumnName" TEXT,

	PRIMARY KEY("Id"),
	CONSTRAINT "fk_databases_charts"
		FOREIGN KEY("DatabaseId") 
			REFERENCES "Databases"("Id")
);
