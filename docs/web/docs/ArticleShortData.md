# ArticleShortData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | article ID | 
**Type** | **string** |  | 
**Attributes** | [**ArticleShortDataAttributes**](ArticleShortDataAttributes.md) |  | 
**Relationships** | Pointer to [**ArticleShortDataRelationships**](ArticleShortDataRelationships.md) |  | [optional] 

## Methods

### NewArticleShortData

`func NewArticleShortData(id string, type_ string, attributes ArticleShortDataAttributes, ) *ArticleShortData`

NewArticleShortData instantiates a new ArticleShortData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleShortDataWithDefaults

`func NewArticleShortDataWithDefaults() *ArticleShortData`

NewArticleShortDataWithDefaults instantiates a new ArticleShortData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ArticleShortData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ArticleShortData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ArticleShortData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *ArticleShortData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ArticleShortData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ArticleShortData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *ArticleShortData) GetAttributes() ArticleShortDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *ArticleShortData) GetAttributesOk() (*ArticleShortDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *ArticleShortData) SetAttributes(v ArticleShortDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *ArticleShortData) GetRelationships() ArticleShortDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *ArticleShortData) GetRelationshipsOk() (*ArticleShortDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *ArticleShortData) SetRelationships(v ArticleShortDataRelationships)`

SetRelationships sets Relationships field to given value.

### HasRelationships

`func (o *ArticleShortData) HasRelationships() bool`

HasRelationships returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


