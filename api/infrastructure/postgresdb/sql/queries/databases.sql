-- name: CreateDatabase :exec
INSERT INTO "Databases"
("Id", "CreatedAt", "ModifiedAt", "IsDeleted", "IsArchived", "Name", "Type", "ConnectionString", "Schema", "LastSync")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: GetDatabases :many
SELECT * FROM "Databases"
WHERE NOT ("IsDeleted" OR "IsArchived");

-- name: GetDatabaseById :one
SELECT * FROM "Databases"
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: UpdateDatabaseById :exec
UPDATE "Databases"
SET
	"ModifiedAt" = $2,
	"IsDeleted" = $3,
	"IsArchived" = $4,
	"Name" = $5,
	"Type" = $6,
	"ConnectionString" = $7,
	"Schema" = $8,
	"LastSync" = $9
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: SoftDeleteDatabaseById :exec
UPDATE "Databases"
SET
	"ModifiedAt" = CURRENT_TIMESTAMP,
	"IsDeleted" = TRUE
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: SoftArchiveDatabaseById :exec
UPDATE "Databases"
SET
	"ModifiedAt" = CURRENT_TIMESTAMP,
	"IsArchived" = TRUE
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: DeleteDatabaseById :exec
DELETE FROM "Databases"
WHERE
	"Id" = $1;
