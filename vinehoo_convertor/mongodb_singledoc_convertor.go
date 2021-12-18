package vinehooconvertor

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func sliceToJSONString(v map[string]interface{}, slice reflect.Value, fieldname string) {
	v[fieldname] = make([]map[string]interface{}, slice.Len())

	for i := 0; i < slice.Len(); i++ {
		item := reflect.ValueOf(slice.Index(i).Interface())
		switch item.Kind() {
		case reflect.Slice:
			m_values := make(map[string]interface{})

			for j := 0; j < item.Len(); j++ {
				subitem := reflect.ValueOf(item.Index(j).Interface())
				subitem_k := reflect.ValueOf(subitem.Field(0).Interface()).String()
				subitem_v := reflect.ValueOf(subitem.Field(1).Interface())
				switch reflect.ValueOf(subitem.Field(1).Interface()).Kind() {
				case reflect.String:
					m_values[subitem_k] = subitem_v.String()
				case reflect.Float64, reflect.Float32:
					m_values[subitem_k] = subitem_v.Float()
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					m_values[subitem_k] = subitem_v.Int()
				case reflect.Uint:
					m_values[subitem_k] = subitem_v.Uint()
				case reflect.Slice:
					sliceToJSONString(m_values, subitem_v, subitem_k)
				}
				// fmt.Println("subitem:", subitem.Field(0), subitem.Field(1))
			}
			reflect.ValueOf(v[fieldname]).Index(i).Set(reflect.ValueOf(m_values))
		}
	}

}

func MongoDBSingleDocToJSONString(v interface{}) string {
	var (
		m_result = make(map[string]interface{})
		data     = make(map[string]interface{})
	)

	m_result["error_code"] = 0
	m_result["error_msg"] = ""

	t := reflect.ValueOf(v)

	switch t.Kind() {
	case reflect.Slice:
		nums := t.Len()
		for i := 0; i < nums; i++ {
			item := t.Index(i)
			k := reflect.ValueOf(item.Field(0).Interface()).String()
			// fmt.Println(k, v.Kind())
			v := reflect.ValueOf(item.Field(1).Interface())
			switch v.Kind() {
			case reflect.String:
				data[k] = v.String()
			case reflect.Uint, reflect.Uint16, reflect.Uint64, reflect.Uint8:
				data[k] = v.Uint()
			case reflect.Int, reflect.Int16, reflect.Int64, reflect.Int32:
				data[k] = v.Int()
			case reflect.Float64, reflect.Float32:
				data[k] = v.Float()
			case reflect.Slice:
				sliceToJSONString(data, v, k)
			case reflect.Array:
				if k == "_id" {
					objectid := reflect.ValueOf(reflect.ValueOf(item.Interface()).Field(1).Interface())
					ids := make([]byte, objectid.Len())
					for i := 0; i < objectid.Len(); i++ {
						ids[i] = byte(objectid.Index(i).Uint())
					}
					data[k] = fmt.Sprintf("%x", ids)
				}
			}
		}
		m_result["data"] = data

		p, _ := json.Marshal(&m_result)
		return string(p)
	default:
		return ""
	}
}
