// Package stdlib covers protection scope: stdlib and third-party packages
// never belong to the current module ("factory"), so their types are
// silent by default, with no fence or other configuration needed.
package stdlib

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

func NoDiagnostics() {
	_ = strings.Builder{}
	_ = &strings.Builder{}
	_ = sync.WaitGroup{}
	_ = &sync.WaitGroup{}
	_ = http.Header{}
	_ = time.Duration(5)
}
