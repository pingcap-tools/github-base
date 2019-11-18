package manager

import (
	"context"
	"github.com/juju/errors"
	"github.com/pingcap/github-base/util"
)

func (mgr *Manager)isMember(login string) (bool, error) {
	isMember, exist := mgr.Members[login]
	if exist {
		return isMember, nil
	}

	isPingCAPMember := false
	err := util.RetryOnError(context.Background(), 3, func() error {
		r, _, err := mgr.Github.Organizations.IsMember(context.Background(), "pingcap", login)
		if err == nil {
			isPingCAPMember = r
		}
		return errors.Trace(err)
	})
	if err == nil && isPingCAPMember {
		mgr.Members[login] = isPingCAPMember
		return isPingCAPMember, nil
	} else if err != nil {
		return false, errors.Trace(err)
	}

	isTikvMember := false
	err = util.RetryOnError(context.Background(), 3, func() error {
		r, _, err := mgr.Github.Organizations.IsMember(context.Background(), "tikv", login)
		if err == nil {
			isTikvMember = r
		}
		return errors.Trace(err)
	})
	if err == nil && isTikvMember {
		mgr.Members[login] = isTikvMember
		return isTikvMember, nil
	} else if err != nil {
		return false, errors.Trace(err)
	}

	mgr.Members[login] = false
	return false, nil
}

// IsMember return if a user is pingcap/tikv member
func (mgr *Manager)IsMember(login string) (bool, error) {
	return mgr.isMember(login)
}
