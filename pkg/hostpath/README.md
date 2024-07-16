# CSI Hostpath driver

## Usage:

### Build hostpathplugin
```
$ make
```

### Start Hostpath driver
```
$ sudo ./bin/hostpathplugin --endpoint tcp://127.0.0.1:10000 --nodeid CSINode -v=5
```

### Test using csc
Get ```csc``` tool from https://github.com/rexray/gocsi/tree/master/csc

#### Get plugin info
```
$ csc identity plugin-info --endpoint tcp://127.0.0.1:10000
"csi-hostpath"  "0.1.0"
```

#### Create a block volume
```
$ csc controller new --endpoint tcp://127.0.0.1:10000 --cap 1,block --req-bytes 1048576 --lim-bytes 1048576 CSIVolumeName
CSIVolumeID
```

#### Create mounted volume
```
$ csc controller new --endpoint tcp://127.0.0.1:10000 --cap MULTI_NODE_MULTI_WRITER,mount,xfs,uid=500,gid=500 CSIVolumeName
CSIVolumeID
```

#### List volumes
```
csc controller list-volumes --endpoint tcp://127.0.0.1:10000
CSIVolumeID  0
CSIVolumeID  0
```

#### Delete a volume
```
$ csc controller del --endpoint tcp://127.0.0.1:10000 CSIVolumeID
CSIVolumeID
```

#### Validate volume capabilities
```
$ csc controller validate-volume-capabilities --endpoint tcp://127.0.0.1:10000 --cap 1,block CSIVolumeID
CSIVolumeID  true
```

#### NodePublish a volume
```
$ csc node publish --endpoint tcp://127.0.0.1:10000 --cap 1,block --target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### NodeUnpublish a volume
```
$ csc node unpublish --endpoint tcp://127.0.0.1:10000 --target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### Get NodeInfo
```
$ csc node get-info --endpoint tcp://127.0.0.1:10000
CSINode
```

### Create snapshot
```
$ csc controller create-snapshot --endpoint tcp://127.0.0.1:10000 --params ignoreFailedRead=true --source-volume CSIVolumeID CSISnapshotName
CSISnapshotID
```

### Delete snapshot
```
csc controller delete-snapshot --endpoint tcp://127.0.0.1:10000 CSISnapshotID
```

### List snapshots
```
csc controller list-snapshots --endpoint tcp://127.0.0.1:10000
```
