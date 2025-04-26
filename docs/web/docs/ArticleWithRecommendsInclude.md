# ArticleWithRecommendsInclude

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authors** | [**[]AuthorData**](AuthorData.md) |  | 
**Tags** | [**[]TagData**](TagData.md) |  | 
**Recommends** | [**[]ArticleShortData**](ArticleShortData.md) |  | 

## Methods

### NewArticleWithRecommendsInclude

`func NewArticleWithRecommendsInclude(authors []AuthorData, tags []TagData, recommends []ArticleShortData, ) *ArticleWithRecommendsInclude`

NewArticleWithRecommendsInclude instantiates a new ArticleWithRecommendsInclude object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleWithRecommendsIncludeWithDefaults

`func NewArticleWithRecommendsIncludeWithDefaults() *ArticleWithRecommendsInclude`

NewArticleWithRecommendsIncludeWithDefaults instantiates a new ArticleWithRecommendsInclude object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthors

`func (o *ArticleWithRecommendsInclude) GetAuthors() []AuthorData`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ArticleWithRecommendsInclude) GetAuthorsOk() (*[]AuthorData, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ArticleWithRecommendsInclude) SetAuthors(v []AuthorData)`

SetAuthors sets Authors field to given value.


### GetTags

`func (o *ArticleWithRecommendsInclude) GetTags() []TagData`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ArticleWithRecommendsInclude) GetTagsOk() (*[]TagData, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ArticleWithRecommendsInclude) SetTags(v []TagData)`

SetTags sets Tags field to given value.


### GetRecommends

`func (o *ArticleWithRecommendsInclude) GetRecommends() []ArticleShortData`

GetRecommends returns the Recommends field if non-nil, zero value otherwise.

### GetRecommendsOk

`func (o *ArticleWithRecommendsInclude) GetRecommendsOk() (*[]ArticleShortData, bool)`

GetRecommendsOk returns a tuple with the Recommends field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecommends

`func (o *ArticleWithRecommendsInclude) SetRecommends(v []ArticleShortData)`

SetRecommends sets Recommends field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


