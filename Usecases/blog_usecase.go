// this package will be the interface betweeen the controller/s and the repository
package usecases

import (
	repositories "g6/blog-api/Repositories"
	"time"
)

type blogUsecase struct {
	blogRepo repositories.BlogRepository
	ctxtimeout  time.Duration
}
