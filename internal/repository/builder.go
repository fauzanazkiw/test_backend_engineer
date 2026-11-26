package repository

import (
	"github.com/Masterminds/squirrel"
)

// builder dipakai bersama agar seluruh query memakai placeholder PostgreSQL ($1, $2, ...).
var builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
