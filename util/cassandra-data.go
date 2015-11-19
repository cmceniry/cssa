package util

type SnapFile struct {
	Filename	string
	Filepath	string
	Inode		uint64
	Size		int64
	Keyspace	string
	ColumnFamily	string
}

type SnapManifest struct {
	// TODO: add some date information here
	SnapshotName	string
	Files		[]*SnapFile
}
