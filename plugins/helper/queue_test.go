package helper

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueue(t *testing.T) {
	it := NewQueueIterator()
	it.Push(NewQueueIteratorNode("a"))
	it.Push(NewQueueIteratorNode("b"))
	require.True(t, it.HasNext())
	folderRaw, err := it.Fetch()
	require.NoError(t, err)
	node := folderRaw.(*QueueIteratorNode)
	require.Equal(t, "a", node.Data())
	require.True(t, it.HasNext())
	folderRaw, err = it.Fetch()
	require.NoError(t, err)
	node = folderRaw.(*QueueIteratorNode)
	require.Equal(t, "b", node.Data())
	require.False(t, it.HasNext())
}
