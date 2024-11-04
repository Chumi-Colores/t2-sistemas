package global

import "sync"

var Buffer_size int
var Food_channels []chan Food
var Command_channels []chan Command_data
var Final_snapshots = make(map[int]map[int]Snapshot_to_write)
var Garlic_owner int

const FileName = "snapshot_"

var Wg sync.WaitGroup
var Mu sync.Mutex
