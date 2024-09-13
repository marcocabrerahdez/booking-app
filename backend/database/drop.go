package database

func Drop() error {
	if err := DB.Exec("DROP TABLE IF EXISTS users CASCADE").Error; err != nil {
		return err
	}

	return nil
}
