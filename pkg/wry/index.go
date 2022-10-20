package wry

import (
	"encoding/binary"
)

func (db *IPDB[uint32]) SearchIndexV4(ip uint32) uint32 {
	ipLen := db.IPLen
	entryLen := uint32(db.OffLen + db.IPLen)

	buf := make([]byte, entryLen)
	l, r, mid, ipc := db.IdxStart, db.IdxEnd, uint32(0), uint32(0)

	for {
		mid = (r-l)/entryLen/2*entryLen + l
		buf = db.Data[mid : mid+entryLen]
		ipc = uint32(binary.LittleEndian.Uint32(buf[:ipLen]))

		if r-l == entryLen {
			if ip >= uint32(binary.LittleEndian.Uint32(db.Data[r:r+uint32(ipLen)])) {
				buf = db.Data[r : r+entryLen]
			}
			return uint32(Bytes3ToUint32(buf[ipLen:entryLen]))
		}

		if ipc > ip {
			r = mid
		} else if ipc < ip {
			l = mid
		} else if ipc == ip {
			return uint32(Bytes3ToUint32(buf[ipLen:entryLen]))
		}
	}
}

func (db *IPDB[uint64]) SearchIndexV6(ip uint64) uint32 {
	ipLen := db.IPLen
	entryLen := uint64(db.OffLen + db.IPLen)

	buf := make([]byte, entryLen)
	l, r, mid, ipc := db.IdxStart, db.IdxEnd, uint64(0), uint64(0)

	for {
		mid = (r-l)/entryLen/2*entryLen + l
		buf = db.Data[mid : mid+entryLen]
		ipc = uint64(binary.LittleEndian.Uint64(buf[:ipLen]))

		if r-l == entryLen {
			if ip >= uint64(binary.LittleEndian.Uint64(db.Data[r:r+uint64(ipLen)])) {
				buf = db.Data[r : r+entryLen]
			}
			return Bytes3ToUint32(buf[ipLen:entryLen])
		}

		if ipc > ip {
			r = mid
		} else if ipc < ip {
			l = mid
		} else if ipc == ip {
			return Bytes3ToUint32(buf[ipLen:entryLen])
		}
	}
}
