# UpdateArticleContentData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Article ID uuid | 
**Type** | **string** |  | 
**Attributes** | [**UpdateArticleContentDataAttributes**](UpdateArticleContentDataAttributes.md) |  | 

## Methods

### NewUpdateArticleContentData

`func NewUpdateArticleContentData(id string, type_ string, attributes UpdateArticleContentDataAttributes, ) *UpdateArticleContentData`

NewUpdateArticleContentData instantiates a new UpdateArticleContentData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateArticleContentDataWithDefaults

`func NewUpdateArticleContentDataWithDefaults() *UpdateArticleContentData`

NewUpdateArticleContentDataWithDefaults instantiates a new UpdateArticleContentData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateArticleContentData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateArticleContentData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateArticleContentData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateArticleContentData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateArticleContentData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateArticleContentData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateArticleContentData) GetAttributes() UpdateArticleContentDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateArticleContentData) GetAttributesOk() (*UpdateArticleContentDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateArticleContentData) SetAttributes(v UpdateArticleContentDataAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


