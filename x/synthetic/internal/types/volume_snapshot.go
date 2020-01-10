package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VolumeSnapshot struct {
	LongVolume  sdk.Uint `json:"long_volume"`
	ShortVolume sdk.Uint `json:"short_volume"`
}

func NewVolumeSnapshot() VolumeSnapshot {
	return VolumeSnapshot{
		sdk.ZeroUint(),
		sdk.ZeroUint(),
	}
}

// we assume that you might want to check
func (v *VolumeSnapshot) AddAndVerify(snap *VolumeSnapshot) {
	v.ShortVolume = v.ShortVolume.Add(snap.ShortVolume)
	v.LongVolume = v.LongVolume.Add(snap.LongVolume)
}

func (v *VolumeSnapshot) Add(snap *VolumeSnapshot) {
	v.ShortVolume = v.ShortVolume.Add(snap.ShortVolume)
	v.LongVolume = v.LongVolume.Add(snap.LongVolume)
}

func (v *VolumeSnapshot) AddVolumes(LongVolume, ShortVolume sdk.Uint) {
	v.ShortVolume = v.ShortVolume.Add(ShortVolume)
	v.LongVolume = v.LongVolume.Add(LongVolume)
}

type VolumeSnapshots struct {
	maxSnapshotNumber int
	coefficients      []sdk.Uint
	snapshots         []VolumeSnapshot
}

// Creates new VolumeSnapshot storage
// A maximum snapshot number equals to a length of a coefficients array
func NewVolumeSnapshots(coefficients []sdk.Uint) VolumeSnapshots {
	maxSnapshotNumber := len(coefficients)
	if maxSnapshotNumber == 0 {
		panic("coefficients cannot be empty")
	}

	v := VolumeSnapshots{maxSnapshotNumber, coefficients, nil}
	return v
}

func (v *VolumeSnapshots) AddSnapshot(LongVolume, ShortVolume sdk.Uint) {
	v.snapshots = append(v.snapshots, VolumeSnapshot{LongVolume, ShortVolume})
}

func (v *VolumeSnapshots) GetWeightedVolumes() *VolumeSnapshot {
	var weightedVolumes = NewVolumeSnapshot()
	for i := 0; i < len(v.snapshots); i++ {
		coefficient := v.coefficients[i]

		weightedLongs := v.snapshots[i].LongVolume.Mul(coefficient)
		weightedShorts := v.snapshots[i].ShortVolume.Mul(coefficient)

		weightedVolumes.AddVolumes(weightedLongs, weightedShorts)
	}
	return &weightedVolumes
}
