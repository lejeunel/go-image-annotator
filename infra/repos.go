package infra

import (
	af_store "github.com/lejeunel/go-image-annotator-v2/app/file-store"
	im_store "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	db "github.com/lejeunel/go-image-annotator-v2/infra/db"
	an "github.com/lejeunel/go-image-annotator-v2/infra/db/annotation"
	clc "github.com/lejeunel/go-image-annotator-v2/infra/db/collection"
	im "github.com/lejeunel/go-image-annotator-v2/infra/db/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/infra/db/label"
	scr "github.com/lejeunel/go-image-annotator-v2/infra/db/scroll"
)

type SQLiteInfra struct {
	ImageRepo      *im.SQLiteImageRepo
	CollectionRepo *clc.SQLiteCollectionRepo
	LabelRepo      *lbl.SQLiteLabelRepo
	ImageStore     *im_store.ImageStore
	FileStore      *af_store.FileStore
	AnnotationRepo *an.SQLiteAnnotationRepo
	ScrollerRepo   *scr.SQLiteScrollerRepo
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
	}

}
