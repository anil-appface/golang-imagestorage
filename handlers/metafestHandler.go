package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anil/golang-imagestorage/store"
	util "github.com/anil/golang-imagestorage/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

//MetafestHandler all the albums & image handlers
type MetafestHandler struct {
	db *sql.DB
}

func NewMetafestHandler(_db *sql.DB) *MetafestHandler {
	return &MetafestHandler{
		db: _db,
	}
}

//InsertAlbumHandler insert handler
func (mh *MetafestHandler) InsertAlbumHandler(w http.ResponseWriter, r *http.Request) {

	var p store.MetafestModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	mrs := store.NewMetafestRepoSQL(mh.db)

	if err := mrs.Insert(&p); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, p)
}

//GetAllAlbumHandler get all the albums.
func (mh *MetafestHandler) GetAllAlbumHandler(w http.ResponseWriter, r *http.Request) {

	log.Tracef("getting the album from db: %#v")
	mrs := store.NewMetafestRepoSQL(mh.db)

	album, err := mrs.MatchParent("")
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			util.RespondWithError(w, http.StatusNotFound, "No albums found")
		default:
			util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.RespondWithJSON(w, http.StatusOK, album)
}

//GetAlbumHandler get all the albums.
func (mh *MetafestHandler) GetAlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumID := vars["albumId"]

	log.Tracef("getting the album from db: %#v", albumID)
	mrs := store.NewMetafestRepoSQL(mh.db)

	album, err := mrs.Get(albumID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			util.RespondWithError(w, http.StatusNotFound, "Album not found")
		default:
			util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.RespondWithJSON(w, http.StatusOK, album)
}

//DeleteAlbumHandler get all the albums.
func (mh *MetafestHandler) DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	albumID := vars["albumId"]

	log.Tracef("getting the album from db: %#v", albumID)
	mrs := store.NewMetafestRepoSQL(mh.db)

	if err := mrs.Delete(albumID); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusOK, nil)
}

//InsertImageHandler insert handler
func (mh *MetafestHandler) InsertImageHandler(w http.ResponseWriter, r *http.Request) {

}
