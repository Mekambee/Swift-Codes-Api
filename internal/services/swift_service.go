package services

import (
	"fmt"

	"github.com/Mekambee/Swift-Codes-Api/internal/database"
	"github.com/Mekambee/Swift-Codes-Api/internal/models"
)

func SaveSwiftCodes(data []models.SwiftCodeData) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
      INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
      VALUES ($1, $2, $3, $4, $5, $6)
      ON CONFLICT (swift_code) DO NOTHING
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, sc := range data {
		_, err = stmt.Exec(
			sc.SwiftCode,
			sc.BankName,
			sc.Address,
			sc.CountryISO2,
			sc.CountryName,
			sc.IsHeadquarter,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetSwiftCode(swiftCode string) (*models.SwiftCodeData, error) {
	query := `
      SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter
      FROM swift_codes
      WHERE swift_code = $1
      LIMIT 1
    `
	row := database.DB.QueryRow(query, swiftCode)
	var sc models.SwiftCodeData
	err := row.Scan(
		&sc.ID,
		&sc.SwiftCode,
		&sc.BankName,
		&sc.Address,
		&sc.CountryISO2,
		&sc.CountryName,
		&sc.IsHeadquarter,
	)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func GetBranchesByHQ(swiftHQ string) ([]models.SwiftCodeData, error) {
	base := swiftHQ[0:8]
	query := `
        SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter
        FROM swift_codes
        WHERE LEFT(swift_code, 8) = $1 AND swift_code != $2
    `
	rows, err := database.DB.Query(query, base, swiftHQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.SwiftCodeData{}

	for rows.Next() {
		var sc models.SwiftCodeData
		err := rows.Scan(
			&sc.ID,
			&sc.SwiftCode,
			&sc.BankName,
			&sc.Address,
			&sc.CountryISO2,
			&sc.CountryName,
			&sc.IsHeadquarter,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	return result, nil
}

func GetSwiftByCountryISO2(iso2 string) ([]models.SwiftCodeData, error) {
	query := `
      SELECT id, swift_code, bank_name, address, country_iso2, country_name, is_headquarter
      FROM swift_codes
      WHERE country_iso2 = $1
    `
	rows, err := database.DB.Query(query, iso2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.SwiftCodeData
	for rows.Next() {
		var sc models.SwiftCodeData
		err := rows.Scan(
			&sc.ID,
			&sc.SwiftCode,
			&sc.BankName,
			&sc.Address,
			&sc.CountryISO2,
			&sc.CountryName,
			&sc.IsHeadquarter,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	return result, nil
}

func DeleteSwiftCode(swiftCode, bankName, iso2 string) error {
	query := `
      DELETE FROM swift_codes
      WHERE swift_code = $1 AND bank_name = $2 AND country_iso2 = $3
    `
	res, err := database.DB.Exec(query, swiftCode, bankName, iso2)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no matching record found")
	}
	return nil
}
