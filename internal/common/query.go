package common

// QueryOptions opciones genéricas para consultas
type QueryOptions struct {
	// Paginación
	Page   int `form:"page,default=1" binding:"min=1"`
	Limit  int `form:"limit,default=20" binding:"min=1,max=100"`
	Offset int `form:"-"` // Calculado automáticamente

	// Ordenamiento
	OrderBy  string `form:"order_by,default=created_at"`
	OrderDir string `form:"order_dir,default=desc" binding:"omitempty,oneof=asc desc"`

	// Búsqueda
	Search string `form:"search"`

	// Filtros genéricos
	Filters map[string]interface{} `form:"-"`

	// Control de preloads
	Preloads []string `form:"-"`

	// Filtros de usuario
	UserContext *UserContext `form:"-"`
}

// Validate valida y normaliza las opciones
func (qo *QueryOptions) Validate() error {
	if qo.Page < 1 {
		qo.Page = 1
	}
	if qo.Limit < 1 || qo.Limit > 100 {
		qo.Limit = 20
	}
	qo.Offset = (qo.Page - 1) * qo.Limit

	if qo.OrderDir != "asc" && qo.OrderDir != "desc" {
		qo.OrderDir = "desc"
	}

	if qo.Filters == nil {
		qo.Filters = make(map[string]interface{})
	}

	return nil
}

// AddFilter agrega un filtro
func (qo *QueryOptions) AddFilter(key string, value interface{}) {
	if qo.Filters == nil {
		qo.Filters = make(map[string]interface{})
	}
	qo.Filters[key] = value
}

// GetFilter obtiene un filtro
func (qo *QueryOptions) GetFilter(key string) (interface{}, bool) {
	if qo.Filters == nil {
		return nil, false
	}
	value, exists := qo.Filters[key]
	return value, exists
}

// PaginationMeta metadatos de paginación
type PaginationMeta struct {
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Total    int64 `json:"total"`
	Pages    int64 `json:"pages"`
	HasNext  bool  `json:"has_next"`
	HasPrev  bool  `json:"has_prev"`
	NextPage *int  `json:"next_page,omitempty"`
	PrevPage *int  `json:"prev_page,omitempty"`
}

// NewPaginationMeta crea metadatos de paginación
func NewPaginationMeta(page, limit int, total int64) *PaginationMeta {
	pages := (total + int64(limit) - 1) / int64(limit)
	if pages == 0 {
		pages = 1
	}

	meta := &PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   pages,
		HasNext: int64(page) < pages,
		HasPrev: page > 1,
	}

	if meta.HasNext {
		next := page + 1
		meta.NextPage = &next
	}

	if meta.HasPrev {
		prev := page - 1
		meta.PrevPage = &prev
	}

	return meta
}
