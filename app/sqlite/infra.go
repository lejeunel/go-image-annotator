package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	db "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	an "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	ev "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/event"
	grp "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	im "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	lbl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	md "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/metadata"
	r "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	usr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
)

type Infra struct {
	im.ImageRepo
	clc.CollectionRepo
	an.AnnotationRepo
	lbl.LabelRepo
	grp.GroupRepo
	r.RoleRepo
	usr.UserRepo
	ev.EventRepo
	md.MetaRepo
	ImageFileStore  fs.FileStore
	TempFileStore   fs.LocalFileStore
	PolicyFileStore fs.FileStore
	qu.FilterStrParser
	qu.OrderingStrConverter
	*sqlx.DB
}

func BuildInfra(path string) Infra {
	filterParser, orderingParser := im.MakeQueryParsers()
	db := db.NewSQLiteDB(fmt.Sprintf("%v/%v", path, "db.sqlite"))
	return Infra{
		im.NewImageRepo(db, filterParser, orderingParser),
		clc.NewCollectionRepo(db),
		an.NewAnnotationRepo(db),
		lbl.NewLabelRepo(db),
		grp.NewGroupRepo(db),
		r.NewRoleRepo(db),
		usr.NewUserRepo(db),
		ev.NewEventRepo(db),
		md.NewMetaRepo(db),
		fs.NewLocalFileStore(fmt.Sprintf("%v/%v", path, "images")),
		fs.NewLocalFileStore(fmt.Sprintf("%v/%v", path, "tmp")),
		fs.NewLocalFileStore(fmt.Sprintf("%v/%v", path, "assets")),
		filterParser,
		orderingParser,
		db,
	}

}
