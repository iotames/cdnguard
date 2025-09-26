package db

func (d DB) GetOne(sqltxt string, dlist []any, args ...any) error {
	return d.edb.GetOne(sqltxt, dlist, args...)
}

func (d DB) GetMany(sqltxt string, dest any, args ...any) error {
	return d.edb.GetMany(sqltxt, dest, args...)
}

func (d DB) GetOneBySqlFile(filename string, dlist []any, args ...any) error {
	var err error
	var sqltxt string
	sqltxt, err = d.sqlDir.GetSQL(filename)
	if err != nil {
		return err
	}
	return d.GetOne(sqltxt, dlist, args...)
}

func (d DB) GetAllBySqlFile(filename string, dest any, args ...any) error {
	var err error
	var sqltxt string
	sqltxt, err = d.sqlDir.GetSQL(filename)
	if err != nil {
		return err
	}
	return d.GetMany(sqltxt, dest, args...)
}

func (d DB) GetAllBySqlFileReplace(filename string, dest any, args ...string) error {
	var err error
	var sqltxt string
	sqltxt, err = d.sqlDir.GetSQL(filename, args...)
	if err != nil {
		return err
	}
	return d.GetMany(sqltxt, dest)
}
