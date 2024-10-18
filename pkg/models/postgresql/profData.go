package postgresql

import (
	"database/sql"
	"dem3_demo_v2/pkg/models"
	"errors"
	"log"
)

type ProfDataModel struct {
	DB *sql.DB
}

func (prof *ProfDataModel) InsertProfData(
	D2Number,
	D2Profstroi,
	D2Object,
	D2Manager,
	D2Kontragent,
	D2ID,
	D2Diler,
	D2City,
	D2Napr,
	D2SumProjToSk,
	D2SumSkidka,
	D2SumProjWithSkidka,
	D2SumConstrWithSkidka,
	D2SumRabWithSkidka,
	D2Status,
	NoteOrder string) (int, error) {

	tx, err := prof.DB.Begin()
	if err != nil {
		return 0, err
	}

	stmtInsert := `INSERT INTO prof_data (d2_number, d2_profstroi, d2_object, d2_manager, d2_kontragent, d2_id, d2_diler, d2_city, d2_napr, d2_sum_proj_to_sk, d2_sum_skidka, d2_sum_proj_with_skidka,
	d2_sum_constr_with_skidka, d2_sum_rab_with_skidka, d2_status, note_order)
	VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING ID`

	stmtSelect := `SELECT id, d2_number FROM prof_data WHERE d2_number = $1`

	stmtUpdate := `UPDATE prof_data SET d2_number = $1, d2_profstroi = $2, d2_object = $3, d2_manager = $4, d2_kontragent = $5, d2_id = $6, d2_diler = $7, d2_city = $8, d2_napr = $9,
	               d2_sum_proj_to_sk = $10, d2_sum_skidka = $11, d2_sum_proj_with_skidka = $12, d2_sum_constr_with_skidka = $13, d2_sum_rab_with_skidka = $14, d2_status = $15, 
	               note_order = $16 WHERE id = $17 RETURNING id`

	row := tx.QueryRow(stmtSelect, D2Number)

	s := &models.ProfData{}

	err = row.Scan(&s.ID, &s.D2Number)

	var id int
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := tx.QueryRow(stmtInsert, D2Number, D2Profstroi, D2Object, D2Manager, D2Kontragent, D2ID, D2Diler, D2City, D2Napr, D2SumProjToSk, D2SumSkidka, D2SumProjWithSkidka,
				D2SumConstrWithSkidka, D2SumRabWithSkidka, D2Status, NoteOrder).Scan(&id)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	} else {
		err = tx.QueryRow(stmtUpdate, D2Number, D2Profstroi, D2Object, D2Manager, D2Kontragent, D2ID, D2Diler, D2City, D2Napr, D2SumProjToSk, D2SumSkidka, D2SumProjWithSkidka,
			D2SumConstrWithSkidka, D2SumRabWithSkidka, D2Status, NoteOrder, &s.ID).Scan(&id)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	return id, err
}

func (prof *ProfDataModel) InsertDemMaterial(id int, size, name, count, allowances, color, height string) {
	stmtInsertDet := "INSERT INTO dem_klaes_materials(order_id, size, name, count, allowances, color, height) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := prof.DB.Exec(stmtInsertDet, id, size, name, count, allowances, color, height)
	if err != nil {
		log.Println(err)
	}
}
