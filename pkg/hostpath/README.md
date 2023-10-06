# CSI Hostpath driver

## Usage:

### Build hostpathplugin
```
$ make
```

### Set endpoint
```
$ endpoint=unix:///tmp/csi.sock # unix (default)
$ #endpoint=tcp://127.0.0.1:10000 # tcp
```

### Start Hostpath driver
```
$ sudo ./bin/hostpathplugin --endpoint $endpoint --nodeid CSINode -v=5
```

### Test using csc
Get ```csc``` tool from https://github.com/rexray/gocsi/tree/master/csc

#### Get plugin info
```
$ sudo csc identity plugin-info --endpoint $endpoint
"hostpath.csi.k8s.io"  "v1.x.x-xx-xxx"
```

#### Create a volume
```
$ sudo csc controller new --endpoint $endpoint --cap 1,block --req-bytes 10240000 CSIVolumeName
CSIVolumeID
```

#### Delete a volume
```
$ sudo csc controller del --endpoint $endpoint CSIVolumeID
CSIVolumeID
```

#### Validate volume capabilities
```
$ sudo csc controller validate-volume-capabilities --endpoint $endpoint --cap 1,block CSIVolumeID
CSIVolumeID  volume_capabilities:<block:<> access_mode:<mode:SINGLE_NODE_WRITER > >
```

#### NodeStage a volume
```
$ sudo csc node stage --endpoint $endpoint --cap 1,block --staging-target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### NodeUnstage a volume
```
$ sudo csc node unstage --endpoint $endpoint --cap 1,block --staging-target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### NodePublish a volume
```
$ sudo csc node publish --endpoint $endpoint --cap 1,block --staging-target-path /mnt/hostpath --target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### NodeUnpublish a volume
```
$ sudo csc node unpublish --endpoint $endpoint --target-path /mnt/hostpath CSIVolumeID
CSIVolumeID
```

#### Get NodeInfo
```
$ sudo csc node get-info --endpoint $endpoint
CSINode
```
