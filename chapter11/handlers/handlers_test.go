package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter11/handlers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestIndexIntegration in short mode.")
	}
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	// Arrange
	bs := db.NewBookService(testDB, nil)
	book := bs.Upsert(db.Book{
		Name:   "My first integration test",
		Status: db.Available.String(),
	})
	ha := handlers.NewHandler(bs, nil, nil)
	svr := httptest.NewServer(http.HandlerFunc(ha.Index))
	defer svr.Close()

	// Act
	r, err := http.Get(svr.URL)

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response[db.Book]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Contains(t, resp.Items, book)
}

func TestListBooksIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListBooksIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	eb := bs.Upsert(db.Book{
		Name:   "My first integration test",
		Status: db.Available.String(),
	})
	ha := handlers.NewHandler(bs, nil, nil)
	svr := httptest.NewServer(http.HandlerFunc(ha.ListBooks))
	defer svr.Close()

	// Act
	r, err := http.Get(svr.URL)

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response[db.Book]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Contains(t, resp.Items, eb)
}

func TestListMagazinesIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListMagazinesIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	ms := db.NewMagazineService(testDB, nil)
	em := ms.Upsert(db.Magazine{
		Name:   "My integration test",
		Status: db.Available.String(),
	})
	ha := handlers.NewHandler(nil, nil, ms)
	svr := httptest.NewServer(http.HandlerFunc(ha.ListMagazines))
	defer svr.Close()

	// Act
	r, err := http.Get(svr.URL)

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response[db.Magazine]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Contains(t, resp.Items, em)
}

func TestUserUpsertIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestIndexIntegration in short mode.")
	}
	// Arrange
	newUser := db.User{
		Name: "New user",
	}
	userPayload, err := json.Marshal(newUser)
	require.Nil(t, err)
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	us := db.NewUserService(testDB, nil, nil)
	ha := handlers.NewHandler(nil, us, nil)
	svr := httptest.NewServer(http.HandlerFunc(ha.UserUpsert))
	defer svr.Close()

	// Act
	r, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(userPayload))

	// Assert
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response[db.Book]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Equal(t, newUser.Name, resp.User.Name)
}

func TestBookUpsertIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestBookUpsertIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	ms := db.NewMagazineService(testDB, nil)
	us := db.NewUserService(testDB, bs, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	newBook := db.Book{
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	}
	bookPayload, err := json.Marshal(newBook)
	require.Nil(t, err)

	ha := handlers.NewHandler(bs, us, nil)
	svr := httptest.NewServer(http.HandlerFunc(ha.BookUpsert))
	defer svr.Close()

	// Act
	r, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(bookPayload))

	// Assert
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)
	var resp handlers.Response[db.Book]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, newBook.Name, resp.Items[0].Name)
	assert.Equal(t, db.Available.String(), resp.Items[0].Status)
}

func TestMagazineUpsertIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestMagazineUpsertIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	ms := db.NewMagazineService(testDB, nil)
	us := db.NewUserService(testDB, nil, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	newMag := db.Magazine{
		Name:    "Existing mag",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	}
	magPayload, err := json.Marshal(newMag)
	require.Nil(t, err)

	ha := handlers.NewHandler(nil, us, ms)
	svr := httptest.NewServer(http.HandlerFunc(ha.MagazineUpsert))
	defer svr.Close()

	// Act
	r, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(magPayload))

	// Assert
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)
	var resp handlers.Response[db.Magazine]
	err = json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, newMag.Name, resp.Items[0].Name)
	assert.Equal(t, db.Available.String(), resp.Items[0].Status)
}

