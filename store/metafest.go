package store

import (
	"database/sql"
	"time"
)

const (
	FolderMode = 421
	FileMode   = 420
)

// MetafestModel represents the files that are stored in a album.
type MetafestModel struct {
	ParentQID string `json:"parentQid,omitempty"` // Folder QID
	QID       string `json:"fileQid,omitempty"`   // File QID

	// Common metadata
	Mode    uint32    `json:"fileMode,omitempty"`
	Name    string    `json:"fileName"`
	Size    uint64    `json:"size,omitempty"`
	ModTime time.Time `json:"modTime,omitempty"`
}

// MetafestRepo represents a metafest repository.
// To keep all scaling options open the metafest is immutable,
// there is (no update method).
type MetafestRepo interface {
	Get(qid string) (*MetafestModel, error)
	MatchParent(parentQid string) ([]*MetafestModel, error)
	Insert(r *MetafestModel) error
	Delete(qid string)
	IsExists(qid string) (bool, error)
}

// MetafestRepoSQL is an SQL implementation of the MetafestRepo
type MetafestRepoSQL struct {
	db Prepper
}

//NewMetafestRepoSQL creates a new instance of MetafestRepoSQL
func NewMetafestRepoSQL(db Prepper) *MetafestRepoSQL {
	return &MetafestRepoSQL{
		db: db,
	}
}

//Get retrieve the file from the 'metafest' table by QID
func (mrs *MetafestRepoSQL) Get(qid string) (*MetafestModel, error) {
	log.Tracef("opting to retrieve metafest with QID >> %s", qid)
	metafestModel := &MetafestModel{
		QID: qid,
	}

	stmt, err := mrs.db.Prepare("SELECT parent_qid,mode, name, size, timestamp FROM metafest WHERE qid = ?;")
	if err != nil {
		log.Errorf("error while prepraring the statement")
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(qid)

	if row != nil {
		err = row.Scan(&metafestModel.ParentQID, &metafestModel.Mode, &metafestModel.Name, &metafestModel.Size, &metafestModel.ModTime)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debugf("qid >> %q no rows", qid)
			return nil, err
		}
		return nil, err
	}
	log.Tracef("qid of metafest %s retrieved", metafestModel.QID)
	return metafestModel, nil
}

//MatchParent : retrieve the file from the 'metafest' table by QID
func (mrs *MetafestRepoSQL) MatchParent(parentQid string) ([]*MetafestModel, error) {

	log.Tracef("getting records from folder id %s", parentQid)

	metafestModels := make([]*MetafestModel, 0)
	stmt, err := mrs.db.Prepare("SELECT qid, parent_qid,mode, name, size, timestamp FROM metafest WHERE parent_qid = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(parentQid)
	if err != nil {
		log.Errorf("could not execute the query")
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	for rows.Next() {
		metafestModel := &MetafestModel{}
		rows.Scan(&metafestModel.QID, &metafestModel.ParentQID, &metafestModel.Mode, &metafestModel.Name, &metafestModel.Size, &metafestModel.ModTime)
		metafestModels = append(metafestModels, metafestModel)
	}
	log.Tracef("retrieved %d record(s)", len(metafestModels))
	return metafestModels, nil
}

//Insert : inserts a new metafest to 'metafest' table
func (mrs *MetafestRepoSQL) Insert(metafest *MetafestModel) error {
	log.Tracef("opting to insert resource... >> %#v", metafest)
	stmt, err := mrs.db.Prepare("INSERT INTO metafest(qid, parent_qid ,mode, name, size, timestamp) VALUES(?,?,?,?,?,NOW());")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(metafest.QID, metafest.ParentQID, metafest.Mode, metafest.Name, metafest.Size)
	if err != nil {
		return err
	}

	log.Tracef("insert metafest... complete.")
	return nil
}

//Delete :: deletes records from manifest based on qid
func (mrs *MetafestRepoSQL) Delete(qid string) error {
	log.Tracef("abort deleting metafest id %q", qid)
	stmt, err := mrs.db.Prepare("DELETE FROM metafest WHERE qid = ?;")
	if err != nil {
		log.Errorf("abort deleting metafest id %#v", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(qid)
	if err != nil {
		log.Errorf("abort deleting metafest id %#v", err)
		return err
	}
	return nil
}
