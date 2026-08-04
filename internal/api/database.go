package api

import "context"

type Database interface {
	// Ping checks the database connection.
	Ping(ctx context.Context) error
}
