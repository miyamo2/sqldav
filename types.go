package sqldav

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"math"
	"reflect"
	"slices"
)

var (
	ErrValueIsIncompatibleOfStringSlice  = errors.New("value is incompatible of string slice")
	ErrValueIsIncompatibleOfIntSlice     = errors.New("value is incompatible of int slice")
	ErrValueIsIncompatibleOfFloat64Slice = errors.New("value is incompatible of float64 slice")
	ErrValueIsIncompatibleOfBinarySlice  = errors.New("value is incompatible of []byte slice")
	ErrCollectionAlreadyContainsItem     = errors.New("collection already contains item")
	ErrFailedToCast                      = errors.New("failed to cast")
)

// SetSupportable are the types that support the Set
type SetSupportable interface {
	string | []byte | int | float64
}

// compatibility check
var (
	_ driver.Valuer = (*Set[string])(nil)
	_ sql.Scanner   = (*Set[string])(nil)
)

// Set is a DynamoDB set type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Set[T SetSupportable] []T

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (s *Set[T]) Scan(value interface{}) error {
	if len(*s) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	if value == nil {
		*s = nil
		return nil
	}
	switch (interface{})(s).(type) {
	case *Set[int]:
		return scanAsIntSet((interface{})(s).(*Set[int]), value)
	case *Set[float64]:
		return scanAsFloat64Set((interface{})(s).(*Set[float64]), value)
	case *Set[string]:
		return scanAsStringSet((interface{})(s).(*Set[string]), value)
	case *Set[[]byte]:
		return scanAsBinarySet((interface{})(s).(*Set[[]byte]), value)
	}
	return nil
}

func (s Set[T]) Value() (v driver.Value, err error) {
	switch s := (interface{})(s).(type) {
	case Set[int]:
		v, err = numericSetToAttributeValue(s)
	case Set[float64]:
		v, err = numericSetToAttributeValue(s)
	case Set[string]:
		v, err = stringSetToAttributeValue(s)
	case Set[[]byte]:
		v, err = binarySetToAttributeValue(s)
	}
	return
}

// GormDataType returns the data type for Gorm.
func (s *Set[T]) GormDataType() string {
	var t T
	switch (interface{})(t).(type) {
	case string:
		return "SS"
	case int:
		return "NS"
	case float64:
		return "NS"
	case []byte:
		return "BS"
	}
	return "SS"
}

func numericSetToAttributeValue[T Set[int] | Set[float64]](s T) (*types.AttributeValueMemberNS, error) {
	return ToDocumentAttributeValue[*types.AttributeValueMemberNS](s)
}

func stringSetToAttributeValue(s Set[string]) (*types.AttributeValueMemberSS, error) {
	return ToDocumentAttributeValue[*types.AttributeValueMemberSS](s)
}

func binarySetToAttributeValue(s Set[[]byte]) (*types.AttributeValueMemberBS, error) {
	return ToDocumentAttributeValue[*types.AttributeValueMemberBS](s)
}

// scanAsIntSet scans the value as Set[int]
func scanAsIntSet(s *Set[int], value interface{}) error {
	sv, ok := value.([]float64)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfIntSlice
	}
	for _, v := range sv {
		if math.Floor(v) != v {
			*s = nil
			return ErrValueIsIncompatibleOfIntSlice
		}
		*s = append(*s, int(v))
	}
	return nil
}

// scanAsFloat64Set scans the value as Set[float64]
func scanAsFloat64Set(s *Set[float64], value interface{}) error {
	sv, ok := value.([]float64)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfFloat64Slice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

// scanAsStringSet scans the value as Set[string]
func scanAsStringSet(s *Set[string], value interface{}) error {
	sv, ok := value.([]string)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfStringSlice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

// scanAsBinarySet scans the value as Set[[]byte]
func scanAsBinarySet(s *Set[[]byte], value interface{}) error {
	sv, ok := value.([][]byte)
	if !ok {
		*s = nil
		return ErrValueIsIncompatibleOfBinarySlice
	}
	for _, v := range sv {
		*s = append(*s, v)
	}
	return nil
}

func isCompatibleWithSet[T SetSupportable](value interface{}) (compatible bool) {
	var t T
	switch (interface{})(t).(type) {
	case string:
		compatible = isStringSetCompatible(value)
	case int:
		compatible = isIntSetCompatible(value)
	case float64:
		compatible = isFloat64SetCompatible(value)
	case []byte:
		compatible = isBinarySetCompatible(value)
	}
	return
}

func isIntSetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]int); ok {
		compatible = true
		return
	}
	if value, ok := value.([]float64); ok {
		compatible = true
		for _, v := range value {
			if math.Floor(v) == v {
				compatible = true
				continue
			}
			compatible = false
			return
		}
	}
	return
}

func isStringSetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]string); ok {
		compatible = true
	}
	return
}

func isFloat64SetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([]float64); ok {
		compatible = true
	}
	return
}

func isBinarySetCompatible(value interface{}) (compatible bool) {
	if _, ok := value.([][]byte); ok {
		compatible = true
	}
	return
}

func newSet[T SetSupportable]() Set[T] {
	return Set[T]{}
}

// compatibility check
var (
	_ driver.Valuer = (*List)(nil)
	_ sql.Scanner   = (*List)(nil)
)

