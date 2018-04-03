package crdb

type Update interface {
	IsUpdating() bool
	Value() interface{}
}

type StringUpdate struct {
	Updating bool
	NewVal   string
}

func (su StringUpdate) IsUpdating() bool   { return su.Updating }
func (su StringUpdate) Value() interface{} { return su.NewVal }

type BoolUpdate struct {
	Updating bool
	NewVal   bool
}

func (su BoolUpdate) IsUpdating() bool   { return su.Updating }
func (su BoolUpdate) Value() interface{} { return su.NewVal }

type Int64Update struct {
	Updating bool
	NewVal   int64
}

func (su Int64Update) IsUpdating() bool   { return su.Updating }
func (su Int64Update) Value() interface{} { return su.NewVal }

type Float64Update struct {
	Updating bool
	NewVal   float64
}

func (su Float64Update) IsUpdating() bool   { return su.Updating }
func (su Float64Update) Value() interface{} { return su.NewVal }

// AppendUpdate adds the value of upd to updArgs and updCol (column to be updated) to updCols
// if upd IsUpdating.
// The result of appending to updCols and updArgs are returned respectively.
// If upd IsUpdating is false, updCols and updArgs are returned unmodified.
func AppendUpdate(upd Update, updCol, updCols string, updArgs []interface{}) (string, []interface{}) {
	if upd.IsUpdating() {
		updCols = ColDesc(updCols, updCol)
		updArgs = append(updArgs, upd.Value())
	}
	return updCols, updArgs
}
