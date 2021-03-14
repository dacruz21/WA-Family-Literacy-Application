package controllers_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TritonSE/words-alive/internal/controllers"
	"github.com/TritonSE/words-alive/internal/database"
	"github.com/TritonSE/words-alive/internal/models"
	"github.com/TritonSE/words-alive/internal/testutils"
)

var conn = database.GetConnection()
var r = controllers.GetRouter()
var ts = httptest.NewServer(r)

// Set up database
func TestMain(m *testing.M) {
	database.Migrate("../../migrations")

	_, _ = conn.Exec("DELETE FROM books")
	_, _ = conn.Exec("DELETE FROM book_contents")

	// Populate database
	conn.Exec("INSERT INTO books (id, title, author) values ('c_id', 'c','c1');")
	conn.Exec("INSERT INTO books (id, title, author) values ('a_id', 'a','a1');")
	conn.Exec("INSERT INTO books (id, title, author) values ('b_id', 'b','b1');")
	conn.Exec("INSERT INTO books (id, title, author) VALUES ('d_id'," +
		"'d', 'd1');")
	conn.Exec("INSERT INTO books (id, title, author) VALUES ('catcher'," +
		"'catcher in the rye', 'a');")
	conn.Exec("INSERT INTO books (id, title, author) VALUES ('update_me'," +
		"'not_updated', 'update_author');")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('a_id', 'en', 'a_en_rv', 'a_en_rb', 'a_en_ev', 'a_en_eb', " +
		"'a_en_lv', 'a_en_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('b_id', 'en', 'b_en_rv', 'b_en_rb', 'b_en_ev', 'b_en_eb', " +
		"'b_en_lv', 'b_en_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('d_id', 'en', 'd_en_rv', 'd_en_rb', 'd_en_ev', 'd_en_eb', " +
		"'d_en_lv', 'd_en_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('c_id', 'en', 'c_en_rv', 'c_en_rb', 'c_en_ev', 'c_en_eb', " +
		"'c_en_lv', 'c_en_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('a_id', 'es', 'a_es_rv', 'a_es_rb', 'a_es_ev', 'a_es_eb', " +
		"'a_es_lv', 'a_es_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('catcher', 'en', 'catcher_rv', 'catcher_rb', 'catcher_ev', " +
		" 'catcher_eb', 'catcher_lv', 'catcher_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('catcher', 'es', 'catcher_es_rv', 'catcher_es_rb', 'catcher_es_ev', " +
		" 'catcher_es_eb', 'catcher_es_lv', 'catcher_es_lb')")

	conn.Exec("INSERT INTO book_contents (id, lang, read_video, read_body, " +
		"explore_video, explore_body, learn_video, learn_body) VALUES " +
		"('update_me', 'en', 'updt_en_rv', 'updt_en_rb', 'updt_en_ev', " +
		" 'updt_en_eb', 'updt_en_lv', 'updt_en_lb')")

	// Close the server
	defer ts.Close()
	os.Exit(m.Run())
}

// Test for the book list function
func TestGetBooks(t *testing.T) {

	var response []models.Book
	testutils.MakeHttpRequest("GET", ts.URL+"/books", nil, &response, t)

	// Check for correct number of elements, and sort alphabetically
	require.Len(t, response, 6)
	require.Equal(t, "a", response[0].Title)
	require.Equal(t, "b", response[1].Title)
	require.Equal(t, "c", response[2].Title)

}

// Test for the book details function
func TestGetBookDetails(t *testing.T) {

	var response models.BookDetails
	testutils.MakeHttpRequest("GET", ts.URL+"/books/catcher/en", nil, &response, t)

	// Check title and content
	require.Equal(t, "catcher in the rye", response.Title)
	require.Equal(t, "catcher_rb", response.Read.Body)
}

// Test for non-existent language
func TestGetBookNullLang(t *testing.T) {

	var response string
	testutils.MakeHttpRequest("GET", ts.URL+"/books/catcher/fr", nil, &response, t)

	require.Equal(t, "book does not exist in specified language", response)
}

// Test for non-existent book
func TestGetNullBook(t *testing.T) {

	var response string
	testutils.MakeHttpRequest("GET", ts.URL+"/books/nonexistent/en", nil, &response, t)

	require.Equal(t, "book not found", response)
}

