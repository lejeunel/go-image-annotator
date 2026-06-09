package infra

import (
	"github.com/jmoiron/sqlx"
	db "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	an "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	grp "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/group"
	im "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/label"
	scr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/scroll"
	usr "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/user"
	af_store "github.com/lejeunel/go-image-annotator/app/file-store"
	im_store "github.com/lejeunel/go-image-annotator/app/image-store"
)

type SQLiteInfra struct {
	Image      *im.SQLiteImageRepo
	Collection *clc.SQLiteCollectionRepo
	Label      *lbl.SQLiteLabelRepo
	ImageStore *im_store.ImageStore
	FileStore  *af_store.FileStore
	Annotation *an.SQLiteAnnotationRepo
	Scroller   *scr.SQLiteScrollerRepo
	Group      *grp.SQLiteGroupRepo
	User       *usr.SQLiteUserRepo
	Db         *sqlx.DB
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
	grprepo := grp.NewSQLiteGroupRepo(db)
	usrrepo := usr.NewSQLiteUserRepo(db)
	imstorerepo := SQLiteImageStoreRepo{imrepo, clrepo, lbrepo, anrepo}
	imstore := im_store.New(imstorerepo, afrepo)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	return &SQLiteInfra{
		Image:      imrepo,
		Collection: clrepo,
		Label:      lbrepo,
		ImageStore: imstore,
		FileStore:  afrepo,
		Annotation: anrepo,
		Scroller:   scrrepo,
		Group:      grprepo,
		User:       usrrepo,
		Db:         db,
	}

}
