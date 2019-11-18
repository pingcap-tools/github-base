package types

// Repo struct
type Repo struct {
	Owner string
	Repo  string
}

// GetOwner get repo's owner
func (r *Repo)GetOwner() string {
	return r.Owner
}

// GetRepo get repo's repo name
func (r *Repo)GetRepo() string {
	return r.Repo
}

