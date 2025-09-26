package contract

type ISqlDir interface {
	GetSQL(fpath string, replaceList ...string) (string, error)
}
