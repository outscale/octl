/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/package builder

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/outscale/gli/pkg/openapi"
	"github.com/outscale/gli/pkg/options"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type Builder struct {
	spec *openapi.Spec
}

func NewBuilder(spec *openapi.Spec) *Builder {
	return &Builder{spec: spec}
}

func (b *Builder) FromStruct(cmd *cobra.Command, arg reflect.Type, prefix string) {
	typeName := arg.Name()
	fs := cmd.Flags()
	for i := range arg.NumField() {
		f := arg.Field(i)
		t := f.Type
		ot := t
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		help := b.spec.SummaryForAttribute(typeName, f.Name)
		switch t.Kind() {
		case reflect.Bool:
			fs.Bool(prefix+f.Name, false, help)
		case reflect.String:
			fs.String(prefix+f.Name, "", help)
			if t.Implements(reflect.TypeFor[enum]()) {
				values := reflect.New(t).Interface().(enum).Values()
				_ = cmd.RegisterFlagCompletionFunc(prefix+f.Name, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
					return lo.Map(values, func(v string, _ int) cobra.Completion { return cobra.Completion(v) }), cobra.ShellCompDirectiveDefault
				})
			}
		case reflect.Int:
			fs.Int(prefix+f.Name, 0, help)
		case reflect.Slice:
			switch t.Elem().Kind() {
			case reflect.Bool:
				fs.BoolSlice(prefix+f.Name, nil, help)
			case reflect.String:
				fs.StringSlice(prefix+f.Name, nil, help)
				if t.Elem().Implements(reflect.TypeFor[enum]()) {
					values := reflect.New(t.Elem()).Interface().(enum).Values()
					_ = cmd.RegisterFlagCompletionFunc(prefix+f.Name, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
						return lo.Map(values, func(v string, _ int) cobra.Completion { return cobra.Completion(v) }), cobra.ShellCompDirectiveDefault
					})
				}
			case reflect.Int:
				fs.IntSlice(prefix+f.Name, nil, help)
			case reflect.Struct:
				if t.Elem().Implements(reflect.TypeFor[json.Marshaler]()) {
					fs.StringSlice(prefix+f.Name, nil, help)
				} else {
					for i := range options.NumEntriesInSlices {
						b.FromStruct(cmd, t.Elem(), prefix+f.Name+"."+strconv.Itoa(i)+".")
					}
				}
			}
		case reflect.Struct:
			if ot.Implements(reflect.TypeFor[json.Marshaler]()) {
				fs.String(prefix+f.Name, "", help)
			} else {
				b.FromStruct(cmd, t, prefix+f.Name+".")
			}
		}
	}
}
