package middleware

import "reflect"

// mockFieldLevel implements validator.FieldLevel for testing purposes
type mockFieldLevel struct {
	field string
}

func (m *mockFieldLevel) Top() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Parent() reflect.Value {
	return reflect.Value{}
}

func (m *mockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.field)
}

func (m *mockFieldLevel) FieldName() string {
	return ""
}

func (m *mockFieldLevel) StructFieldName() string {
	return ""
}

func (m *mockFieldLevel) Param() string {
	return ""
}

func (m *mockFieldLevel) GetTag() string {
	return ""
}

func (m *mockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), false
}

func (m *mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func (m *mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}
