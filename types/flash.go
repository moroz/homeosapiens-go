package types

//go:generate stringer -type=FlashLevel -trimprefix=FlashLevel_
type FlashLevel int

const (
	FlashLevel_Info FlashLevel = iota
	FlashLevel_Success
	FlashLevel_Danger
)

type FlashMessage struct {
	Level   FlashLevel
	Message string
}
