package collections

import (
	pro "datahub/domain/annotation_profiles"
	e "datahub/errors"
	g "datahub/generic"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

type CollectionRecord struct {
	Id          CollectionId             `db:"id"`
	Name        string                   `db:"name"`
	Description string                   `db:"description"`
	Group       string                   `db:"group_name"`
	CreatedAt   time.Time                `db:"created_at"`
	UpdatedAt   time.Time                `db:"updated_at"`
	ProfileId   *pro.AnnotationProfileId `db:"profile_id"`
}

type SQLCollectionRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
}

func NewSQLiteCollectionRepo(db *sqlx.DB) *SQLCollectionRepo {

	return &SQLCollectionRepo{Db: db, ErrorConverter: &e.SQLiteErrorConverter{}}

}

func NewPostgreSQLCollectionRepo(db *sqlx.DB) *SQLCollectionRepo {

	return &SQLCollectionRepo{Db: db, ErrorConverter: &e.PostgreSQLErrorConverter{}}

}

func (r *SQLCollectionRepo) Create(collection *Collection) error {

	query := "INSERT INTO collections (id,name,description,group_name,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.Db.Exec(query, collection.Id, collection.Name,
		collection.Description, collection.Group, collection.CreatedAt, collection.UpdatedAt)

	if err != nil {
		err = r.ErrorConverter.Convert(err)
		if errors.Is(err, e.ErrDBUniqueConstraint) {
			err = e.ErrDuplication
		}

		return err
	}

	return nil
}

func FromRecord(record CollectionRecord) Collection {
	return Collection{
		Id:          record.Id,
		Name:        record.Name,
		Description: record.Description,
		Group:       record.Group,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
		ProfileId:   record.ProfileId,
	}
}

func (r *SQLCollectionRepo) Find(id CollectionId) (*Collection, error) {
	record := CollectionRecord{}
	if err := r.Db.Get(&record, "SELECT id,name,description,group_name,created_at,updated_at,profile_id FROM collections WHERE id=$1", id); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}

	collection := FromRecord(record)
	return &collection, nil
}

func (r *SQLCollectionRepo) GetByName(name string) (*Collection, error) {
	record := CollectionRecord{}
	err := r.Db.Get(&record, "SELECT id,name,description,group_name,created_at,updated_at,profile_id FROM collections WHERE name=$1", name)

	if err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}

	collection := FromRecord(record)
	return &collection, nil
}

func (r *SQLCollectionRepo) Touch(id CollectionId, now time.Time) error {
	_, err := r.Db.Exec("UPDATE collections SET updated_at=$1 WHERE id=$2",
		now, id)

	if err != nil {
		return e.ErrDB
	}

	return nil

}

func (r *SQLCollectionRepo) Update(id CollectionId, c CollectionUpdatables) error {
	_, err := r.Db.Exec("UPDATE collections SET name=$1,description=$2,group_name=$3 WHERE id=$4",
		c.Name, c.Description, c.Group, id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLCollectionRepo) Delete(collection *Collection) error {
	id := collection.Id.String()
	_, err := r.Db.Exec("DELETE FROM collections WHERE id=$1", id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLCollectionRepo) Count() (int64, error) {

	var count int64

	query := "SELECT COUNT(*) FROM collections"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (r *SQLCollectionRepo) List(ordering OrderingArgs, pagination g.PaginationParams) ([]Collection, *g.PaginationMeta, error) {

	query := g.SqlBuilder.Select("*").From("collections")
	if ordering.Name == true {
		query = query.OrderBy("name")
	}

	query = query.Limit(uint64(pagination.Limit())).Offset(uint64(pagination.Offset()))
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("building query: %w", e.ErrDB)
	}

	var records []CollectionRecord
	err = r.Db.Select(&records, sql, args...)
	if err != nil {
		return nil, nil, e.ErrDB
	}

	count, err := r.Count()
	if err != nil {
		return nil, nil, err
	}
	paginationMeta := g.NewPaginationMeta(pagination.Page, count, int64(pagination.PageSize))

	collections := []Collection{}
	for _, record := range records {
		collection := FromRecord(record)
		collections = append(collections, collection)
	}

	return collections, &paginationMeta, nil

}

func (r *SQLCollectionRepo) AssignProfile(collection *Collection, profile *pro.AnnotationProfile) error {
	query := g.SqlBuilder.Update("collections").
		Set("profile_id", profile.Id).Where(sq.Eq{"id": collection.Id})
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", e.ErrDB)
	}

	_, err = r.Db.Exec(sql, args...)
	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLCollectionRepo) UnassignProfile(collection *Collection) error {
	_, err := r.Db.Exec("UPDATE collections SET profile_id=NULL WHERE id=$1", collection.Id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}
