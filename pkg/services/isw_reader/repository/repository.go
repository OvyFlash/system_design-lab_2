package repository

import (
	"lab_2/internal/databases/postgres"
	"lab_2/pkg/services/isw_reader/models"
)

type Repository struct {
	*postgres.Storage
}

func NewRepository(p *postgres.Storage) (*Repository, error) {
	r := &Repository{
		Storage: p,
	}
	err := r.Write.AutoMigrate(&models.RawISWPage{}, &models.ParsedISWPage{})
	return r, err
}

func (r *Repository) WriteRawData(page models.RawISWPage) error {
	return r.Write.Create(&page).Error
}

func (r *Repository) GetPresentPagesDates() (res []models.RawISWPage, err error) {
	err = r.Read.Select("date").Find(&res).Error
	return
}

func (r *Repository) GetUnprocessedPages() (res []models.RawISWPage, err error) {
	err = r.Read.Where("date NOT IN (?)", r.Read.Select("date").
		Model(&models.ParsedISWPage{})).
		Find(&res).Error
	return
}

func (r *Repository) GetLastPage() (res models.RawISWPage, err error) {
	err = r.Read.Select("date").Order("date DESC").Take(&res).Error
	return
}

func (r *Repository) WriteParsedPage(page models.ParsedISWPage) error {
	return r.Write.Create(&page).Error
}
