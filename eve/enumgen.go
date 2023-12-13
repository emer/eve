// Code generated by "goki generate"; DO NOT EDIT.

package eve

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync/atomic"

	"goki.dev/enums"
	"goki.dev/ki/v2"
)

var _NodeTypesValues = []NodeTypes{0, 1, 2}

// NodeTypesN is the highest valid value
// for type NodeTypes, plus one.
const NodeTypesN NodeTypes = 3

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _NodeTypesNoOp() {
	var x [1]struct{}
	_ = x[BODY-(0)]
	_ = x[GROUP-(1)]
	_ = x[JOINT-(2)]
}

var _NodeTypesNameToValueMap = map[string]NodeTypes{
	`BODY`:  0,
	`body`:  0,
	`GROUP`: 1,
	`group`: 1,
	`JOINT`: 2,
	`joint`: 2,
}

var _NodeTypesDescMap = map[NodeTypes]string{
	0: `note: uppercase required to not conflict with type names`,
	1: ``,
	2: ``,
}

var _NodeTypesMap = map[NodeTypes]string{
	0: `BODY`,
	1: `GROUP`,
	2: `JOINT`,
}

// String returns the string representation
// of this NodeTypes value.
func (i NodeTypes) String() string {
	if str, ok := _NodeTypesMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the NodeTypes value from its
// string representation, and returns an
// error if the string is invalid.
func (i *NodeTypes) SetString(s string) error {
	if val, ok := _NodeTypesNameToValueMap[s]; ok {
		*i = val
		return nil
	}
	if val, ok := _NodeTypesNameToValueMap[strings.ToLower(s)]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type NodeTypes")
}

// Int64 returns the NodeTypes value as an int64.
func (i NodeTypes) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the NodeTypes value from an int64.
func (i *NodeTypes) SetInt64(in int64) {
	*i = NodeTypes(in)
}

// Desc returns the description of the NodeTypes value.
func (i NodeTypes) Desc() string {
	if str, ok := _NodeTypesDescMap[i]; ok {
		return str
	}
	return i.String()
}

// NodeTypesValues returns all possible values
// for the type NodeTypes.
func NodeTypesValues() []NodeTypes {
	return _NodeTypesValues
}

// Values returns all possible values
// for the type NodeTypes.
func (i NodeTypes) Values() []enums.Enum {
	res := make([]enums.Enum, len(_NodeTypesValues))
	for i, d := range _NodeTypesValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type NodeTypes.
func (i NodeTypes) IsValid() bool {
	_, ok := _NodeTypesMap[i]
	return ok
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i NodeTypes) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *NodeTypes) UnmarshalText(text []byte) error {
	if err := i.SetString(string(text)); err != nil {
		log.Println(err)
	}
	return nil
}

var _NodeFlagsValues = []NodeFlags{7}

// NodeFlagsN is the highest valid value
// for type NodeFlags, plus one.
const NodeFlagsN NodeFlags = 8

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _NodeFlagsNoOp() {
	var x [1]struct{}
	_ = x[Dynamic-(7)]
}

var _NodeFlagsNameToValueMap = map[string]NodeFlags{
	`Dynamic`: 7,
	`dynamic`: 7,
}

var _NodeFlagsDescMap = map[NodeFlags]string{
	7: `Dynamic means that this node can move -- if not so marked, it is a Static node. Any top-level group that is not Dynamic is immediately pruned from further consideration, so top-level groups should be separated into Dynamic and Static nodes at the start.`,
}

var _NodeFlagsMap = map[NodeFlags]string{
	7: `Dynamic`,
}

// String returns the string representation
// of this NodeFlags value.
func (i NodeFlags) String() string {
	str := ""
	for _, ie := range ki.FlagsValues() {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	for _, ie := range _NodeFlagsValues {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// BitIndexString returns the string
// representation of this NodeFlags value
// if it is a bit index value
// (typically an enum constant), and
// not an actual bit flag value.
func (i NodeFlags) BitIndexString() string {
	if str, ok := _NodeFlagsMap[i]; ok {
		return str
	}
	return ki.Flags(i).BitIndexString()
}

// SetString sets the NodeFlags value from its
// string representation, and returns an
// error if the string is invalid.
func (i *NodeFlags) SetString(s string) error {
	*i = 0
	return i.SetStringOr(s)
}

// SetStringOr sets the NodeFlags value from its
// string representation while preserving any
// bit flags already set, and returns an
// error if the string is invalid.
func (i *NodeFlags) SetStringOr(s string) error {
	flgs := strings.Split(s, "|")
	for _, flg := range flgs {
		if val, ok := _NodeFlagsNameToValueMap[flg]; ok {
			i.SetFlag(true, &val)
		} else if val, ok := _NodeFlagsNameToValueMap[strings.ToLower(flg)]; ok {
			i.SetFlag(true, &val)
		} else if flg == "" {
			continue
		} else {
			err := (*ki.Flags)(i).SetStringOr(flg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Int64 returns the NodeFlags value as an int64.
func (i NodeFlags) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the NodeFlags value from an int64.
func (i *NodeFlags) SetInt64(in int64) {
	*i = NodeFlags(in)
}

// Desc returns the description of the NodeFlags value.
func (i NodeFlags) Desc() string {
	if str, ok := _NodeFlagsDescMap[i]; ok {
		return str
	}
	return ki.Flags(i).Desc()
}

// NodeFlagsValues returns all possible values
// for the type NodeFlags.
func NodeFlagsValues() []NodeFlags {
	es := ki.FlagsValues()
	res := make([]NodeFlags, len(es))
	for i, e := range es {
		res[i] = NodeFlags(e)
	}
	res = append(res, _NodeFlagsValues...)
	return res
}

// Values returns all possible values
// for the type NodeFlags.
func (i NodeFlags) Values() []enums.Enum {
	es := ki.FlagsValues()
	les := len(es)
	res := make([]enums.Enum, les+len(_NodeFlagsValues))
	for i, d := range es {
		res[i] = d
	}
	for i, d := range _NodeFlagsValues {
		res[i+les] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type NodeFlags.
func (i NodeFlags) IsValid() bool {
	_, ok := _NodeFlagsMap[i]
	if !ok {
		return ki.Flags(i).IsValid()
	}
	return ok
}

// HasFlag returns whether these
// bit flags have the given bit flag set.
func (i NodeFlags) HasFlag(f enums.BitFlag) bool {
	return atomic.LoadInt64((*int64)(&i))&(1<<uint32(f.Int64())) != 0
}

// SetFlag sets the value of the given
// flags in these flags to the given value.
func (i *NodeFlags) SetFlag(on bool, f ...enums.BitFlag) {
	var mask int64
	for _, v := range f {
		mask |= 1 << v.Int64()
	}
	in := int64(*i)
	if on {
		in |= mask
		atomic.StoreInt64((*int64)(i), in)
	} else {
		in &^= mask
		atomic.StoreInt64((*int64)(i), in)
	}
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i NodeFlags) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *NodeFlags) UnmarshalText(text []byte) error {
	if err := i.SetString(string(text)); err != nil {
		log.Println(err)
	}
	return nil
}
