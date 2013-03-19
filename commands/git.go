package commands

import (
	"github.com/jgrocho/go-git2"
)

func initRepo(prefix string) error {
	_, err := git2.InitRepository(prefix, false)
	return err
}

func addAndCommit(prefix, file, message string) error {
	repo, err := git2.Open(prefix)
	if err != nil {
		return err
	}
	defer repo.Free()

	config, err := git2.OpenGlobalConfig()
	if err != nil {
		return err
	}
	defer config.Free()

	index, err := repo.Index()
	if err != nil {
		return err
	}
	defer index.Free()

	if err := index.Add(file, 0); err != nil {
		return err
	}
	if err := index.Write(); err != nil {
		return err
	}

	treeOid, err := index.CreateTree()
	if err != nil {
		return err
	}
	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		return err
	}
	defer tree.Free()

	userName, err := config.GetString("user.name")
	if err != nil {
		return err
	}
	userEmail, err := config.GetString("user.email")
	if err != nil {
		return err
	}
	author, err := git2.SignatureNow(userName, userEmail)
	if err != nil {
		return err
	}
	defer author.Free()

	committer, err := git2.SignatureNow("pass", "pass@github.com")
	if err != nil {
		return err
	}
	defer committer.Free()

	if repo.Empty() {
		_, err = repo.CreateCommit("HEAD", author, committer, "UTF-8", message, tree)
		if err != nil {
			return err
		}
	} else {
		head, err := repo.Head()
		if err != nil {
			return err
		}
		defer head.Free()

		parent, err := repo.LookupCommit(head.Oid())
		if err != nil {
			return err
		}
		defer parent.Free()

		_, err = repo.CreateCommit("HEAD", author, committer, "UTF-8", message, tree, parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAndCommit(prefix, file, message string) error {
	repo, err := git2.Open(prefix)
	if err != nil {
		return err
	}
	defer repo.Free()

	config, err := git2.OpenGlobalConfig()
	if err != nil {
		return err
	}
	defer config.Free()

	index, err := repo.Index()
	if err != nil {
		return err
	}
	defer index.Free()

	if err := index.Add(file, 0); err != nil {
		return err
	}
	if err := index.Remove(0); err != nil {
		return err
	}
	if err := index.Write(); err != nil {
		return err
	}

	treeOid, err := index.CreateTree()
	if err != nil {
		return err
	}
	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		return err
	}
	defer tree.Free()

	userName, err := config.GetString("user.name")
	if err != nil {
		return err
	}
	userEmail, err := config.GetString("user.email")
	if err != nil {
		return err
	}
	author, err := git2.SignatureNow(userName, userEmail)
	if err != nil {
		return err
	}
	defer author.Free()

	committer, err := git2.SignatureNow("pass", "pass@github.com")
	if err != nil {
		return err
	}
	defer committer.Free()

	if repo.Empty() {
		_, err = repo.CreateCommit("HEAD", author, committer, "UTF-8", message, tree)
		if err != nil {
			return err
		}
	} else {
		head, err := repo.Head()
		if err != nil {
			return err
		}
		defer head.Free()

		parent, err := repo.LookupCommit(head.Oid())
		if err != nil {
			return err
		}
		defer parent.Free()

		_, err = repo.CreateCommit("HEAD", author, committer, "UTF-8", message, tree, parent)
		if err != nil {
			return err
		}
	}

	return nil
}