// Test if creating a book and its corresponding contents works
func TestCreateBookandBookDetails(t *testing.T) {
	var book = models.APICreateBook{
		Title:  "Harry Potter",
		Author: "JK Rowling",
		Image:  nil,
	}

	var bookDetails = models.APICreateBookContents{
		Language: "en",
		Read: models.TabContent{
			Video: nil,
			Body:  "read_body",
		},
		Explore: models.TabContent{
			Video: nil,
			Body:  "explore_body",
		},
		Learn: models.TabContent{
			Video: nil,
			Body:  "learn_body",
		},
	}

	var createdBook models.Book
	var createdBookDetails models.BookDetails
	var jsonString = testutils.MakeJSONBody(book, t)

	testutils.MakeHttpRequest("POST", ts.URL+"/books", jsonString, &createdBook, t)

	require.Equal(t, "Harry Potter", createdBook.Title)
	require.Equal(t, "JK Rowling", createdBook.Author)
	require.Equal(t, []string{}, createdBook.Languages)

	jsonString = testutils.MakeJSONBody(bookDetails, t)

	testutils.MakeHttpRequest("POST",
		ts.URL+"/books/"+createdBook.ID, jsonString, &createdBookDetails, t)

}

// Test if inserting into book_contents with an id not in books throws an error
func TestBookDetailsErrorWithInvalidID(t *testing.T) {
	var badBookDetails = models.APICreateBookContents{

		Language: "en",
		Read: models.TabContent{
			Video: nil,
			Body:  "read_body",
		},
		Explore: models.TabContent{
			Video: nil,
			Body:  "explore_body",
		},
		Learn: models.TabContent{
			Video: nil,
			Body:  "learn_body",
		},
	}
	var response string
	var jsonString = testutils.MakeJSONBody(badBookDetails, t)

	testutils.MakeHttpRequest("POST",
		ts.URL+"/books/nonexistant", jsonString, &response, t)

	require.Equal(t, "error", response)

}

// Test if deleting from book_contents works
func TestBookDetailDelete(t *testing.T) {
	var response interface{}
	testutils.MakeHttpRequest("DELETE",
		ts.URL+"/books/c_id/en", nil, &response, t)

	require.Equal(t, nil, response)
}

// Testif deleting book works
func TestBookDelete(t *testing.T) {
	var response interface{}
	testutils.MakeHttpRequest("DELETE", ts.URL+"/books/d_id", nil, &response, t)

	require.Equal(t, nil, response)
}

// Testing if delete on an invalid id returns an error
func TestBookDeleteOnInvalidId(t *testing.T) {
	var response string
	testutils.MakeHttpRequest("DELETE", ts.URL+"/books/nonexistant", nil, &response, t)
	require.Equal(t, "error", response)
}

// Test if book detail delete on an invalid id returns an error
func TestBookDetailDeleteOnInvalidId(t *testing.T) {
	var response string
	testutils.MakeHttpRequest("DELETE", ts.URL+"/books/nonexistant/en", nil, &response, t)
	require.Equal(t, "error", response)
}

// Test book detail on invalid language returns an error
func TestBookDetailDeleteOnInvalidLanguage(t *testing.T) {
	var response string
	testutils.MakeHttpRequest("DELETE", ts.URL+"/books/catcher/ge", nil, &response, t)
	require.Equal(t, "error", response)
}

// Test if updating entry in books work
func TestUpdateBook(t *testing.T) {
	var str = "updated_title"
	var updatedBook = models.APIUpdateBook{
		Title:  &str,
		Author: nil,
		Image:  nil,
	}
	jsonStr := testutils.MakeJSONBody(updatedBook, t)
	var response models.Book
	testutils.MakeHttpRequest("PATCH", ts.URL+"/books/update_me", jsonStr, &response, t)
	require.Equal(t, "updated_title", response.Title)
	require.Equal(t, "update_author", response.Author)

}

// Test if updating entry in book_contents works
func TestUpdateBookDetails(t *testing.T) {
	var read_vid string = "new_read_video"
	var updatedBook = models.APIUpdateBookDetails{
		Language: nil,
		Read: models.APIUpdateTabContents{
			Video: &read_vid,
			Body:  nil,
		},
		Explore: models.APIUpdateTabContents{
			Video: nil,
			Body:  nil,
		},
		Learn: models.APIUpdateTabContents{
			Video: nil,
			Body:  nil,
		},
	}

	jsonStr := testutils.MakeJSONBody(updatedBook, t)
	var response models.BookDetails
	testutils.MakeHttpRequest("PATCH", ts.URL+"/books/update_me/en",
		jsonStr, &response, t)
	require.Equal(t, "new_read_video", *response.Read.Video)
}
