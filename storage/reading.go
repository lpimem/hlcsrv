package storage

import (
	"database/sql"
	"github.com/lpimem/hlcsrv/util"
	"time"
)

type reading struct{}

// Reading API namespace
var Reading reading

// Book a book
type Book struct {
	Title     string
	BasePath  string
	PageCount int
	BookAdded time.Time
}

// ReadingProgress Book reading progress of a user
type ReadingProgress struct {
	Book
	Progress int
	Display  bool
	NextPage string

	StartRead time.Time
	LastRead  time.Time
}

// QueryBooks -
func (r *reading) QueryBooks(titleCond string) (books []*Book, err error) {
	query := `select title, base_path, page_count, ctime from book
	where title like ?`
	err = util.QueryDb(storage.DB, query, []interface{}{titleCond},
		func(idx int, rows *sql.Rows) error {
			var b Book
			err := rows.Scan(&b.Title, &b.BasePath, &b.PageCount, &b.BookAdded)
			if err != nil {
				return err
			}
			books = append(books, &b)
			return nil
		})
	return books, err
}

// QueryAllProgress -
func (r *reading) QueryAllProgress(uid UserID) (ret []*ReadingProgress, err error) {
	query := `select title, base_path, page_count, b.progress, b.display, 
	case when b.page_uri is null then toc_page else b.page_uri end as next_page, 
	ctime as book_added, b.ctime as start_read, b.mtime as last_read, 
	from book
	left join reading b on book.id = b.bid
	where b.uid = ?`
	err = util.QueryDb(storage.DB, query, []interface{}{uid},
		func(idx int, rows *sql.Rows) error {
			var p ReadingProgress
			var display int
			err := rows.Scan(&p.Title, &p.BasePath, &p.PageCount, &p.Progress,
				&display, &p.NextPage, &p.BookAdded, &p.StartRead, &p.LastRead)
			if err != nil {
				return err
			}
			p.Display = display > 0
			ret = append(ret, &p)
			return nil
		})
	return ret, err
}
