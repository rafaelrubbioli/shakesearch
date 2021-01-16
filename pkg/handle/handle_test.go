package handle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"pulley.com/shakesearch/mocks"
)

func TestSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	searchMock := mocks.NewMockSearcher(ctrl)

	server := httptest.NewServer(http.HandlerFunc(Search(searchMock)))
	defer server.Close()

	expect := httpexpect.New(t, server.URL)
	t.Run("success", func(t *testing.T) {
		expected := []string{"This is the expected result for the search hamlet."}
		searchMock.EXPECT().
			Search("hamlet").
			Return(expected)

		e := expect.GET("/search").
			WithQuery("q", "hamlet").
			Expect()

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(expected)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, e.Raw().StatusCode)
		require.Equal(t, buf.String(), e.Body().Raw())
	})

	t.Run("multiple query search", func(t *testing.T) {
		expected := []string{"This is the expected result for the search hamlet.", "This is the expected result for the search second."}
		searchMock.EXPECT().
			Search("hamlet").
			Return([]string{expected[0]})

		searchMock.EXPECT().
			Search("second").
			Return([]string{expected[1]})

		e := expect.GET("/search").
			WithQuery("q", "hamlet").
			WithQuery("q", "second").
			Expect()

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(expected)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, e.Raw().StatusCode)
		require.Equal(t, buf.String(), e.Body().Raw())
	})

	t.Run("query not given", func(t *testing.T) {
		e := expect.GET("/search").
			Expect()

		require.Equal(t, http.StatusBadRequest, e.Raw().StatusCode)
		require.Equal(t, "missing search query in URL params", e.Body().Raw())
	})
}
