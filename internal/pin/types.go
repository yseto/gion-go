package pin

type ReadFlag string

const (
	Unseen ReadFlag = "Unseen"
	Seen   ReadFlag = "Seen"
	Setpin ReadFlag = "Setpin"
)
