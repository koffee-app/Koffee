package formatter

// u32ToInterface copies items of data to a new array, useful for conversion between types!
func u32ToInterface(data []uint32) []interface{} {
	newData := make([]interface{}, len(data))
	for i := range data {
		newData[i] = data[i]
	}
	return newData
}
