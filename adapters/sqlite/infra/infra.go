package infra

import (
	"github.com/jmoiron/sqlx"
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
	Image      im.SQLiteImageRepo
	Collection clc.SQLiteCollectionRepo
	Label      lbl.SQLiteLabelRepo
	ImageStore im_store.ImageStore
	FileStore  af_store.Interface
	Annotation an.SQLiteAnnotationRepo
	Scroller   scr.SQLiteScrollerRepo
	Group      grp.SQLiteGroupRepo
	User       usr.SQLiteUserRepo
	Db         *sqlx.DB
}

type SQLiteImageStoreRepo struct {
	im.SQLiteImageRepo
	clc.SQLiteCollectionRepo
	lbl.SQLiteLabelRepo
	an.SQLiteAnnotationRepo
}

func NewSQLiteInfra(db *sqlx.DB, fstore af_store.Interface) SQLiteInfra {
	imrepo := im.NewSQLiteImageRepo(db)
	anrepo := an.NewSQLiteAnnotationRepo(db)
	clrepo := clc.NewSQLiteCollectionRepo(db)
	lbrepo := lbl.NewSQLiteLabelRepo(db)
	grprepo := grp.NewSQLiteGroupRepo(db)
	usrrepo := usr.NewSQLiteUserRepo(db)
	imstorerepo := SQLiteImageStoreRepo{imrepo, clrepo, lbrepo, anrepo}
	imstore := im_store.New(imstorerepo, fstore)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	return SQLiteInfra{
		Image:      imrepo,
		Collection: clrepo,
		Label:      lbrepo,
		ImageStore: imstore,
		FileStore:  fstore,
		Annotation: anrepo,
		Scroller:   scrrepo,
		Group:      grprepo,
		User:       usrrepo,
		Db:         db,
	}

}
