package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	grp "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	r "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	scr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/scroll"
	usr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	af_store "github.com/lejeunel/go-image-annotator/modules/file-store"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	im_store "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type LocalFSInfra struct {
	Image           im.SQLiteImageRepo
	Collection      clc.SQLiteCollectionRepo
	Label           lbl.SQLiteLabelRepo
	ImageStore      im_store.ImageStore
	ImageFileStore  af_store.Interface
	PolicyFileStore af_store.Interface
	Annotation      an.SQLiteAnnotationRepo
	Scroller        scr.SQLiteScrollerRepo
	Group           grp.SQLiteGroupRepo
	Role            r.SQLiteRoleRepo
	User            usr.SQLiteUserRepo
	Db              *sqlx.DB
}

func NewLocalFSInfra(db *sqlx.DB, basePath string) LocalFSInfra {
	imrepo := im.NewSQLiteImageRepo(db)
	anrepo := an.NewSQLiteAnnotationRepo(db)
	clrepo := clc.NewSQLiteCollectionRepo(db)
	lbrepo := lbl.NewSQLiteLabelRepo(db)
	grprepo := grp.NewSQLiteGroupRepo(db)
	rlrepo := r.NewSQLiteRoleRepo(db)
	usrrepo := usr.NewSQLiteUserRepo(db)
	imageFileStore := fs.NewFileStore(fmt.Sprintf("%v/%v", basePath, "images"))
	policyFileStore := fs.NewFileStore(fmt.Sprintf("%v/%v", basePath, "assets"))
	imstore := im_store.New(imrepo, clrepo, anrepo, imageFileStore)
	scrrepo := scr.NewSQLiteScrollerRepo(db)
	return LocalFSInfra{
		Image:           imrepo,
		Collection:      clrepo,
		Label:           lbrepo,
		ImageStore:      imstore,
		ImageFileStore:  imageFileStore,
		PolicyFileStore: policyFileStore,
		Annotation:      anrepo,
		Scroller:        scrrepo,
		Group:           grprepo,
		Role:            rlrepo,
		User:            usrrepo,
		Db:              db,
	}

}
