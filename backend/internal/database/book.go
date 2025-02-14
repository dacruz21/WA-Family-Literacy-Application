package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/TritonSE/words-alive/internal/models"
)

type BookDatabase struct {
	Conn *pgxpool.Pool
}

/*
 * Get book list for the main page
 * Only need to show previews (no read/explore/learn)
 */
func (db *BookDatabase) FetchBookList(ctx context.Context) ([]models.Book, error) {

	books := make([]models.Book, 0)

	var query string = "SELECT books.id, title, author, image, created_at, " +
		"array_remove(array_agg(lang), NULL) as languages " +
		"FROM books LEFT JOIN book_contents ON books.id = " +
		"book_contents.id GROUP BY books.id ORDER BY title"

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error on SELECT FROM books in FetchBookList")
	}

	defer rows.Close()

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Image,
			&book.CreatedAt, &book.Languages); err != nil {
			fmt.Print(err)
			return nil, errors.Wrap(err, "error scanning result of"+
				" SELECT FROM books in FetchBookList")
		}

		books = append(books, book)
	}

	return books, nil
}

/*
 * Fetch a full book including all read/explore/learn content
 * For use after a user has selected a book to read
 */
func (db *BookDatabase) FetchBookDetails(ctx context.Context,
	id string, lang string) (*models.BookDetails, bool, error) {

	var book models.BookDetails

	var query string = "SELECT books.id, title, author, image, " +
		"created_at, read_video, read_body, explore_video, explore_body, " +
		"learn_video, learn_body FROM books LEFT JOIN book_contents ON " +
		"books.id = book_contents.id WHERE books.id = $1 AND lang = $2"
	rows, err := db.Conn.Query(ctx, query, id, lang)
	if err != nil {
		return &book, false, errors.Wrap(err, "error on query for book details")
	}

	defer rows.Close()

	if !rows.Next() {
		// Check if user has selected a wrong language
		var count int
		var queryID string = "SELECT count(*) FROM book_contents WHERE " +
			"id = $1"

		err = db.Conn.QueryRow(ctx, queryID, id).Scan(&count)

		if err != nil {
			return nil, false, errors.Wrap(err, "error on scan for ID")
		}

		return nil, count > 0, nil
	}

	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Image,
		&book.CreatedAt, &book.Read.Video, &book.Read.Body, &book.Explore.Video,
		&book.Explore.Body, &book.Learn.Video, &book.Learn.Body)

	if err != nil {
		return &book, false, errors.Wrap(err, "error on Scan into book details")
	}

	return &book, false, nil
}

/*
 * Fetches a book entry from the database given the id
 */
func (db *BookDatabase) FetchBook(ctx context.Context, id string) (*models.Book, error) {
	var book models.Book
	var query = "SELECT books.id, title, author, image, array_remove(array_agg(lang), NULL) as languages," +
		"created_at FROM books LEFT JOIN book_contents ON " +
		"books.id = book_contents.id WHERE books.id = $1 GROUP BY books.id"

	err := db.Conn.QueryRow(ctx, query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Image,
		&book.Languages, &book.CreatedAt)

	if err != nil {
		return nil, errors.Wrap(err, "error on FetchBook")
	}

	return &book, nil
}

/*
 * Inserts a book into the books table
 */
func (db *BookDatabase) InsertBook(ctx context.Context,
	book models.APICreateBook) (*models.Book, error) {
	var newBookId string

	var query string = "INSERT INTO books (title, author, image) " +
		"VALUES ($1, $2, $3) RETURNING id"
	err := db.Conn.QueryRow(ctx, query, book.Title, book.Author, book.Image).Scan(&newBookId)

	if err != nil {
		return nil, errors.Wrap(err, "error on INSERT INTO books in InsertBook")
	}

	return db.FetchBook(ctx, newBookId)

}

/*
 * Inserts a books details into the book_contents table. Returns the
 * complete books details
 */
func (db *BookDatabase) InsertBookDetails(ctx context.Context, id string,
	book models.APICreateBookContents) (*models.BookDetails, error) {
	var newBookDetail *models.BookDetails
	var query string = "INSERT INTO book_contents " +
		"(id, lang, read_video, read_body, explore_video, explore_body, " +
		"learn_video, learn_body) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := db.Conn.Exec(ctx, query, id, book.Language,
		book.Read.Video, book.Read.Body, book.Explore.Video, book.Explore.Body,
		book.Learn.Video, book.Learn.Body)

	if err != nil {
		return nil, errors.Wrap(err, "error on INSERT INTO book_contents in InsertBookDetails")
	}

	newBookDetail, _, err = db.FetchBookDetails(ctx, id, book.Language)

	return newBookDetail, err
}

/*
 * Deletes a language entry of a book in the book_contents table.
 */
func (db *BookDatabase) DeleteBookContent(ctx context.Context, id string, lang string) error {
	var query string = "DELETE from book_contents WHERE id = $1 AND lang = $2"

	commandTag, err := db.Conn.Exec(ctx, query, id, lang)

	if err != nil {
		return errors.Wrap(err, "error on delete from book_contents")
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}

	return nil
}

/*
 * Deletes a book from the books table
 */
func (db *BookDatabase) DeleteBook(ctx context.Context, id string) error {
	var query string = "DELETE from books WHERE id = $1"

	commandTag, err := db.Conn.Exec(ctx, query, id)

	if err != nil {
		return errors.Wrap(err, "error on delete from book")
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}

	return nil
}

/*
 * Updates a book in the books table
 */
func (db *BookDatabase) UpdateBook(ctx context.Context, id string,
	updates models.APIUpdateBook) (*models.Book, error) {
	var query string = "UPDATE books " +
		"SET title = COALESCE($1, title), " +
		"author = COALESCE($2, author), " +
		"image = COALESCE($3, image) " +
		"WHERE id = $4"
	_, err := db.Conn.Exec(ctx, query, updates.Title, updates.Author,
		updates.Image, id)
	if err != nil {
		return nil, errors.Wrap(err, "error on update book")
	}

	return db.FetchBook(ctx, id)

}

/*
 * Updates a row in the book_contents table
 */
func (db *BookDatabase) UpdateBookDetails(ctx context.Context, id string,
	lang string, book models.APIUpdateBookDetails) (*models.BookDetails, error) {
	var updatedBookDetails *models.BookDetails
	var query string = "UPDATE book_contents SET " +
		"read_video = COALESCE($1, read_video), " +
		"read_body = COALESCE($2, read_body), " +
		"explore_video = COALESCE($3, explore_video), " +
		"explore_body = COALESCE($4, explore_body), " +
		"learn_video = COALESCE($5, learn_video), " +
		"learn_body = COALESCE($6, learn_body) " +
		"WHERE id = $7 AND lang = $8"

	_, err := db.Conn.Exec(ctx, query,
		book.Read.Video, book.Read.Body,
		book.Explore.Video, book.Explore.Body,
		book.Learn.Video, book.Learn.Body, id, lang)

	if err != nil {
		fmt.Print(err)
		return nil, errors.Wrap(err, "error on updating book_contents")
	}

	updatedBookDetails, _, err = db.FetchBookDetails(ctx, id, lang)

	return updatedBookDetails, err

}
