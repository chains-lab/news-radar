/*
REST API

REST API

API version: 0.0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package resources

import (
	"encoding/json"
)

// checks if the UpdateArticleContentDataAttributes type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateArticleContentDataAttributes{}

// UpdateArticleContentDataAttributes struct for UpdateArticleContentDataAttributes
type UpdateArticleContentDataAttributes struct {
	Content []Section `json:"content,omitempty"`
}

// NewUpdateArticleContentDataAttributes instantiates a new UpdateArticleContentDataAttributes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateArticleContentDataAttributes() *UpdateArticleContentDataAttributes {
	this := UpdateArticleContentDataAttributes{}
	return &this
}

// NewUpdateArticleContentDataAttributesWithDefaults instantiates a new UpdateArticleContentDataAttributes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateArticleContentDataAttributesWithDefaults() *UpdateArticleContentDataAttributes {
	this := UpdateArticleContentDataAttributes{}
	return &this
}

// GetContent returns the Content field value if set, zero value otherwise.
func (o *UpdateArticleContentDataAttributes) GetContent() []Section {
	if o == nil || IsNil(o.Content) {
		var ret []Section
		return ret
	}
	return o.Content
}

// GetContentOk returns a tuple with the Content field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateArticleContentDataAttributes) GetContentOk() ([]Section, bool) {
	if o == nil || IsNil(o.Content) {
		return nil, false
	}
	return o.Content, true
}

// HasContent returns a boolean if a field has been set.
func (o *UpdateArticleContentDataAttributes) HasContent() bool {
	if o != nil && !IsNil(o.Content) {
		return true
	}

	return false
}

// SetContent gets a reference to the given []Section and assigns it to the Content field.
func (o *UpdateArticleContentDataAttributes) SetContent(v []Section) {
	o.Content = v
}

func (o UpdateArticleContentDataAttributes) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateArticleContentDataAttributes) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Content) {
		toSerialize["content"] = o.Content
	}
	return toSerialize, nil
}

type NullableUpdateArticleContentDataAttributes struct {
	value *UpdateArticleContentDataAttributes
	isSet bool
}

func (v NullableUpdateArticleContentDataAttributes) Get() *UpdateArticleContentDataAttributes {
	return v.value
}

func (v *NullableUpdateArticleContentDataAttributes) Set(val *UpdateArticleContentDataAttributes) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateArticleContentDataAttributes) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateArticleContentDataAttributes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateArticleContentDataAttributes(val *UpdateArticleContentDataAttributes) *NullableUpdateArticleContentDataAttributes {
	return &NullableUpdateArticleContentDataAttributes{value: val, isSet: true}
}

func (v NullableUpdateArticleContentDataAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateArticleContentDataAttributes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


