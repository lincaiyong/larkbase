package larkbase

import (
	"context"
)

func CreateAll(ctx context.Context, url string, data []map[string]string, tosFields string) error {
	conn, err := ConnectAny(ctx, url)
	if err != nil {
		return err
	}
	var fields []string
	var records []*AnyRecord
	for i, row := range data {
		var r AnyRecord
		r.Data = make(map[string]string)
		for k, v := range row {
			if i == 0 {
				fields = append(fields, k)
			}
			r.Data[k] = v
			r.Update(k, v)
		}
		records = append(records, &r)
	}
	_, err = conn.CreateAllAny(fields, records)
	if err != nil {
		return err
	}
	return nil
}
