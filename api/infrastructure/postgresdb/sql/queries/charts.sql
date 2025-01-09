-- name: CreateChart :exec
INSERT INTO "Charts"
("Id", "CreatedAt", "ModifiedAt", "IsDeleted", "IsArchived", "DatabaseId", "Title", "Type", "Query", "CategoricalColumnName", "ValueColumnName")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);

-- name: GetCharts :many
SELECT * FROM "Charts"
WHERE NOT ("IsDeleted" OR "IsArchived");

-- name: GetChartById :one
SELECT * FROM "Charts"
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: UpdateChartById :exec
UPDATE "Charts"
SET
	"ModifiedAt" = $2,
	"IsDeleted" = $3,
	"IsArchived" = $4,
	"Title" = $5,
	"Type" = $6,
	"Query" = $7,
	"CategoricalColumnName" = $8,
	"ValueColumnName" = $9
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: SoftDeleteChartById :exec
UPDATE "Charts"
SET
	"ModifiedAt" = CURRENT_TIMESTAMP,
	"IsDeleted" = TRUE
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: SoftArchiveChartById :exec
UPDATE "Charts"
SET
	"ModifiedAt" = CURRENT_TIMESTAMP,
	"IsArchived" = TRUE
WHERE
	"Id" = $1
	AND NOT ("IsDeleted" OR "IsArchived");

-- name: DeleteChartById :exec
DELETE FROM "Charts"
WHERE
	"Id" = $1;
