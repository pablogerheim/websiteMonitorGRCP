package repositories

import (
	"websiteMonitor/repositories/spanner"

	"gorm.io/gorm"
)

type Repositories struct {
	Spanner spanner.ISpannerSite
}

func (r *Repositories) Inject(DB *gorm.DB) {
	r.Spanner = spanner.NewSpannerRepository(DB)
}

