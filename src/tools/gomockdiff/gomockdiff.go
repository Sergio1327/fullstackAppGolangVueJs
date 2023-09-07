package gomockdiff

import (
	"fmt"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

type eqMatcher struct {
	want interface{}
}

func EqMatcher(want interface{}) eqMatcher {
	return eqMatcher{
		want: want,
	}
}

func (eq eqMatcher) Matches(got interface{}) bool {
	return gomock.Eq(eq.want).Matches(got)
}

func (eq eqMatcher) Got(got interface{}) string {
	return fmt.Sprintf("%v (%T)\nDiff (-got +want):\n%s", got, got, strings.TrimSpace(cmp.Diff(got, eq.want)))
}

func (eq eqMatcher) String() string {
	return fmt.Sprintf("%v (%T)\n", eq.want, eq.want)
}
