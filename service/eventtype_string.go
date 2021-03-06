// Code generated by "stringer -type=EventType"; DO NOT EDIT.

package service

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownEvent-0]
	_ = x[ConnectedEvent-1]
	_ = x[DisconnectedEvent-2]
	_ = x[MessageEvent-3]
}

const _EventType_name = "UnknownEventConnectedEventDisconnectedEventMessageEvent"

var _EventType_index = [...]uint8{0, 12, 26, 43, 55}

func (i EventType) String() string {
	if i < 0 || i >= EventType(len(_EventType_index)-1) {
		return "EventType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EventType_name[_EventType_index[i]:_EventType_index[i+1]]
}
