package manager

import (
	"context"
	"time"

	"github.com/juju/errors"
	"github.com/pingcap/github-base/util"
)

const EXPIRE = time.Hour

type IsMember struct {
	expected time.Time
	is       bool
}

func (mgr *Manager) isMember(login string) (bool, error) {
	isMember, exist := mgr.Members[login]
	if exist && isMember.expected.After(time.Now()) {
		return isMember.is, nil
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
		mgr.Members[login] = IsMember{
			expected: time.Now().Add(EXPIRE),
			is:       isPingCAPMember,
		}
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
		mgr.Members[login] = IsMember{
			expected: time.Now().Add(EXPIRE),
			is:       isPingCAPMember,
		}
		return isTikvMember, nil
	} else if err != nil {
		return false, errors.Trace(err)
	}

	mgr.Members[login] = IsMember{
		expected: time.Now().Add(EXPIRE),
		is:       isPingCAPMember,
	}
	return false, nil
}

// IsMember return if a user is pingcap/tikv member
func (mgr *Manager) IsMember(login string) (bool, error) {
	return mgr.isMember(login)
}
