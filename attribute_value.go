package sqldav

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/iancoleman/strcase"
	"reflect"
	"regexp"
	"strings"
)

// ErrNestedStructHasIncompatibleAttributes occurs when the nested struct has incompatible attributes.
var ErrNestedStructHasIncompatibleAttributes = errors.New("nested struct has incompatible attributes")

// AssignMapValueToReflectValue assigns the map type value to the reflect.Value
func AssignMapValueToReflectValue(rt reflect.Type, rv reflect.Value, mv map[string]interface{}) error {
	for i := 0; i < rt.NumField(); i++ {
		tf := rt.Field(i)
		vf := func() reflect.Value {
			if rv.Kind() == reflect.Pointer {
				return rv.Elem().Field(i)
			}
			return rv.Field(i)
		}()
		name := getColumnNameFromStructField(tf)
		a, ok := mv[name]
		if !ok {
			continue
		}
		err := assignInterfaceValueToReflectValue(tf.Type, vf, a)
		if err != nil {
			return err
		}
	}
	return nil
}

// ErrDocumentAttributeValueIsIncompatible occurs when an incompatible conversion to following:
//   - *types.AttributeValueMemberL
//   - *types.AttributeValueMemberM
//   - *types.AttributeValueMemberSS
//   - *types.AttributeValueMemberNS
//   - *types.AttributeValueMemberBS
var ErrDocumentAttributeValueIsIncompatible = fmt.Errorf("document-attribute-value is incompatible")

// DocumentAttributeMember represents Document Attribute Member.
type DocumentAttributeMember interface {
	*types.AttributeValueMemberL | *types.AttributeValueMemberM | *types.AttributeValueMemberSS | *types.AttributeValueMemberNS | *types.AttributeValueMemberBS
}

// ToDocumentAttributeValue converts given interface to a DocumentAttributeMember.
//
// NOTE: this function returns a typed-nil if the conversion is incompatible.
// therefore, nil-check is not guaranteed to work.
func ToDocumentAttributeValue[T DocumentAttributeMember](value interface{}) (T, error) {
	v, err := toAttibuteValue(value)
	if err != nil {
		return nil, err
	}
	if v, ok := v.(T); ok {
		return v, nil
	}
	return nil, ErrDocumentAttributeValueIsIncompatible
}

// reXORMColumnName matches column name from xorm tag
var reXORMColumnName = regexp.MustCompile(`'(.*?)'`)

// getColumnNameFromStructField returns the column name from the struct field
func getColumnNameFromStructField(sf reflect.StructField) string {
	tag := sf.Tag
	for _, value := range strings.Split(tag.Get("gorm"), ";") {
		if value == "" {
			continue
		}
		kv := strings.Split(value, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "column":
			return kv[1]
		}
	}
	if name := tag.Get("db"); name != "" {
		return name
	}
	if xt := tag.Get("xorm"); xt != "" {
		matches := reXORMColumnName.FindAllString(xt, 1)
		if len(matches) > 0 {
			return matches[1]
		}
	}
	return strcase.ToSnake(sf.Name)
}

// assignInterfaceValueToReflectValue assigns the value to the reflect.Value
func assignInterfaceValueToReflectValue(rt reflect.Type, rv reflect.Value, value interface{}) error {
	if rv.CanAddr() {
		switch sc := rv.Addr().Interface().(type) {
		case sql.Scanner:
			return sc.Scan(value)
		}
	} else {
		switch sc := rv.Interface().(type) {
		case sql.Scanner:
			return sc.Scan(value)
		}
	}
	switch rt.Kind() {
	case reflect.String:
		str, ok := value.(string)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible string and %T", value))
		}
		rv.SetString(str)
	case reflect.Int:
		f64, ok := value.(float64)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible int and %T", value))
		}
		rv.Set(reflect.ValueOf(int(f64)))
	case reflect.Bool:
		b, ok := value.(bool)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible bool and %T", value))
		}
		rv.SetBool(b)
	case reflect.Float64:
		f64, ok := value.(float64)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible float64 and %T", value))
		}
		rv.SetFloat(f64)
	case reflect.Slice:
		if rt.Elem().Kind() != reflect.Uint8 {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible []byte and %T", value))
		}
		b, ok := value.([]byte)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible []byte and %T", value))
		}
		rv.SetBytes(b)
	case reflect.Struct:
		mv, ok := value.(map[string]interface{})
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible struct and %T", value))
		}
		err := AssignMapValueToReflectValue(rt, rv, mv)
		if err != nil {
			return err
		}
	case reflect.Pointer:
		if value == nil {
			return nil
		}
		rv.Set(reflect.New(rt.Elem()))
		// NOTE: even return error, it will not be returned to the caller.
		// Only expect the attribute to be nil.
		assignInterfaceValueToReflectValue(rt.Elem(), rv.Elem(), value)
	}
	return nil
}

// toAttibuteValue converts the value to a types.AttributeValue
func toAttibuteValue(value interface{}) (types.AttributeValue, error) {
	switch value := value.(type) {
	case List:
		avs := make([]types.AttributeValue, 0, len(value))
		for _, v := range value {
			av, err := toAttibuteValue(v)
			if err != nil {
				return nil, err
			}
			avs = append(avs, av)
		}
		return &types.AttributeValueMemberL{Value: avs}, nil
	case Map:
		avm := make(map[string]types.AttributeValue)
		for k, v := range value {
			av, err := toAttibuteValue(v)
			if err != nil {
				return nil, err
			}
			avm[k] = av
		}
		return &types.AttributeValueMemberM{Value: avm}, nil
	case Set[string]:
		return &types.AttributeValueMemberSS{Value: value}, nil
	case Set[int]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Set[float64]:
		ss := make([]string, 0, len(value))
		for _, v := range value {
			ss = append(ss, fmt.Sprintf("%v", v))
		}
		return &types.AttributeValueMemberNS{Value: ss}, nil
	case Set[[]byte]:
		return &types.AttributeValueMemberBS{Value: value}, nil
	default:
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Struct:
			avm := make(map[string]types.AttributeValue)
			for i := 0; i < rv.NumField(); i++ {
				fv := rv.Field(i)
				ft := rv.Type().Field(i)
				if fv.CanInterface() {
					av, err := toAttibuteValue(fv.Interface())
					if err != nil {
						return nil, err
					}
					avm[getColumnNameFromStructField(ft)] = av
				} else if fv.CanAddr() {
					av, err := toAttibuteValue(fv.Addr().Interface())
					if err != nil {
						return nil, err
					}
					avm[getColumnNameFromStructField(ft)] = av
				}
			}
			return &types.AttributeValueMemberM{Value: avm}, nil
		case reflect.Ptr:
			if rv.IsNil() {
				return &types.AttributeValueMemberNULL{}, nil
			}
			if !rv.CanAddr() {
				return attributevalue.Marshal(value)
			}
			return toAttibuteValue(rv.Addr().Interface())
		}
		return attributevalue.Marshal(value)
	}
}
