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
	pr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/profile"
	r "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	tra "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/transactors"
	usr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	q "github.com/lejeunel/go-image-annotator/modules/query"
	usrs "github.com/lejeunel/go-image-annotator/modules/user-store"
)

type Infra struct {
	im.ImageRepo
	clc.CollectionRepo
	an.AnnotationRepo
	lbl.LabelRepo
	grp.GroupRepo
	r.RoleRepo
	usr.UserRepo
	usrs.UserStore
	pr.ProfileRepo
	ev.EventRepo
	md.MetaRepo
	ImageFileStore fs.FileStore
	TempFileStore  fs.LocalFileStore
	q.FilterParser
	q.OrderParser
	*sqlx.DB
}

func BuildInfra(localPath string, imageStore fs.FileStore) Infra {
	filterParser, orderingParser := im.MakeQueryParsers()
	db := db.NewSQLiteConnection(db.DBPath(localPath))
	userRepo := usr.NewUserRepo(db)
	roleRepo := r.NewRoleRepo(db)
	groupRepo := grp.NewGroupRepo(db)
	localFileStore, err := fs.NewLocalFileStore(fmt.Sprintf("%v/%v", localPath, "tmp"))
	if err != nil {
		panic(err)
	}
	return Infra{
		im.NewImageRepo(db, filterParser, orderingParser),
		clc.NewCollectionRepo(db),
		an.NewAnnotationRepo(db),
		lbl.NewLabelRepo(db),
		groupRepo,
		roleRepo,
		userRepo,
		usrs.NewUserStore(userRepo, roleRepo, groupRepo),
		pr.NewProfileRepo(db),
		ev.NewEventRepo(db),
		md.NewMetaRepo(db),
		imageStore,
		*localFileStore,
		filterParser,
		orderingParser,
		db,
	}
}

func BuildImageStoreFromInfra(infra Infra) ims.ImageStore {
	return ims.New(
		ims.Repos{
			ImageRepo:      infra.ImageRepo,
			CollectionRepo: infra.CollectionRepo,
			AnnotationRepo: infra.AnnotationRepo,
			MetaRepo:       infra.MetaRepo,
		},
		tra.NewStoreTransactor(infra.DB, infra.FilterParser, infra.OrderStrParser),
		infra.ImageFileStore, infra.FilterParser, infra.OrderParser)
}
