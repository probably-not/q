package consts

const (
	CommitPopU32              = uint32(0x10000)
	PushOverflowCheckU32      = uint32(0x8000)
	PushOverflowProtectionU32 = uint32(-0x8000 & 0xffffffff)
)
