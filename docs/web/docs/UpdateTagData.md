# UpdateTagData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | tag name | 
**Type** | **string** |  | 
**Attributes** | [**UpdateTagDataAttributes**](UpdateTagDataAttributes.md) |  | 

## Methods

### NewUpdateTagData

`func NewUpdateTagData(id string, type_ string, attributes UpdateTagDataAttributes, ) *UpdateTagData`

NewUpdateTagData instantiates a new UpdateTagData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateTagDataWithDefaults

`func NewUpdateTagDataWithDefaults() *UpdateTagData`

NewUpdateTagDataWithDefaults instantiates a new UpdateTagData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateTagData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateTagData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateTagData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateTagData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateTagData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateTagData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateTagData) GetAttributes() UpdateTagDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateTagData) GetAttributesOk() (*UpdateTagDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateTagData) SetAttributes(v UpdateTagDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


