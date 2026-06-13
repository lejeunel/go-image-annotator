package sqlite

import (
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	grp "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	scr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/scroll"
	usr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	af_store "github.com/lejeunel/go-image-annotator/modules/file-store"
	im_store "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type SQLiteRepos struct {
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

func NewSQLiteRepos(db *sqlx.DB, fstore af_store.Interface) SQLiteRepos {
	imrepo := im.NewSQLiteImageRepo(db)
	anrepo := an.NewSQLiteAnnotationRepo(db)
	clrepo := clc.NewSQLiteCollectionRepo(db)
	lbrepo := lbl.NewSQLiteLabelRepo(db)
	grprepo := grp.NewSQLiteGroupRepo(db)
	usrrepo := usr.NewSQLiteUserRepo(db)
	imstorerepo := SQLiteImageStoreRepo{imrepo, clrepo, lbrepo, anrepo}
	imstore := im_store.New(imstorerepo, fstore)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	return SQLiteRepos{
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
