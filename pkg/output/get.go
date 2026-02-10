/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package output

import (
	"fmt"
	"reflect"

	"github.com/outscale/gli/pkg/config"
)

func GetRow(v any, cols config.Columns) ([]string, error) {
	row := make([]string, len(cols))
	for i, c := range cols {
		val, err := c.Get(v)
		if err != nil {
			return nil, err
		}
		switch {
		case val == nil || reflect.ValueOf(val).IsZero():
			continue
		case reflect.TypeOf(val).Kind() == reflect.Ptr:
			val = reflect.ValueOf(val).Elem().Interface()
		}
		row[i] = fmt.Sprint(val)
	}
	return row, nil
}
