// Code generated by ogen, DO NOT EDIT.

package api

// setDefaults set default value of fields.
func (s *ExecutionRequest) setDefaults() {
	{
		val := int(5)
		s.MaxRetries.SetTo(val)
	}
	{
		val := int64(0)
		s.Timeout.SetTo(val)
	}
}