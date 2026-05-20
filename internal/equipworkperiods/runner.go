package equipworkperiods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"oiltech/internal/algokit"
)

type Runner struct {
	kit *algokit.Kit
	db  *sqlx.DB
}

func NewRunner(kit *algokit.Kit, db *sqlx.DB) *Runner {
	return &Runner{kit: kit, db: db}
}

type stateRow struct {
	FacilityID    string
	EquipmentID   string
	EquipmentType string
	State         string
	Quality       string
	Moment        time.Time
}

type workPeriod struct {
	FacilityID      string
	EquipmentID     string
	EquipmentType   string
	PeriodStart     time.Time
	PeriodEnd       time.Time
	State           string
	DurationSeconds int64
	Quality         string
}

func (r *Runner) Handle(w http.ResponseWriter, req *http.Request) {
	var runReq algokit.RunRequest
	if err := json.NewDecoder(req.Body).Decode(&runReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	transformation, err := r.kit.LoadTransformation(ctx, runReq.TransformationCode)
	if err != nil {
		writeError(w, err)
		return
	}

	rows, err := r.readStates(ctx, transformation.SourceDatasetCodes)
	if err != nil {
		writeError(w, err)
		return
	}

	periods := buildPeriods(rows)

	written, err := r.writePeriods(ctx, transformation.TargetDatasetCode, periods)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := algokit.RunResponse{
		Status:      "ok",
		RowsRead:    int64(len(rows)),
		RowsWritten: written,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (r *Runner) readStates(ctx context.Context, datasetCodes []string) ([]stateRow, error) {
	if len(datasetCodes) == 0 {
		return nil, fmt.Errorf("no source datasets")
	}
	code := datasetCodes[0]
	query := fmt.Sprintf(
		`select facility_id, equipment_id, equipment_type, state, quality, moment
		 from data."ds_%s"
		 order by facility_id, equipment_id, moment asc`,
		code,
	)
	sqlRows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	var result []stateRow
	for sqlRows.Next() {
		var sr stateRow
		if err := sqlRows.Scan(&sr.FacilityID, &sr.EquipmentID, &sr.EquipmentType, &sr.State, &sr.Quality, &sr.Moment); err != nil {
			return nil, err
		}
		result = append(result, sr)
	}
	return result, sqlRows.Err()
}

func buildPeriods(rows []stateRow) []workPeriod {
	type key struct {
		FacilityID  string
		EquipmentID string
	}
	groups := map[key][]stateRow{}
	for _, row := range rows {
		k := key{row.FacilityID, row.EquipmentID}
		groups[k] = append(groups[k], row)
	}

	var periods []workPeriod
	for _, group := range groups {
		sort.Slice(group, func(i, j int) bool {
			return group[i].Moment.Before(group[j].Moment)
		})
		for i := 0; i+1 < len(group); i++ {
			curr := group[i]
			next := group[i+1]
			dur := int64(next.Moment.Sub(curr.Moment).Seconds())
			if dur < 0 {
				dur = 0
			}
			periods = append(periods, workPeriod{
				FacilityID:      curr.FacilityID,
				EquipmentID:     curr.EquipmentID,
				EquipmentType:   curr.EquipmentType,
				PeriodStart:     curr.Moment,
				PeriodEnd:       next.Moment,
				State:           curr.State,
				DurationSeconds: dur,
				Quality:         curr.Quality,
			})
		}
	}
	return periods
}

func (r *Runner) writePeriods(ctx context.Context, targetCode string, periods []workPeriod) (int64, error) {
	if len(periods) == 0 {
		return 0, nil
	}
	var written int64
	for _, p := range periods {
		_, err := r.db.ExecContext(ctx,
			fmt.Sprintf(`insert into data."ds_%s"
				(facility_id, equipment_id, equipment_type, period_start, period_end, state, duration_seconds, quality)
				values ($1,$2,$3,$4,$5,$6,$7,$8)
				on conflict do nothing`, targetCode),
			p.FacilityID, p.EquipmentID, p.EquipmentType,
			p.PeriodStart, p.PeriodEnd, p.State, p.DurationSeconds, p.Quality,
		)
		if err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	resp := algokit.RunResponse{Status: "error", ErrorMessage: err.Error()}
	_ = json.NewEncoder(w).Encode(resp)
}
