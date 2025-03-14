# Article

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**ArticleData**](ArticleData.md) |  | 
**Included** | [**ArticleInclude**](ArticleInclude.md) |  | 

## Methods

### NewArticle

`func NewArticle(data ArticleData, included ArticleInclude, ) *Article`

NewArticle instantiates a new Article object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleWithDefaults

`func NewArticleWithDefaults() *Article`

NewArticleWithDefaults instantiates a new Article object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *Article) GetData() ArticleData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Article) GetDataOk() (*ArticleData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Article) SetData(v ArticleData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *Article) GetIncluded() ArticleInclude`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *Article) GetIncludedOk() (*ArticleInclude, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *Article) SetIncluded(v ArticleInclude)`

SetIncluded sets Included field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


