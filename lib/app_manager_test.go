package steamreviewfetcher

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppManager_loadOrSaveGameDetails(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testappmanager")
	require.Nil(t, err)
	defer os.RemoveAll(tmpdir)

	var manager = AppManager{Directory: tmpdir}
	manager.game.Name = "Test App"
	manager.game.Platforms.Windows = true

	err = manager.saveGameDetails()
	require.Nil(t, err)

	manager.game = Game{}
	require.Equal(t, manager.game.Name, "")
	require.False(t, manager.game.Platforms.Windows)

	err = manager.loadGameDetails()
	require.Nil(t, err)
	require.Equal(t, manager.game.Name, "Test App")
	require.True(t, manager.game.Platforms.Windows)
}
