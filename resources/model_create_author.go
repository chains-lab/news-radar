/*
REST API

REST API

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package resources

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the CreateAuthor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateAuthor{}

// CreateAuthor struct for CreateAuthor
type CreateAuthor struct {
	Data CreateAuthorData `json:"data"`
}

type _CreateAuthor CreateAuthor

// NewCreateAuthor instantiates a new CreateAuthor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateAuthor(data CreateAuthorData) *CreateAuthor {
	this := CreateAuthor{}
	this.Data = data
	return &this
}

// NewCreateAuthorWithDefaults instantiates a new CreateAuthor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateAuthorWithDefaults() *CreateAuthor {
	this := CreateAuthor{}
	return &this
}

// GetData returns the Data field value
func (o *CreateAuthor) GetData() CreateAuthorData {
	if o == nil {
		var ret CreateAuthorData
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *CreateAuthor) GetDataOk() (*CreateAuthorData, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *CreateAuthor) SetData(v CreateAuthorData) {
	o.Data = v
}

func (o CreateAuthor) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateAuthor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

func (o *CreateAuthor) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"data",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varCreateAuthor := _CreateAuthor{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCreateAuthor)

	if err != nil {
		return err
	}

	*o = CreateAuthor(varCreateAuthor)

	return err
}

type NullableCreateAuthor struct {
	value *CreateAuthor
	isSet bool
}

func (v NullableCreateAuthor) Get() *CreateAuthor {
	return v.value
}

func (v *NullableCreateAuthor) Set(val *CreateAuthor) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateAuthor) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateAuthor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateAuthor(val *CreateAuthor) *NullableCreateAuthor {
	return &NullableCreateAuthor{value: val, isSet: true}
}

func (v NullableCreateAuthor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateAuthor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


