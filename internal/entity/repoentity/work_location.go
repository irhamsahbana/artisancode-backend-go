package repoentity

type WorkLocation struct {
	ID           string   `db:"id"`
	TenantID     string   `db:"tenant_id"`
	Name         string   `db:"name"`
	Address      *string  `db:"address"`
	Timezone     string   `db:"timezone"`
	Latitude     *float64 `db:"latitude"`
	Longitude    *float64 `db:"longitude"`
	RadiusMeters *int     `db:"radius_meters"`
}
