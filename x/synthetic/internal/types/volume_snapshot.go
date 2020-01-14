package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VolumeSnapshot struct {
	LongVolume  sdk.Int `json:"long_volume" yaml:"long_volume"`
	ShortVolume sdk.Int `json:"short_volume" yaml:"short_volume"`
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

// TODO: Make fields private?
type VolumeSnapshots struct {
	SnapshotLimit      int              `json:"snapshot_limit" yaml:"snapshot_limit"`
	Coefficients       []sdk.Int        `json:"coefficients" yaml:"coefficients"`
	DefaultCoefficient sdk.Int          `json:"default_coefficient"yaml:"default_coefficient"`
	Snapshots          []VolumeSnapshot `json:"snapshots" yaml:"snapshots"`
}

// Creates new VolumeSnapshot storage
// A maximum snapshot number equals to a length of a Coefficients array
func NewVolumeSnapshots(snapshotLimit int, coefficients []sdk.Int) VolumeSnapshots {
	if snapshotLimit == 0 {
		panic("VolumeSnapshots cannot be empty")
	}

	return VolumeSnapshots{
		snapshotLimit,
		coefficients,
		sdk.OneInt(), // we assume that coefficients can only increase snapshots' volumes. Thus using OneInt as a default coefficient is the same as saving original value of a snapshot
		nil,
	}
}

func VolumeSnapshotsWithCustomCoeff(maxSnapshotNumber int, coefficients []sdk.Int, defaultCoeff sdk.Int) VolumeSnapshots {
	if maxSnapshotNumber == 0 {
		panic("VolumeSnapshots cannot be empty")
	}

	return VolumeSnapshots{
		maxSnapshotNumber,
		coefficients,
		defaultCoeff,
		nil,
	}
}

func (v *VolumeSnapshots) AddSnapshotValues(LongVolume, ShortVolume sdk.Int) {
	v.AddSnapshot(VolumeSnapshot{LongVolume, ShortVolume})
}

func (v *VolumeSnapshots) AddSnapshot(snapshot VolumeSnapshot) {
	if len(v.Snapshots) == v.SnapshotLimit {
		v.Snapshots = v.Snapshots[1:]
	}
	v.Snapshots = append(v.Snapshots, snapshot)
}

// returns a copy of an object with a new snapshot appended
func (v VolumeSnapshots) AppendSnapshot(snapshot VolumeSnapshot) VolumeSnapshots {
	if len(v.Snapshots) == v.SnapshotLimit {
		v.Snapshots = v.Snapshots[1:]
	}
	v.Snapshots = append(v.Snapshots, snapshot)
	return v
}

func (v *VolumeSnapshots) GetWeightedVolumes() *VolumeSnapshot {
	var weightedVolumes = NewEmptyVolumeSnapshot()
	for i := 0; i < len(v.Snapshots); i++ {
		coefficient := v.newOrDefaultCoefficient(i)

		weightedLongs := v.Snapshots[i].LongVolume.Mul(coefficient)
		weightedShorts := v.Snapshots[i].ShortVolume.Mul(coefficient)

		weightedVolumes.AddVolumes(weightedLongs, weightedShorts)
	}
	return &weightedVolumes
}

func (v *VolumeSnapshots) newOrDefaultCoefficient(index int) sdk.Int {
	if len(v.Coefficients) > index {
		return v.Coefficients[index]
	}

	return v.DefaultCoefficient
}
