# TagUpdateData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | tag name | 
**Type** | **string** |  | 
**Attributes** | [**TagUpdateDataAttributes**](TagUpdateDataAttributes.md) |  | 

## Methods

### NewTagUpdateData

`func NewTagUpdateData(id string, type_ string, attributes TagUpdateDataAttributes, ) *TagUpdateData`

NewTagUpdateData instantiates a new TagUpdateData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTagUpdateDataWithDefaults

`func NewTagUpdateDataWithDefaults() *TagUpdateData`

NewTagUpdateDataWithDefaults instantiates a new TagUpdateData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TagUpdateData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TagUpdateData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TagUpdateData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *TagUpdateData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *TagUpdateData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *TagUpdateData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *TagUpdateData) GetAttributes() TagUpdateDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *TagUpdateData) GetAttributesOk() (*TagUpdateDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *TagUpdateData) SetAttributes(v TagUpdateDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


