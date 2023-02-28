package xcontext

import (
	"context"
)

type ctxKeyValue struct {
	key, value any
}

type ValueBag interface {
	AddValue(key, val any)
	value(key any) any
}

type valueBag struct {
	values []*ctxKeyValue
}

type ValueBagInjector func(ValueBag)

type valuesContext struct {
	context.Context
	ValueBag
}

func NewValueBag() ValueBag {
	return &valueBag{
		values: []*ctxKeyValue{},
	}
}

func (vb *valueBag) AddValue(key, val any) {
	vb.values = append(vb.values, &ctxKeyValue{key, val})
}

func (vb *valueBag) value(key any) any {
	for _, kv := range vb.values {
		if kv.key == key {
			return kv.value
		}
	}
	return nil
}

func NewValuesContext(ctx context.Context, valueBag ValueBag) context.Context {
	return &valuesContext{
		Context:  ctx,
		ValueBag: valueBag,
	}
}

func (v *valuesContext) Value(key interface{}) interface{} {
	val := v.value(key)
	if val != nil {
		return val
	}
	return v.Context.Value(key)
}