func TestListUserByID_Books_Integration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListUserByID_Books_Integration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	ms := db.NewMagazineService(testDB, nil)
	us := db.NewUserService(testDB, bs, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	eb := bs.Upsert(db.Book{
		ID:      uuid.New().String(),
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	})
	ha := handlers.NewHandler(bs, us, nil)

	// Act
	path := fmt.Sprintf("/users/%s/books", eu.ID)
	req, err := http.NewRequest("GET", path, nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}/books", ha.ListUserByID_Books)
	router.ServeHTTP(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.Response[db.Book]
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.Nil(t, err)
	assert.Equal(t, eu.Name, resp.User.Name)
	assert.Equal(t, eu.ID, resp.User.ID)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, eb.Name, resp.Items[0].Name)
	assert.Equal(t, eb.ID, resp.Items[0].ID)
}
func TestListUserByID_Magazines_Integration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListUserByID_Magazines_Integration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	ms := db.NewMagazineService(testDB, nil)
	us := db.NewUserService(testDB, bs, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	em := ms.Upsert(db.Magazine{
		Name:    "Existing mag",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	})
	ha := handlers.NewHandler(bs, us, nil)

	// Act
	path := fmt.Sprintf("/users/%s/magazines", eu.ID)
	req, err := http.NewRequest("GET", path, nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}/magazines", ha.ListUserByID_Magazines)
	router.ServeHTTP(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.Response[db.Magazine]
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.Nil(t, err)
	assert.Equal(t, eu.Name, resp.User.Name)
	assert.Equal(t, eu.ID, resp.User.ID)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, em.Name, resp.Items[0].Name)
	assert.Equal(t, em.ID, resp.Items[0].ID)
}

func TestSwapBookIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestSwapBookIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	ps := db.NewPostingService()
	bs := db.NewBookService(testDB, ps)
	ms := db.NewMagazineService(testDB, ps)
	us := db.NewUserService(testDB, bs, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	swapUser, err := us.Upsert(db.User{
		Name: "Swap user",
	})
	require.Nil(t, err)
	eb := bs.Upsert(db.Book{
		Name:    "Existing book",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	})
	ha := handlers.NewHandler(bs, us, ms)

	// Act
	path := fmt.Sprintf("/books/%s?user=%s", eb.ID, swapUser.ID)
	req, err := http.NewRequest("POST", path, nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Methods("POST").Path("/books/{id}").Handler(http.HandlerFunc(ha.SwapBook))
	router.ServeHTTP(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.Response[db.Book]
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.Nil(t, err)
	assert.Equal(t, swapUser.Name, resp.User.Name)
	assert.Equal(t, swapUser.ID, resp.User.ID)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, eb.Name, resp.Items[0].Name)
	assert.Equal(t, eb.ID, resp.Items[0].ID)
	assert.Equal(t, db.Swapped.String(), resp.Items[0].Status)
}

func TestSwapMagazineIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestSwapMagazineIntegration in short mode.")
	}
	// Arrange
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	ps := db.NewPostingService()
	bs := db.NewBookService(testDB, ps)
	ms := db.NewMagazineService(testDB, ps)
	us := db.NewUserService(testDB, bs, ms)
	eu, err := us.Upsert(db.User{
		Name: "Existing user",
	})
	require.Nil(t, err)
	swapUser, err := us.Upsert(db.User{
		Name: "Swap user",
	})
	require.Nil(t, err)
	em := ms.Upsert(db.Magazine{
		Name:    "Existing mag",
		Status:  db.Available.String(),
		OwnerID: eu.ID,
	})
	ha := handlers.NewHandler(bs, us, ms)

	// Act
	path := fmt.Sprintf("/magazines/%s?user=%s", em.ID, swapUser.ID)
	req, err := http.NewRequest("POST", path, nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Methods("POST").Path("/magazines/{id}").Handler(http.HandlerFunc(ha.SwapMagazine))
	router.ServeHTTP(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.Response[db.Magazine]
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.Nil(t, err)
	assert.Equal(t, swapUser.Name, resp.User.Name)
	assert.Equal(t, swapUser.ID, resp.User.ID)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, em.Name, resp.Items[0].Name)
	assert.Equal(t, em.ID, resp.Items[0].ID)
	assert.Equal(t, db.Swapped.String(), resp.Items[0].Status)
}
