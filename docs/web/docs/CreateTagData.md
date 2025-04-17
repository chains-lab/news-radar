# CreateTagData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Attributes** | [**CreateTagDataAttributes**](CreateTagDataAttributes.md) |  | 

## Methods

### NewCreateTagData

`func NewCreateTagData(type_ string, attributes CreateTagDataAttributes, ) *CreateTagData`

NewCreateTagData instantiates a new CreateTagData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateTagDataWithDefaults

`func NewCreateTagDataWithDefaults() *CreateTagData`

NewCreateTagDataWithDefaults instantiates a new CreateTagData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *CreateTagData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CreateTagData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CreateTagData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *CreateTagData) GetAttributes() CreateTagDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *CreateTagData) GetAttributesOk() (*CreateTagDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *CreateTagData) SetAttributes(v CreateTagDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


