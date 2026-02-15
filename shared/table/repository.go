package table

import (
	"errors"
	"path"

	"github.com/rs/zerolog/log"
)

type Repository struct {
	ItemTable ItemTable
}

func (r *Repository) Load(csvDirPath string) error {
	var errs error

	r.ItemTable = *NewItemTable()
	if err := r.ItemTable.Load(path.Join(csvDirPath, "item.csv")); err != nil {
		log.Err(err).Msg("failed to load item table")
		errs = errors.Join(errs, err)
	}

	return errs
}
