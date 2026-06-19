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

## Simulating unhealthy volumes / storage (development only)

By default every hostpath volume reports as healthy through the volume health
RPCs (`ControllerGetVolumeHealth`, `ControllerListVolumeHealth`,
`NodeGetVolumeHealth`, `NodeGetStorageHealth`). To exercise the unhealthy path
during development and testing, a separate `healthctl` binary writes and
removes marker files in the driver's state directory. The driver reads these
markers on every health RPC.

The `healthctl` binary is built alongside `hostpathplugin` by `make` and
operates directly on the state directory (no gRPC connection to the driver
needed). It is intended for manual use against a running driver; markers
survive driver restarts.

### Build

```
$ make
```

This produces `./bin/healthctl` and `./bin/hostpathplugin`.

### Subcommands

The `-statedir` flag (default `/csi-data-dir`) must point at the same
directory the driver was started with.

```
# Mark a volume unhealthy on both controller and node (default scope)
./bin/healthctl -statedir <stateDir> mark-volume-unhealthy <volID> \
    --status INACCESSIBLE --reason simulated --message "disk gone"

# Mark a volume unhealthy only on the controller side
./bin/healthctl -statedir <stateDir> mark-volume-unhealthy <volID> \
    --scope controller --status INACCESSIBLE --reason simulated

# Mark a volume unhealthy only on the node side
./bin/healthctl -statedir <stateDir> mark-volume-unhealthy <volID> \
    --scope node --status DATA_LOSS --reason "node-local issue"

# Clear a volume's unhealthy markers (all scopes by default)
./bin/healthctl -statedir <stateDir> mark-volume-healthy <volID>

# Clear only the controller-side marker for a volume
./bin/healthctl -statedir <stateDir> mark-volume-healthy <volID> --scope controller

# Mark node storage unhealthy (status defaults to STORAGE_DEGRADED)
./bin/healthctl -statedir <stateDir> mark-storage-unhealthy \
    --status STORAGE_UNREACHABLE --reason "disk gone"

# Clear the node storage unhealthy marker
./bin/healthctl -statedir <stateDir> mark-storage-healthy

# List all current markers
./bin/healthctl -statedir <stateDir> list
```

### Scope

Per-volume markers are scope-aware via `--scope` (default `both`):

- `both` writes `<volID>.health`, read by both `ControllerGetVolumeHealth`/
  `ControllerListVolumeHealth` and `NodeGetVolumeHealth`.
- `controller` writes `<volID>.controller.health`, read only by the
  controller-side RPCs (takes precedence over `both`).
- `node` writes `<volID>.node.health`, read only by `NodeGetVolumeHealth`
  (takes precedence over `both`).

This lets you simulate a volume that is unhealthy on one side but not the
other, matching how the CSI spec distinguishes controller-reported
`INACCESSIBLE` ("not accessible from all nodes") from node-reported
`INACCESSIBLE` ("not accessible from that node").

### Status values

Volume health statuses (`--status` for `mark-volume-unhealthy`):
`DEGRADED`, `INACCESSIBLE`, `DATA_LOSS`.

Storage health statuses (`--status` for `mark-storage-unhealthy`):
`STORAGE_UNREACHABLE`, `STORAGE_DEGRADED`.

### Behaviour notes

- Markers are stored as `<stateDir>/<volID>.health`,
  `<stateDir>/<volID>.controller.health`, `<stateDir>/<volID>.node.health`
  (per volume, per scope) and `<stateDir>/storage.health` (node storage).
  They are written atomically.
- Deleting a volume removes all of its health markers (all scopes).
- `ControllerListVolumeHealth` only returns volumes that have a
  controller-effective marker (controller-scoped or `both`) and still exist
  in the driver's state, so stale or node-only markers are skipped.
- An empty `HealthStatuses` list (no marker) means healthy per the CSI spec.

