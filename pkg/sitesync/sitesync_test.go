package sitesync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyVersion(t *testing.T) {
	is_true := CompareCurrentVersion("")
	assert.Equal(t, is_true, false, `currentVersion must not be empty.`)
}

func TestSetVersion(t *testing.T) {
	is_true := CompareCurrentVersion("bla")
	assert.Equal(t, is_true, true, `Set version.`)
}

func TestUpdateVersion(t *testing.T) {
	InitConfig("name", "namespace")
	updateCurrentVersion("blub")
	is_false := CompareCurrentVersion("bla")
	assert.Equal(t, is_false, false, `Update version.`)
}

func TestEqualVersions(t *testing.T) {
	InitConfig("name", "namespace")
	updateCurrentVersion("bla")
	is_true := CompareCurrentVersion("bla")
	assert.Equal(t, is_true, true, `Version is equal.`)
}
