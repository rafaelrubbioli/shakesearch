package searcher

import (
	"index/suffixarray"
	"strings"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		searcher, err := New("../../completeworks.txt")
		require.NoError(t, err)
		require.NotNil(t, searcher)
	})

	t.Run("error reading input file", func(t *testing.T) {
		searcher, err := New("notFound.txt")
		require.Error(t, err)
		require.Nil(t, searcher)
	})
}

func TestSearcher_Search(t *testing.T) {
	text := `PAROLLES. Well, thou hast a son shall take this disgrace off me; scurvy, old, filthy, scurvy lord! Well, I must be patient; there is no fettering of ` +
		`authority. I’ll beat him, by my life, if I can meet him with any convenience, an he were double and double a lord. I’ll have no more pity of his age than I ` +
		`would have of—I’ll beat him, and if I could but meet him again. Enter Lafew. LAFEW. Sirrah, your lord and master’s married; there’s news for you; you have a ` +
		`new mistress.`

	searcher := searcher{
		CompleteWorks: text,
		SuffixArray:   suffixarray.New([]byte(strings.ToLower(text))),
		cache:         cache.New(time.Hour, time.Hour),
	}

	t.Run("searching for first or last characters", func(t *testing.T) {
		results := searcher.Search("mistress")
		require.Len(t, results, 1)
		require.Equal(t, `" Sirrah, your lord and master’s married; there’s news for you; you have a new mistress."`, results[0])

		results = searcher.Search("PAROLLES")
		require.Len(t, results, 1)
		require.Equal(t, `"PAROLLES."`, results[0])
	})

	t.Run("search case-insensitive", func(t *testing.T) {
		results := searcher.Search("parolles")
		require.Len(t, results, 1)
		require.Equal(t, `"PAROLLES."`, results[0])
	})

	t.Run("result is successfully cached and returned", func(t *testing.T) {
		before := searcher.Search("parolles")
		require.Len(t, before, 1)
		require.Equal(t, `"PAROLLES."`, before[0])

		cached, ok := searcher.cache.Get("parolles")
		require.True(t, ok)
		require.Equal(t, cached, before)

		after := searcher.Search("parolles")
		require.Len(t, after, 1)
		require.Equal(t, before, after)
	})
}
