package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Saintrad/todo-server-client/internal/richerror"
)

type fileState struct {
	NextID int    `json:"next_id"`
	Users  []User `json:"users"`
}

type FileUserRepo struct {
	mu       sync.Mutex
	filePath string
	state    fileState
}

// NewFileUserRepo loads state from file if present, otherwise starts empty.
func NewFileUserRepo(path string) (*FileUserRepo, error) {
	r := &FileUserRepo{
		filePath: path,
		state: fileState{
			NextID: 1,
			Users:  make([]User, 0),
		},
	}

	// Ensure parent dir exists (e.g., data/)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Missing file is not an error: start empty
		if errors.Is(err, os.ErrNotExist) {
			return r, nil
		}
		return nil, err
	}

	// Empty file: treat as empty state
	if len(data) == 0 {
		return r, nil
	}

	var st fileState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Defensive defaults
	if st.NextID <= 0 {
		st.NextID = computeNextID(st.Users)
	}
	if st.Users == nil {
		st.Users = make([]User, 0)
	}

	r.state = st
	return r, nil
}

func computeNextID(Users []User) int {
	max := 0
	for _, t := range Users {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

// Create assigns an ID, stores the User, and persists to disk.
func (r *FileUserRepo) Create(User User) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	User.ID = r.state.NextID
	r.state.NextID++

	r.state.Users = append(r.state.Users, User)

	if err := r.saveLocked(); err != nil {
		return -1, err
	}
	return User.ID, nil
}

// saveLocked persists r.state to disk atomically.
// Call only while holding r.mu.
func (r *FileUserRepo) saveLocked() error {
	b, err := json.MarshalIndent(r.state, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	dir := filepath.Dir(r.filePath)
	tmp, err := os.CreateTemp(dir, "Users-*.tmp")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName) // no-op if rename succeeded
	}()

	if _, err := tmp.Write(b); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// Atomic replace on most OS/filesystems when same directory
	if err := os.Rename(tmpName, r.filePath); err != nil {
		return err
	}
	return nil
}

func (r *FileUserRepo) List() ([]User, *richerror.AppError) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]User, 0)

	for _, User := range r.state.Users {

		list = append(list, User)
	}

	return list, nil
}

func (r *FileUserRepo) GetByEmail(email string) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	Users := r.state.Users

	for _, User := range Users {
		if User.Email == email {
			return User, nil
		}
	}

	return User{}, richerror.NotFound("User not found", nil)
}

func (r *FileUserRepo) EmailExists(email string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	Users := r.state.Users

	for _, User := range Users {
		if User.Email == email {
			return true, nil
		}
	}

	return false, nil
}