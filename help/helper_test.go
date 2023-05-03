package help

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayContains(t *testing.T) {
	array := []string{"test", "test2"}
	assert.True(t, ArrayContainsString(array, "test"))
	assert.False(t, ArrayContainsString(array, "test3"))
}

func TestGetPositionOfString(t *testing.T) {
	array := []string{"test", "test2"}
	assert.Equal(t, 0, GetPositionOfString(array, "test"))
	assert.Equal(t, -1, GetPositionOfString(array, "test3"))
}

func TestRemoveFromArray(t *testing.T) {
	array := []string{"test", "test2"}
	assert.Equal(t, array, RemoveFromArray(array, -1))
	assert.Equal(t, []string{"test2"}, RemoveFromArray(array, 0))
}

func TestRemoveStringFromArray(t *testing.T) {
	array := []string{"test", "test2", "test2", "test", "test"}
	assert.Equal(t, array, RemoveFirstStringOccurrenceFromArray(array, "test3"))
	assert.Equal(t, array[1:], RemoveFirstStringOccurrenceFromArray(array, "test"))
}

func TestTrimSuffix(t *testing.T) {
	str := "aasdasd"
	assert.Equal(t, "aasd", TrimSuffix(str, "asd"))
	assert.Equal(t, "aasd", TrimSuffix(str, "bvs"))
	assert.Equal(t, "aas", TrimSuffix(str, "aasd"))
}

func TestTrimPrefix(t *testing.T) {
	str := "aasdasd"
	assert.Equal(t, "dasd", TrimPrefix(str, "asd"))
	assert.Equal(t, "dasd", TrimPrefix(str, "bvs"))
	assert.Equal(t, "asd", TrimPrefix(str, "aasd"))
}

func TestDuplicationCheck(t *testing.T) {
	array := []string{"duplicat", "duplicat", "test"}
	array = RemoveDuplicates(array)
	assert.Equal(t, []string{"duplicat", "test"}, array)
}

func TestClearNamesArray(t *testing.T) {
	array := []string{"", "duplicat", "duplicat", "test", "", "duplicat"}
	ClearStringArray(&array)
	assert.Equal(t, []string{"duplicat", "test"}, array)
}

func TestTransformFromInputToNames(t *testing.T) {
	array := []string{"", "duplicat", "duplicat", "test", "", "duplicat"}
	array = DeleteMultiplesAndEmpty(array)
	assert.Equal(t, []string{"duplicat", "test"}, array)
}
