package storage

import (
	"github.com/pingcap/github-base/config"
	"github.com/pingcap/github-base/pkg/storage/basic"
	"github.com/pingcap/github-base/pkg/storage/mysql"

	"github.com/juju/errors"
)

func New(config *config.Database) (basic.Driver, error) {
	driver, err := mysql.Connect(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return driver, nil
}
