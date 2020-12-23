package collector

import (
    "github.com/panjf2000/ants/v2"
    "sync"
    
    "github.com/ul-gaul/go-basestation/data/packet"
)

type DataCallback func(packet.PacketList)

type IDataCollector interface {
    // Packets returns all the data added since the application started or since the last Clear.
    Packets() packet.PacketList
    
    // Clear clears all the data collected
    Clear()
    
    // AddPackets appends the packets to the list
    AddPackets(packets ...packet.RocketPacket)
    
    // AddCallback adds a callback which is called when new packets are received/read.
    AddCallback(cb DataCallback) uint
    
    // RemoveCallback removes the specified callback
    RemoveCallback(id uint)
}

var _ IDataCollector = (*dataCollector)(nil)

type dataCollector struct {
    mutCallbacks, mutData sync.RWMutex
    data, newData         packet.PacketList
    chDataChanged         chan struct{}
    callbacks             map[uint]DataCallback
    lastID                uint
}

func New() (IDataCollector, error) {
    dc := &dataCollector{
        chDataChanged: make(chan struct{}, 1),
        callbacks:     make(map[uint]DataCallback),
    }
    if err := ants.Submit(dc.run); err != nil {
        return nil, err
    }
    return dc, nil
}

func (dc *dataCollector) run() {
    var chunk []packet.RocketPacket
    var wg sync.WaitGroup
    
    for range dc.chDataChanged {
        dc.mutData.Lock()
        chunk = dc.newData
        dc.newData = nil
        dc.mutData.Unlock()
        
        dc.mutCallbacks.RLock()
        wg.Add(len(dc.callbacks))
        for _, callback := range dc.callbacks {
            go func(cb DataCallback) {
                cb(chunk)
                wg.Done()
            }(callback)
        }
        dc.mutCallbacks.RUnlock()
        wg.Wait()
    }
}

func (dc *dataCollector) AddPackets(packets ...packet.RocketPacket) {
    dc.mutData.Lock()
    dc.newData = append(dc.newData, packets...)
    dc.data = append(dc.data, packets...)
    dc.mutData.Unlock()
    
    select {
    case dc.chDataChanged <- struct{}{}:
    default:
    }
}

func (dc *dataCollector) Packets() packet.PacketList {
    dc.mutData.RLock()
    defer dc.mutData.RUnlock()
    return dc.data
}

func (dc *dataCollector) Clear() {
    dc.mutData.Lock()
    defer dc.mutData.Unlock()
    dc.data = nil
    dc.newData = nil
}

func (dc *dataCollector) AddCallback(cb DataCallback) uint {
    dc.mutCallbacks.Lock()
    defer dc.mutCallbacks.Unlock()
    dc.lastID++
    dc.callbacks[dc.lastID] = cb
    return dc.lastID
}

func (dc *dataCollector) RemoveCallback(cb uint) {
    dc.mutCallbacks.Lock()
    defer dc.mutCallbacks.Unlock()
    if _, ok := dc.callbacks[cb]; ok {
        delete(dc.callbacks, cb)
    }
}
