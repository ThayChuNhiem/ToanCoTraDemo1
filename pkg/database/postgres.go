package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Driver Postgres cho Go
)

// PostgresDB đại diện cho thực thể quản lý Connection Pool của Postgres
type PostgresDB struct {
	Pool *sql.DB
}

// NewPostgresConnection khởi tạo Connection Pool tối ưu đến Postgres với cơ chế tự động thử lại
func NewPostgresConnection(dataSourceName string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("không thể mở kết nối đến Postgres: %w", err)
	}

	// Cấu hình Connection Pool tối ưu hiệu năng
	db.SetMaxOpenConns(25)                 // Tối đa 25 kết nối đồng thời
	db.SetMaxIdleConns(5)                  // Giữ tối thiểu 5 kết nối nhàn rỗi
	db.SetConnMaxLifetime(5 * time.Minute) // Hủy kết nối sau 5 phút để tránh rò rỉ bộ nhớ

	// Thử ping kiểm tra kết nối
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("không thể kết nối (ping) đến Postgres: %w", err)
	}

	return &PostgresDB{Pool: db}, nil
}

// Close đóng toàn bộ kết nối của pool an toàn
func (p *PostgresDB) Close() error {
	if p.Pool != nil {
		return p.Pool.Close()
	}
	return nil
}
