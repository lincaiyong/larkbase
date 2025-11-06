package larkfield

import (
	"fmt"
	"strconv"
)

type NumberField struct {
	*BaseField
}

func (f *NumberField) SetIntValue(v int) {
	f.SetUnderlayValue(float64(v))
}

func (f *NumberField) SetValue(v float64) {
	f.SetUnderlayValue(v)
}

func (f *NumberField) Is(value int) *Condition {
	return conditionIs(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsNot(value int) *Condition {
	return conditionIsNot(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsGreater(value int) *Condition {
	return conditionIsGreater(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsGreaterEqual(value int) *Condition {
	return conditionIsGreaterEqual(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsLess(value int) *Condition {
	return conditionIsLess(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsLessEqual(value int) *Condition {
	return conditionIsLessEqual(f.id, f.name, strconv.Itoa(value))
}

func (f *NumberField) IsF(value float64) *Condition {
	return conditionIs(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsNotF(value float64) *Condition {
	return conditionIsNot(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsGreaterF(value float64) *Condition {
	return conditionIsGreater(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsGreaterEqualF(value float64) *Condition {
	return conditionIsGreaterEqual(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsLessF(value float64) *Condition {
	return conditionIsLess(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsLessEqualF(value float64) *Condition {
	return conditionIsLessEqual(f.id, f.name, fmt.Sprintf("%g", value))
}

func (f *NumberField) IsEmpty() *Condition {
	return conditionIsEmpty(f.id, f.name)
}

func (f *NumberField) IsNotEmpty() *Condition {
	return conditionIsNotEmpty(f.id, f.name)
}