// List is a DynamoDB list type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type List []interface{}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (l *List) Scan(value interface{}) error {
	if len(*l) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	sv, ok := value.([]interface{})
	if !ok {
		return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", l, value))
	}
	*l = sv
	return resolveCollectionsNestedInList(l)
}

// Value implements the [driver.Valuer] interface.
//
// [driver.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (l List) Value() (v driver.Value, err error) {
	if err = resolveCollectionsNestedInList(&l); err != nil {
		return
	}
	v, err = ToDocumentAttributeValue[*types.AttributeValueMemberL](l)
	return
}

// GormDataType returns the data type for Gorm.
func (l *List) GormDataType() string {
	return "L"
}

// resolveCollectionsNestedInList resolves nested collection type attribute.
func resolveCollectionsNestedInList(l *List) error {
	for i, v := range *l {
		if v, ok := v.(map[string]interface{}); ok {
			m := Map{}
			err := m.Scan(v)
			if err != nil {
				*l = nil
				return err
			}
			(*l)[i] = m
			continue
		}
		if isCompatibleWithSet[int](v) {
			s := newSet[int]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[float64](v) {
			s := newSet[float64]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[string](v) {
			s := newSet[string]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if isCompatibleWithSet[[]byte](v) {
			s := newSet[[]byte]()
			if err := s.Scan(v); err == nil {
				(*l)[i] = s
				continue
			}
		}
		if v, ok := v.([]interface{}); ok {
			il := List{}
			err := il.Scan(v)
			if err != nil {
				*l = nil
				return err
			}
			(*l)[i] = il
		}
	}
	return nil
}

// compatibility check
var (
	_ driver.Valuer = (*Map)(nil)
	_ sql.Scanner   = (*Map)(nil)
)

// Map is a DynamoDB map type.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type Map map[string]interface{}

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (m *Map) Scan(value interface{}) error {
	if len(*m) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	mv, ok := value.(map[string]interface{})
	if !ok {
		*m = nil
		return ErrFailedToCast
	}
	*m = mv
	return resolveCollectionsNestedInMap(m)
}

// Value implements the [driver.Valuer] interface.
//
// [driver.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (m Map) Value() (v driver.Value, err error) {
	if err = resolveCollectionsNestedInMap(&m); err != nil {
		return
	}
	v, err = ToDocumentAttributeValue[*types.AttributeValueMemberM](m)
	return
}

// GormDataType returns the data type for Gorm.
func (m Map) GormDataType() string {
	return "M"
}

// resolveCollectionsNestedInMap resolves nested document type attribute.
func resolveCollectionsNestedInMap(m *Map) error {
	for k, v := range *m {
		if v, ok := v.(map[string]interface{}); ok {
			im := Map{}
			err := im.Scan(v)
			if err != nil {
				*m = nil
				return err
			}
			(*m)[k] = im
			continue
		}
		if isCompatibleWithSet[int](v) {
			s := newSet[int]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[float64](v) {
			s := newSet[float64]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[string](v) {
			s := newSet[string]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if isCompatibleWithSet[[]byte](v) {
			s := newSet[[]byte]()
			if err := s.Scan(v); err == nil {
				(*m)[k] = s
				continue
			}
		}
		if v, ok := v.([]interface{}); ok {
			l := List{}
			err := l.Scan(v)
			if err != nil {
				*m = nil
				return err
			}
			(*m)[k] = l
		}
	}
	return nil
}

// compatibility check
var (
	_ driver.Valuer = (*TypedList[interface{}])(nil)
	_ sql.Scanner   = (*TypedList[interface{}])(nil)
)

// TypedList is a DynamoDB list type with type specification.
//
// See: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
type TypedList[T any] []T

// Scan implements the [sql.Scanner#Scan]
//
// [sql.Scanner#Scan]: https://golang.org/pkg/database/sql/#Scanner
func (l *TypedList[T]) Scan(value interface{}) error {
	if len(*l) != 0 {
		return ErrCollectionAlreadyContainsItem
	}
	sv, ok := value.([]interface{})
	if !ok {
		return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", l, value))
	}
	*l = slices.Grow(*l, len(sv))
	for _, v := range sv {
		mv, ok := v.(map[string]interface{})
		if !ok {
			var t T
			return errors.Join(ErrFailedToCast, fmt.Errorf("incompatible %T and %T", t, v))
		}
		dest := new(T)
		rv := reflect.ValueOf(dest)
		rt := reflect.TypeOf(*dest)
		err := AssignMapValueToReflectValue(rt, rv, mv)
		if err != nil {
			return err
		}
		*l = append(*l, *dest)
	}
	return nil
}

// Value implements the [driver.Valuer] interface.
//
// [driver.Valuer]: https://pkg.go.dev/gorm.io/gorm#Valuer
func (l TypedList[T]) Value() (v driver.Value, err error) {
	avl := &types.AttributeValueMemberL{Value: make([]types.AttributeValue, 0, len(l))}
	for _, v := range l {
		av, err := ToDocumentAttributeValue[*types.AttributeValueMemberM](v)
		if err != nil {
			return nil, err
		}
		avl.Value = append(avl.Value, av)
	}
	v = avl
	return
}

// GormDataType returns the data type for Gorm.
func (l *TypedList[T]) GormDataType() string {
	return "L"
}
