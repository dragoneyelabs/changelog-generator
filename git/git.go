package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"encoding/hex"
	"sort"
	"time"
)

// Repository data
type Repository struct {
	r *git.Repository
}

// Tag data
type Tag struct {
	Hash Hash
	Name string
	Date time.Time
}

// Commit data
type Commit struct {
	Hash    Hash
	Date    time.Time
	Message string
}

// Hash data
type Hash [20]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// GetLocalRepository returns local repository data given the local path
func GetLocalRepository(path string) (*Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repository{r}, nil
}

// GetRemoteURLs returns the remote urls
func (r *Repository) GetRemoteURLs(remote string) ([]string, error) {
	re, err := r.r.Remote(remote)
	if err != nil {
		return nil, err
	}
	return re.Config().URLs, nil
}

// GetHead returns the HEAD commit
func (r *Repository) GetHead() (*Commit, error) {
	ref, err := r.r.Head()
	if err != nil {
		return nil, err
	}
	c, err := r.r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	return &Commit{
		Hash:    Hash(c.Hash),
		Message: c.Message,
		Date:    c.Committer.When,
	}, nil
}

// GetTagsSortedByDateAsc returns the repository tag list sorted by date ascending
func (r *Repository) GetTagsSortedByDateAsc() ([]Tag, error) {
	tags := make([]Tag, 0)

	t, err := r.r.Tags()
	if err != nil {
		return nil, err
	}
	defer t.Close()
	t.ForEach(func(ref *plumbing.Reference) error { //nolint:errcheck,gosec
		var date time.Time
		var hash Hash
		c, err := r.r.CommitObject(ref.Hash())
		if err != nil {
			tag, err := r.r.TagObject(ref.Hash())
			if err != nil {
				return nil
			}
			date = tag.Tagger.When
			hash, err = toHash(tag.Target.String())
			if err != nil {
				return nil
			}
		} else {
			date = c.Committer.When
			hash = Hash(ref.Hash())
		}
		tags = append(tags, Tag{
			Hash: hash,
			Name: ref.Name().Short(),
			Date: date,
		})
		return nil
	})

	// order by date asc
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.Before(tags[j].Date)
	})

	return tags, nil
}

// GetCommitsSortedByDateDesc returns the repository commits sorted by date ascending from the given hash and ignoring
// the specified hashes
func (r *Repository) GetCommitsSortedByDateDesc(from Hash, ignore ...Hash) ([]Commit, error) {
	commits := make([]Commit, 0)
	it, err := r.log(from, ignore...)
	if err != nil {
		return nil, err
	}
	defer it.Close()
	for c, err := it.Next(); err == nil; c, err = it.Next() {
		commits = append(commits, Commit{
			Hash:    Hash(c.Hash),
			Date:    c.Committer.When,
			Message: c.Message,
		})
	}
	return commits, nil
}

func (r *Repository) log(from Hash, ignore ...Hash) (object.CommitIter, error) {
	h := plumbing.Hash(from)
	commit, err := r.r.CommitObject(h)
	if err != nil {
		return nil, err
	}
	return object.NewCommitIterCTime(commit, nil, toPlumbingHashSlice(ignore)), nil
}

func toPlumbingHashSlice(s []Hash) []plumbing.Hash {
	out := make([]plumbing.Hash, 0)
	for _, h := range s {
		out = append(out, plumbing.Hash(h))
	}
	return out
}

func toHash(s string) (Hash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}
	var arr [20]byte
	copy(arr[:], b)
	return Hash(arr), err
}
