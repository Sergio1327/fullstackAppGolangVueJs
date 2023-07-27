package compare_test

import (
	"fmt"
	"projects_template/tools/compare"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRebuild(t *testing.T) {
	r := require.New(t)

	arrInt := []int{1, 2, 3}
	arrString := []string{"1", "2", "3"}

	data := compare.Rebuild(arrInt, func(t int) string {
		return fmt.Sprintf("%d", t)
	})
	r.Equal(arrString, data)
}
