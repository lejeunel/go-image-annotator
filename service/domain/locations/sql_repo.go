package locations

import (
	e "datahub/errors"
	g "datahub/generic"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SQLLocationRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
}

type SiteRecord struct {
	Id        SiteId    `db:"id"`
	Name      string    `db:"name"`
	Group     string    `db:"group_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CameraRecord struct {
	Id          CameraId  `db:"id"`
	SiteId      SiteId    `db:"site_id"`
	SiteName    string    `db:"site_name"`
	Name        string    `db:"name"`
	Transmitter string    `db:"transmitter"`
	Group       string    `db:"group_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewSQLiteLocationRepo(db *sqlx.DB) *SQLLocationRepo {

	return &SQLLocationRepo{Db: db,
		ErrorConverter: &e.SQLiteErrorConverter{}}

}

func CameraFromRecord(record CameraRecord) Camera {
	return Camera{Id: record.Id,
		Name:        record.Name,
		Group:       record.Group,
		Transmitter: record.Transmitter,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.CreatedAt}

}

func (r *SQLLocationRepo) siteWithNameExists(siteName string) error {
	var count int64
	query := "SELECT COUNT(*) FROM sites WHERE name=$1"
	if err := r.Db.QueryRow(query, siteName).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("creating site with name %v: %w", siteName, e.ErrDuplication)
	}
	return nil
}

func (r *SQLLocationRepo) CreateSite(site *Site) error {

	err := r.siteWithNameExists(site.Name)
	if err != nil {
		return err
	}

	query := "INSERT INTO sites (id, name, group_name, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)"
	_, err = r.Db.Exec(query, site.Id, site.Name, site.Group, site.CreatedAt, site.UpdatedAt)

	if err != nil {
		return fmt.Errorf("creating site: %w", err)
	}

	return nil
}

func (r *SQLLocationRepo) cameraWithNameExists(camera *Camera) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM cameras WHERE name=$1 AND site_id=$2"
	if err := r.Db.QueryRow(query, camera.Name, camera.Site.Id).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r *SQLLocationRepo) CreateCamera(camera *Camera) error {

	hasDuplicate, err := r.cameraWithNameExists(camera)
	if err != nil {
		return err
	}
	if hasDuplicate {
		return fmt.Errorf("creating camera with name %v: %w", camera.Name, e.ErrDuplication)
	}

	query := "INSERT INTO cameras (id,name,site_id,transmitter,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err = r.Db.Exec(query, camera.Id, camera.Name, camera.Site.Id, camera.Transmitter,
		camera.CreatedAt, camera.UpdatedAt)
	if err != nil {
		err = r.ErrorConverter.Convert(err)
		if errors.Is(err, e.ErrDBForeignKeyConstraint) {
			return fmt.Errorf("creating camera in site %v: %w", camera.Site.Id, e.ErrDependency)
		}
		return err
	}

	return nil
}

func (r *SQLLocationRepo) UpdateSite(site *Site) error {
	_, err := r.Db.Exec("UPDATE sites SET name=$1,group_name=$2 WHERE id=$3",
		site.Name, site.Group, site.Id)

	if err != nil {
		return fmt.Errorf("Updating site to %+v: %w", site, err)
	}

	return nil
}

func (r *SQLLocationRepo) FindCamera(id CameraId) (*Camera, error) {
	record := CameraRecord{}
	query := `SELECT c.id,
		c.site_id,
		c.name,
		c.transmitter,
		c.created_at,
		c.updated_at,
		si.name AS site_name,
		si.group_name AS group_name
		FROM cameras AS c
		JOIN sites AS si ON (c.site_id=si.id)
		WHERE c.id=$1
`
	if err := r.Db.Get(&record, query, id); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}

	camera := CameraFromRecord(record)
	site, err := r.FindSite(record.SiteId)
	if err != nil {
		return nil, err
	}
	camera.Site = site
	return &camera, nil
}

func (r *SQLLocationRepo) FindSite(id SiteId) (*Site, error) {

	site := SiteRecord{}
	if err := r.Db.Get(&site, "SELECT * FROM sites WHERE id=$1", id); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}
	return &Site{Id: site.Id,
		Name:      site.Name,
		Group:     site.Group,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}, nil
}

func (r *SQLLocationRepo) FindSiteByName(name string) (*Site, error) {
	site := SiteRecord{}
	if err := r.Db.Get(&site, "SELECT * FROM sites WHERE name=$1", name); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}

	return &Site{Id: site.Id,
		Name:      site.Name,
		Group:     site.Group,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}, nil
}

