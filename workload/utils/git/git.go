package git

import (
	"github.com/yujuncen/brie-bench/workload/utils"
	"os"
)

type Repo struct {
	remote string
	local  string
}

func CloneHash(remote, to, hash string) (*Repo, error) {
	repo, err := Clone(remote, to)
	if err != nil {
		return nil, err
	}
	if hash != "" {
		err2 := repo.ResetHard(hash)
		if err2 != nil {
			return nil, err2
		}
	}
	return repo, nil
}

func Clone(remote, to string) (*Repo, error) {
	if to == "" {
		var err error
		to, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	if err := utils.NewCommand("git",
		"clone", remote, to).Opt(utils.SystemOutput).Run(); err != nil {
		return nil, err
	}
	return &Repo{remote: remote, local: to}, nil
}

func (r *Repo) ResetHard(hash string) error {
	return utils.NewCommand("git", "reset", "--hard", hash).
		Opt(utils.SystemOutput, utils.WorkDir(r.local)).
		Run()
}

func (r *Repo) Make(targets ...string) error {
	return utils.NewCommand("make", targets...).
		Opt(utils.SystemOutput, utils.WorkDir(r.local)).
		Run()
}
