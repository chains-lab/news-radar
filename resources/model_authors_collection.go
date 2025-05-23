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

// checks if the AuthorsCollection type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthorsCollection{}

// AuthorsCollection struct for AuthorsCollection
type AuthorsCollection struct {
	Data AuthorsCollectionData `json:"data"`
}

type _AuthorsCollection AuthorsCollection

// NewAuthorsCollection instantiates a new AuthorsCollection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthorsCollection(data AuthorsCollectionData) *AuthorsCollection {
	this := AuthorsCollection{}
	this.Data = data
	return &this
}

// NewAuthorsCollectionWithDefaults instantiates a new AuthorsCollection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthorsCollectionWithDefaults() *AuthorsCollection {
	this := AuthorsCollection{}
	return &this
}

// GetData returns the Data field value
func (o *AuthorsCollection) GetData() AuthorsCollectionData {
	if o == nil {
		var ret AuthorsCollectionData
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *AuthorsCollection) GetDataOk() (*AuthorsCollectionData, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *AuthorsCollection) SetData(v AuthorsCollectionData) {
	o.Data = v
}

func (o AuthorsCollection) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthorsCollection) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

func (o *AuthorsCollection) UnmarshalJSON(data []byte) (err error) {
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

	varAuthorsCollection := _AuthorsCollection{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAuthorsCollection)

	if err != nil {
		return err
	}

	*o = AuthorsCollection(varAuthorsCollection)

	return err
}

type NullableAuthorsCollection struct {
	value *AuthorsCollection
	isSet bool
}

func (v NullableAuthorsCollection) Get() *AuthorsCollection {
	return v.value
}

func (v *NullableAuthorsCollection) Set(val *AuthorsCollection) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthorsCollection) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthorsCollection) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthorsCollection(val *AuthorsCollection) *NullableAuthorsCollection {
	return &NullableAuthorsCollection{value: val, isSet: true}
}

func (v NullableAuthorsCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthorsCollection) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


