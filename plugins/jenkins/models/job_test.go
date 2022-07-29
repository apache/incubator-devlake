package models

import (
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFolderInput(t *testing.T) {
	it := helper.NewQueueIterator()
	it.Push(NewFolderInput("a"))
	it.Push(NewFolderInput("b"))
	require.True(t, it.HasNext())
	folderRaw, err := it.Fetch()
	require.NoError(t, err)
	folder := folderRaw.(*FolderInput)
	require.Equal(t, "a", folder.Data())
	require.True(t, it.HasNext())
	folderRaw, err = it.Fetch()
	require.NoError(t, err)
	folder = folderRaw.(*FolderInput)
	require.Equal(t, "b", folder.Data())
	require.False(t, it.HasNext())
}
