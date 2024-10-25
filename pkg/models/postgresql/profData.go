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

func (prof *ProfDataModel) Insertik(data *models.ProfData, maps map[int][]string) error {

	stmtInsert := `INSERT INTO prof_data (d2_number, d2_profstroi, d2_object, d2_manager, d2_kontragent, d2_id, d2_diler, d2_city, d2_napr, d2_sum_proj_to_sk, d2_sum_skidka, d2_sum_proj_with_skidka,
	d2_sum_constr_with_skidka, d2_sum_rab_with_skidka, d2_status, note_order)
	VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING ID`

	stmtSelect := `SELECT id, d2_number FROM prof_data WHERE d2_number = $1`

	stmtUpdate := `UPDATE prof_data SET d2_number = $1, d2_profstroi = $2, d2_object = $3, d2_manager = $4, d2_kontragent = $5, d2_id = $6, d2_diler = $7, d2_city = $8, d2_napr = $9,
	               d2_sum_proj_to_sk = $10, d2_sum_skidka = $11, d2_sum_proj_with_skidka = $12, d2_sum_constr_with_skidka = $13, d2_sum_rab_with_skidka = $14, d2_status = $15, 
	               note_order = $16 WHERE id = $17 RETURNING id`

	stmtInsertDet := "INSERT INTO dem_klaes_materials(order_id, size, name, count, allowances, color, height) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	stmtDelete := "DELETE FROM dem_klaes_materials WHERE order_id = $1"

	//stmtUpdateDet := "UPDATE dem_klaes_materials SET size = $1, name = $2, count = $3, allowances = $4, color = $5, height = $6 WHERE order_id = $7 AND name = $8"

	tx, err := prof.DB.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(stmtSelect, &data.D2Number)
	err = row.Scan(&data.ID, &data.D2Number)

	var id int
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := tx.QueryRow(stmtInsert, &data.D2Number, &data.D2Profstroi, &data.D2Object, &data.D2Manager, &data.D2Kontragent, &data.D2ID, &data.D2Diler, &data.D2City, &data.D2Napr,
				&data.D2SumProjToSk, &data.D2SumSkidka, &data.D2SumProjWithSkidka, &data.D2SumConstrWithSkidka, &data.D2SumRabWithSkidka, &data.D2Status, &data.NoteOrder).Scan(&id)
			if err != nil {
				tx.Rollback()
				return err
			}

			//tx1, err := prof.DB.Begin()
			//if err != nil {
			//	return err
			//}
			for _, strings := range maps {
				data.Details.Size = strings[0]
				data.Details.Name = strings[1]
				data.Details.Count = strings[2]
				data.Details.Allowances = strings[4]
				data.Details.Color = strings[5]
				data.Details.Height = strings[6]

				_, err := tx.Exec(stmtInsertDet, id, data.Details.Size, data.Details.Name, data.Details.Count, data.Details.Allowances, data.Details.Color, data.Details.Height)
				if err != nil {
					tx.Rollback()
					log.Println("tx", err)
				}
			}
			//err = tx1.Commit()
			//if err != nil {
			//	tx.Rollback()
			//	return err
			//}
		}
	} else {
		err = tx.QueryRow(stmtUpdate, &data.D2Number, &data.D2Profstroi, &data.D2Object, &data.D2Manager, &data.D2Kontragent, &data.D2ID, &data.D2Diler, &data.D2City, &data.D2Napr,
			&data.D2SumProjToSk, &data.D2SumSkidka, &data.D2SumProjWithSkidka, &data.D2SumConstrWithSkidka, &data.D2SumRabWithSkidka, &data.D2Status, &data.NoteOrder, &data.ID).Scan(&id)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, strings := range maps {
			data.Details.Size = strings[0]
			data.Details.Name = strings[1]
			data.Details.Count = strings[2]
			data.Details.Allowances = strings[4]
			data.Details.Color = strings[5]
			data.Details.Height = strings[6]

			//fmt.Println("RR", data.Details.Name)

			begin, err := prof.DB.Begin()
			if err != nil {
				return err
			}
			_, err = begin.Exec(stmtDelete, &id)
			if err != nil {
				begin.Rollback()
				return err
			}

			begin.Commit()
			//return err

			_, err = tx.Exec(stmtInsertDet, id, data.Details.Size, data.Details.Name, data.Details.Count, data.Details.Allowances, data.Details.Color, data.Details.Height)
			if err != nil {
				tx.Rollback()
				log.Println("tx", err)
			}
		}
	}

	err = tx.Commit()
	return err
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
