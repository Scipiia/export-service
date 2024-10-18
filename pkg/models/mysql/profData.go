package mysql

import (
	"database/sql"
	"dem3_demo_v2/pkg/models"
	"errors"
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
	NoteOrder string) error {

	tx, err := prof.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmtInsert := `INSERT INTO prof_data (d2_number, d2_profstroi, d2_object, d2_manager, d2_kontragent, d2_id, d2_diler, d2_city, d2_napr, d2_sum_proj_to_sk, d2_sum_skidka, d2_sum_proj_with_skidka,
	d2_sum_constr_with_skidka, d2_sum_rab_with_skidka, d2_status, note_order)
	VALUES(?, ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	stmtSelect := `SELECT id, d2_number FROM prof_data WHERE d2_number = ?`

	stmtUpdate := `UPDATE prof_data SET d2_number = ?, d2_profstroi = ?, d2_object = ?, d2_manager = ?, d2_kontragent = ?, d2_id = ?, d2_diler = ?, d2_city = ?, d2_napr = ?, 
                     d2_sum_proj_to_sk = ?, d2_sum_skidka = ?, d2_sum_proj_with_skidka = ?, d2_sum_constr_with_skidka = ?, d2_sum_rab_with_skidka = ?, d2_status = ?, note_order = ? WHERE id = ?`

	//stmtInsertDet := "INSERT INTO dem_klaes_materials(name) VALUES (?)"

	row := tx.QueryRow(stmtSelect, D2Number)

	s := &models.ProfData{}

	err = row.Scan(&s.ID, &s.D2Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := tx.Exec(stmtInsert, D2Number, D2Profstroi, D2Object, D2Manager, D2Kontragent, D2ID, D2Diler, D2City, D2Napr, D2SumProjToSk, D2SumSkidka, D2SumProjWithSkidka,
				D2SumConstrWithSkidka, D2SumRabWithSkidka, D2Status, NoteOrder)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err = tx.Exec(stmtUpdate, D2Number, D2Profstroi, D2Object, D2Manager, D2Kontragent, D2ID, D2Diler, D2City, D2Napr, D2SumProjToSk, D2SumSkidka, D2SumProjWithSkidka,
				D2SumConstrWithSkidka, D2SumRabWithSkidka, D2Status, NoteOrder, &s.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	return err
}

func (prof *ProfDataModel) InsertDemMaterial(id int, size, name, count string) {
	stmtInsertDet := "INSERT INTO dem_klaes_materials(order_id, size, name, count, allowances, color, height) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := prof.DB.Exec(stmtInsertDet, id, size, name, count)
	if err != nil {
		return
	}
}
