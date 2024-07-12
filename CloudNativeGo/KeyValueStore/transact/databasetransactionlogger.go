package transact

import (
	"KeyValueStore/core"
	"database/sql"
	"fmt"
)

type PostgresTransactionLogger struct {
	db     *sql.DB
	events chan<- core.Event // write only channel
	errors <-chan error      // read only channel
}

type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", config.host, config.dbName, config.user, config.password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}
	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_name = transactions
	)`

	exists := false
	err := l.db.QueryRow(query).Scan(&exists)
	return exists, err
}

func (l *PostgresTransactionLogger) createTable() error {
	query := "CREATE TABLE transactions (sequence SERIAL PRIMARY KEY, event_type INT, key VARCHAR(255), value VARCHAR(255))"
	_, err := l.db.Exec(query)
	return err
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions (event_type, key, value) VALUES ($1, $2, $3)`
		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outErrors := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outErrors)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outErrors <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var event core.Event
			if err := rows.Scan(&event.Sequence, &event.EventType, &event.Key, &event.Value); err != nil {
				outErrors <- fmt.Errorf("error reading rows: %w", err)
				return
			}

			outEvent <- event
		}
		if err = rows.Err(); err != nil {
			outErrors <- fmt.Errorf("transaction load read failure: %w", err)
			return
		}
	}()

	return outEvent, outErrors
}

func (l *PostgresTransactionLogger) Close() {
	l.db.Close()
}

func (l *PostgresTransactionLogger) LogPut(key string, value string) {
	ev := core.Event{
		EventType: core.PUT,
		Key:       key,
		Value:     value,
	}
	l.events <- ev
}
func (l *PostgresTransactionLogger) LogDelete(key string) {
	ev := core.Event{
		EventType: core.DELETE,
		Key:       key,
	}
	l.events <- ev
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}
