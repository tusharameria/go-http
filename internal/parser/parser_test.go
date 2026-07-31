package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserFeed_MessageSplitAcrossFeeds(t *testing.T) {
	p := &Parser{}

	msgs := p.Feed([]byte("Hel"))
	require.Empty(t, msgs)

	msgs = p.Feed([]byte("lo\n"))
	require.Len(t, msgs, 1)
	require.Equal(t, []byte("Hello"), msgs[0])
	require.Empty(t, p.aggregator)

	msgs = p.Feed([]byte("Hel\nl\no\n"))
	require.Len(t, msgs, 3)
	require.Equal(t, []byte("Hel"), msgs[0])
	require.Equal(t, []byte("l"), msgs[1])
	require.Equal(t, []byte("o"), msgs[2])
	require.Empty(t, p.aggregator)
}

func TestParserFeed_MultipleMessageSplitAcrossSingleFeed(t *testing.T) {
	p := &Parser{}

	msgs := p.Feed([]byte("Hel\nl\no\n"))
	require.NotEmpty(t, msgs)

	require.Len(t, msgs, 3)
	require.Equal(t, []byte("Hel"), msgs[0])
	require.Equal(t, []byte("l"), msgs[1])
	require.Equal(t, []byte("o"), msgs[2])
	require.Empty(t, p.aggregator)
}

func TestParserFeed_MultipleMessageSplitAcrossMultipleFeed(t *testing.T) {
	p := &Parser{}

	msgs := p.Feed([]byte("Hel\nl\no\n"))
	require.NotEmpty(t, msgs)
	require.Len(t, msgs, 3)
	require.Equal(t, []byte("Hel"), msgs[0])
	require.Equal(t, []byte("l"), msgs[1])
	require.Equal(t, []byte("o"), msgs[2])

	msgs = p.Feed([]byte("Where\nare you\nNOW\n?\n"))
	require.Len(t, msgs, 4)
	require.Equal(t, []byte("Where"), msgs[0])
	require.Equal(t, []byte("are you"), msgs[1])
	require.Equal(t, []byte("NOW"), msgs[2])
	require.Equal(t, []byte("?"), msgs[3])
	require.Empty(t, p.aggregator)
}
