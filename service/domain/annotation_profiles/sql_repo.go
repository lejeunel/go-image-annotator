package annotation_profile

import (
	"database/sql"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLAnnotationProfileRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
}

func NewSQLiteAnnotationProfileRepo(db *sqlx.DB) *SQLAnnotationProfileRepo {

	return &SQLAnnotationProfileRepo{Db: db, ErrorConverter: &e.SQLiteErrorConverter{}}

}

func (r *SQLAnnotationProfileRepo) Save(profile *AnnotationProfile) error {
	query := "INSERT INTO annotation_profiles (id,name) VALUES ($1,$2)"

	_, err := r.Db.Exec(query, profile.Id, profile.Name)
	if err != nil {
		err = r.ErrorConverter.Convert(err)
		if errors.Is(err, e.ErrDBUniqueConstraint) {
			err = e.ErrDuplication
		}

		return err
	}

	return nil
}

func (r *SQLAnnotationProfileRepo) Delete(profile *AnnotationProfile) error {
	_, err := r.Db.Exec("DELETE FROM annotation_profiles WHERE id=$1", profile.Id)
	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationProfileRepo) Find(id AnnotationProfileId) (*AnnotationProfile, error) {

	profile := AnnotationProfile{}
	err := r.Db.Get(&profile, "SELECT id,name FROM annotation_profiles WHERE id=$1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrDB
	}

	return &profile, nil
}

func (r *SQLAnnotationProfileRepo) NumProfiles() (int64, error) {

	var count int64
	query := g.SqlBuilder.Select("COUNT(*)").From("annotation_profiles")
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, e.ErrDB
	}
	if err := r.Db.QueryRow(sql, args...).Scan(&count); err != nil {
		return 0, e.ErrDB
	}

	return count, nil

}

func (r *SQLAnnotationProfileRepo) List(pagination g.PaginationParams) ([]AnnotationProfile, *g.PaginationMeta, error) {
	query := g.SqlBuilder.Select(`id,name`).From("annotation_profiles").OrderBy("name")
	query = query.Limit(uint64(pagination.Limit())).Offset(uint64(pagination.Offset()))

	var profiles []AnnotationProfile
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("building SQL statement: %w", e.ErrDB)
	}
	if err := r.Db.Select(&profiles, sql, args...); err != nil {
		return nil, nil, fmt.Errorf("selecting data: %w", e.ErrDB)
	}

	count, err := r.NumProfiles()
	if err != nil {
		return nil, nil, fmt.Errorf("building profiles pagination meta-data: %w", err)
	}
	paginationMeta := g.NewPaginationMeta(pagination.Page, count, int64(pagination.PageSize))
	return profiles, &paginationMeta, nil
}

func (r *SQLAnnotationProfileRepo) FindByName(name string) (*AnnotationProfile, error) {

	profile := AnnotationProfile{}
	err := r.Db.Get(&profile, "SELECT id,name FROM annotation_profiles WHERE name=$1", name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrDB
	}

	return &profile, nil
}

func (r *SQLAnnotationProfileRepo) AddLabel(profile *AnnotationProfile, label *lbl.Label) error {
	query := "INSERT INTO annotation_profile_label_assoc (id,profile_id,label_id) VALUES ($1,$2,$3)"

	_, err := r.Db.Exec(query, uuid.New(), profile.Id, label.Id)
	if err != nil {
		return e.ErrDB
	}

	return nil
}
func (r *SQLAnnotationProfileRepo) ClearLabels(profile *AnnotationProfile) error {
	query := "DELETE FROM annotation_profile_label_assoc WHERE profile_id=$1"

	_, err := r.Db.Exec(query, profile.Id)
	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationProfileRepo) RemoveLabel(profile *AnnotationProfile, label *lbl.Label) error {
	query := "DELETE FROM annotation_profile_label_assoc WHERE label_id=$1 AND profile_id=$2"

	_, err := r.Db.Exec(query, label.Id, profile.Id)
	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationProfileRepo) GetLabelIds(profile *AnnotationProfile) ([]lbl.LabelId, error) {
	var labelIds []lbl.LabelId

	if err := r.Db.Select(&labelIds, "SELECT label_id FROM annotation_profile_label_assoc WHERE profile_id = $1 ORDER BY id", profile.Id); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}

	return labelIds, nil

}

func (r *SQLAnnotationProfileRepo) Rename(id AnnotationProfileId, name string) error {
	query := "UPDATE annotation_profiles SET name=$1 WHERE id=$2"

	_, err := r.Db.Exec(query, name, id)
	if err != nil {
		return e.ErrDB
	}

	return nil
}
