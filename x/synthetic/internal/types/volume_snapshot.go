package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VolumeSnapshot struct {
	LongVolume  sdk.Int `json:"long_volume"`
	ShortVolume sdk.Int `json:"short_volume"`
}

func NewVolumeSnapshot(long, short sdk.Int) VolumeSnapshot {
	return VolumeSnapshot{
		long,
		short,
	}
}

func NewEmptyVolumeSnapshot() VolumeSnapshot {
	return VolumeSnapshot{
		sdk.ZeroInt(),
		sdk.ZeroInt(),
	}
}

func (v *VolumeSnapshot) Add(snap *VolumeSnapshot) {
	v.ShortVolume = v.ShortVolume.Add(snap.ShortVolume)
	v.LongVolume = v.LongVolume.Add(snap.LongVolume)
}

func (v *VolumeSnapshot) AddVolumes(LongVolume, ShortVolume sdk.Int) {
	v.ShortVolume = v.ShortVolume.Add(ShortVolume)
	v.LongVolume = v.LongVolume.Add(LongVolume)
}

type VolumeSnapshots struct {
	maxSnapshotNumber int
	coefficients      []sdk.Int
	snapshots         []VolumeSnapshot
}

// Creates new VolumeSnapshot storage
// A maximum snapshot number equals to a length of a coefficients array
func NewVolumeSnapshots(maxSnapshotNumber int, coefficients []sdk.Int) VolumeSnapshots {

	if maxSnapshotNumber == 0 {
		panic("snapshots cannot be empty")
	}

	v := VolumeSnapshots{maxSnapshotNumber, coefficients, nil}
	return v
}

func (v *VolumeSnapshots) AddSnapshotValues(LongVolume, ShortVolume sdk.Int) {
	v.AddSnapshot(VolumeSnapshot{LongVolume, ShortVolume})
}


func (v *VolumeSnapshots) AddSnapshot(snapshot VolumeSnapshot) {
	if len(v.snapshots) == v.maxSnapshotNumber {
		v.snapshots = v.snapshots[1:]
	}
	v.snapshots = append(v.snapshots, snapshot)
}

func (v *VolumeSnapshots) GetWeightedVolumes() *VolumeSnapshot {
	var weightedVolumes = NewEmptyVolumeSnapshot()
	for i := 0; i < len(v.snapshots); i++ {
		coefficient := v.coefficients[i]

		weightedLongs := v.snapshots[i].LongVolume.Mul(coefficient)
		weightedShorts := v.snapshots[i].ShortVolume.Mul(coefficient)

		weightedVolumes.AddVolumes(weightedLongs, weightedShorts)
	}
	return &weightedVolumes
}
