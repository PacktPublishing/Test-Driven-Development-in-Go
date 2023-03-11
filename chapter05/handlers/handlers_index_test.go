package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter05/handlers"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handlers integration", func() {
	var svr *httptest.Server
	var eb db.Book

	BeforeEach(func() {
		eb = db.Book{
			ID:     uuid.New().String(),
			Name:   "My first integration test",
			Status: db.Available.String(),
		}
		bs := db.NewBookService([]db.Book{eb}, nil)
		ha := handlers.NewHandler(bs, nil)
		svr = httptest.NewServer(http.HandlerFunc(ha.Index))
	})

	AfterEach(func() {
		svr.Close()
	})

	Describe("Index endpoint", func() {
		Context("with one existing book", func() {
			It("should return book", func() {
				r, err := http.Get(svr.URL)
				Expect(err).To(BeNil())
				Expect(r.StatusCode).To(Equal(http.StatusOK))

				body, err := io.ReadAll(r.Body)
				r.Body.Close()
				Expect(err).To(BeNil())

				var resp handlers.Response
				err = json.Unmarshal(body, &resp)

				Expect(err).To(BeNil())
				Expect(len(resp.Books)).To(Equal(1))
				Expect(resp.Books).To(ContainElement(eb))
			})
		})
	})
})
