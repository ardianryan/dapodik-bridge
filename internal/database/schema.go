package database

import (
	"context"
	"fmt"

	"github.com/ardianryan/dapodik-bridge/internal/models"
)

// TableExists checks whether a specific table exists in the database
func (m *DBManager) TableExists(ctx context.Context, tableName string) (bool, error) {
	if m.Pool() == nil {
		return false, fmt.Errorf("database connection is not available")
	}

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = $1 AND table_schema NOT IN ('pg_catalog', 'information_schema')
		);
	`
	err := m.Pool().QueryRow(ctx, query, tableName).Scan(&exists)
	return exists, err
}

// GetSchemaTables lists all tables in the Dapodik database
func (m *DBManager) GetSchemaTables(ctx context.Context) ([]models.TableInfo, error) {
	if m.Pool() == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	query := `
		SELECT table_schema, table_name, table_type
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name;
	`
	rows, err := m.Pool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.TableInfo
	for rows.Next() {
		var t models.TableInfo
		if err := rows.Scan(&t.Schema, &t.TableName, &t.TableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

// GetSchemaColumns lists columns for a specific table
func (m *DBManager) GetSchemaColumns(ctx context.Context, tableName string) ([]models.ColumnInfo, error) {
	if m.Pool() == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position;
	`
	rows, err := m.Pool().Query(ctx, query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []models.ColumnInfo
	for rows.Next() {
		var c models.ColumnInfo
		if err := rows.Scan(&c.ColumnName, &c.DataType, &c.IsNullable); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}