func (r *SQLLocationRepo) applyFilter(query sq.SelectBuilder, filters FilterArgs) sq.SelectBuilder {
	if filters.Group != nil {
		query = query.Where(sq.Eq{"s.group_name": *filters.Group})
	}

	if filters.Collection != nil {
		query = query.
			Join("cameras AS ca ON ca.site_id = s.id").
			Join("images AS i ON i.camera_id = ca.id").
			Join("image_collection_assoc AS ic ON ic.image_id = i.id").
			Join("collections AS c ON c.id = ic.collection_id").
			Where(sq.Eq{"c.name": *filters.Collection})
	}

	return query

}

func (r *SQLLocationRepo) List(filters FilterArgs, ordering OrderingArgs, pagination g.PaginationParams) ([]Site, *g.PaginationMeta, error) {
	query := NewBaseSiteQuery()
	query = r.applyFilter(query, filters)

	if ordering.Name == true {
		query = query.OrderBy("s.name")
	}

	query = query.Limit(uint64(pagination.Limit())).Offset(uint64(pagination.Offset()))
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("building query: %w", e.ErrDB)
	}

	siteRecords := []SiteRecord{}
	if err := r.Db.Select(&siteRecords, sql, args...); err != nil {
		return nil, nil, r.ErrorConverter.Convert(err)
	}

	count, err := r.NumSites(filters)
	if err != nil {
		return nil, nil, err
	}

	paginationMeta := g.NewPaginationMeta(pagination.Page, count, int64(pagination.PageSize))
	if err != nil {
		return nil, nil, err
	}

	sites := []Site{}
	for _, record := range siteRecords {
		sites = append(sites, Site{Id: record.Id,
			Name:      record.Name,
			Group:     record.Group,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt})
	}

	return sites, &paginationMeta, nil
}

func (r *SQLLocationRepo) NumSites(filters FilterArgs) (int64, error) {
	countQuery := g.SqlBuilder.Select("COUNT(DISTINCT s.id) AS site_count").From("sites AS s")
	countQuery = r.applyFilter(countQuery, filters)

	sql, args, err := countQuery.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query: %w", e.ErrDB)
	}

	var count int64
	if err := r.Db.Get(&count, sql, args...); err != nil {
		return 0, r.ErrorConverter.Convert(err)
	}

	return count, nil

}

func (r *SQLLocationRepo) ListCamerasOfSite(site *Site) ([]*Camera, error) {
	query := g.SqlBuilder.Select("*").From("cameras").Where(sq.Eq{"site_id": site.Id}).
		OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, e.ErrDB
	}

	var cameraRecords []*CameraRecord
	if err := r.Db.Select(&cameraRecords, sql, args...); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}
	var cameras []*Camera
	for _, c := range cameraRecords {
		camera := CameraFromRecord(*c)
		cameras = append(cameras, &camera)
	}

	return cameras, nil

}
func (r *SQLLocationRepo) FindCameraByName(site *Site, name string) (*Camera, error) {
	query := g.SqlBuilder.Select("*").From("cameras").Where(
		sq.And{
			sq.Eq{"site_id": site.Id},
			sq.Eq{"name": name}},
	).OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, e.ErrDB
	}

	record := CameraRecord{}
	if err := r.Db.Get(&record, sql, args...); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}
	camera := CameraFromRecord(record)
	return &camera, nil

}

func (r *SQLLocationRepo) DeleteSite(id SiteId) error {

	_, err := r.Db.Exec("DELETE FROM sites WHERE id=$1", id)
	return err
}

func (r *SQLLocationRepo) DeleteCamera(id CameraId) error {

	_, err := r.Db.Exec("DELETE FROM cameras WHERE id=$1", id)
	return err
}

func (r *SQLLocationRepo) UpdateCamera(id CameraId, camera CameraUpdatables) (*Camera, error) {
	query := "UPDATE cameras SET name=$1, transmitter=$2 WHERE id=$3"
	_, err := r.Db.Exec(query, camera.Name, camera.Transmitter, id)

	if err != nil {
		return nil, fmt.Errorf("updating camera name: %w", err)
	}

	site, err := r.FindSiteByName(camera.SiteName)
	if err != nil {
		return nil, fmt.Errorf("updating site of camera to %v: fetching by name: %w", camera.SiteName, err)
	}
	query = "UPDATE cameras SET site_id=$1 WHERE id=$2"
	_, err = r.Db.Exec(query, site.Id, id)

	if err != nil {
		return nil, fmt.Errorf("updating site of camera name: updating field: %w", err)
	}

	updatedCamera, err := r.FindCamera(id)
	if err != nil {
		return nil, fmt.Errorf("updating camera: %w", err)
	}

	return updatedCamera, nil

}
