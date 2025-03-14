# ArticleData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | article ID | 
**Type** | **string** |  | 
**Attributes** | [**ArticleAttributes**](ArticleAttributes.md) |  | 

## Methods

### NewArticleData

`func NewArticleData(id string, type_ string, attributes ArticleAttributes, ) *ArticleData`

NewArticleData instantiates a new ArticleData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleDataWithDefaults

`func NewArticleDataWithDefaults() *ArticleData`

NewArticleDataWithDefaults instantiates a new ArticleData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ArticleData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ArticleData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ArticleData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *ArticleData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ArticleData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ArticleData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *ArticleData) GetAttributes() ArticleAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ArticleData) GetAttributesOk() (*ArticleAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ArticleData) SetAttributes(v ArticleAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


