// Code generated by "stringer -type=FieldName -linecomment"; DO NOT EDIT.

package integrationvolume

import "strconv"

const _FieldName_name = "timekyber_swap_volumenon_kyber_swap_volume"

var _FieldName_index = [...]uint8{0, 4, 21, 42}

func (i FieldName) String() string {
	if i < 0 || i >= FieldName(len(_FieldName_index)-1) {
		return "FieldName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FieldName_name[_FieldName_index[i]:_FieldName_index[i+1]]
}
