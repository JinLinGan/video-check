package obs

// VirtualMicInfo 麦克风信息
type VirtualMicInfo struct {

	// Name 名字
	Name string
	// ID ID
	ID string
}

// CABLE Input (VB-Audio Virtual Cable)
// {0.0.0.00000000}.{ce4e42b4-c623-41d1-938f-83652535c9d0}

func NewVirtualMicInfo(name, id string) *VirtualMicInfo {
	return &VirtualMicInfo{
		Name: name,
		ID:   id,
	}
}
