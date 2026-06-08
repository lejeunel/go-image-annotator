package infra

import (
	"github.com/jmoiron/sqlx"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	an "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	im "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/label"
	scr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/scroll"
	af_store "github.com/lejeunel/go-image-annotator/app/file-store"
	im_store "github.com/lejeunel/go-image-annotator/app/image-store"
)

type SQLiteInfra struct {
	ImageRepo      *im.SQLiteImageRepo
	CollectionRepo *clc.SQLiteCollectionRepo
	LabelRepo      *lbl.SQLiteLabelRepo
	ImageStore     *im_store.ImageStore
	FileStore      *af_store.FileStore
	AnnotationRepo *an.SQLiteAnnotationRepo
	ScrollerRepo   *scr.SQLiteScrollerRepo
	Db             *sqlx.DB
}

type SQLiteImageStoreRepo struct {
	*im.SQLiteImageRepo
	*clc.SQLiteCollectionRepo
	*lbl.SQLiteLabelRepo
	*an.SQLiteAnnotationRepo
}

func NewSQLiteInfra(dbPath, artefactDir string) *SQLiteInfra {
	db := db.NewSQLiteDB(dbPath)
	imrepo := im.NewSQLiteImageRepo(db)
	anrepo := an.NewSQLiteAnnotationRepo(db)
	clrepo := clc.NewSQLiteCollectionRepo(db)
	lbrepo := lbl.NewSQLiteLabelRepo(db)
	afrepo := af_store.NewFileStore(artefactDir)
	imstorerepo := SQLiteImageStoreRepo{imrepo, clrepo, lbrepo, anrepo}
	imstore := im_store.New(imstorerepo, afrepo)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	return &SQLiteInfra{
		ImageRepo:      imrepo,
		CollectionRepo: clrepo,
		LabelRepo:      lbrepo,
		ImageStore:     imstore,
		FileStore:      afrepo,
		AnnotationRepo: anrepo,
		ScrollerRepo:   scrrepo,
		Db:             db,
	}

}
