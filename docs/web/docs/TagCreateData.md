# TagCreateData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Attributes** | [**TagCreateDataAttributes**](TagCreateDataAttributes.md) |  | 

## Methods

### NewTagCreateData

`func NewTagCreateData(type_ string, attributes TagCreateDataAttributes, ) *TagCreateData`

NewTagCreateData instantiates a new TagCreateData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTagCreateDataWithDefaults

`func NewTagCreateDataWithDefaults() *TagCreateData`

NewTagCreateDataWithDefaults instantiates a new TagCreateData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *TagCreateData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *TagCreateData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *TagCreateData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *TagCreateData) GetAttributes() TagCreateDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *TagCreateData) GetAttributesOk() (*TagCreateDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *TagCreateData) SetAttributes(v TagCreateDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


