# ArticleInclude

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authors** | [**[]AuthorData**](AuthorData.md) |  | 
**Tags** | [**[]TagData**](TagData.md) |  | 

## Methods

### NewArticleInclude

`func NewArticleInclude(authors []AuthorData, tags []TagData, ) *ArticleInclude`

NewArticleInclude instantiates a new ArticleInclude object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleIncludeWithDefaults

`func NewArticleIncludeWithDefaults() *ArticleInclude`

NewArticleIncludeWithDefaults instantiates a new ArticleInclude object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthors

`func (o *ArticleInclude) GetAuthors() []AuthorData`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ArticleInclude) GetAuthorsOk() (*[]AuthorData, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ArticleInclude) SetAuthors(v []AuthorData)`

SetAuthors sets Authors field to given value.


### GetTags

`func (o *ArticleInclude) GetTags() []TagData`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ArticleInclude) GetTagsOk() (*[]TagData, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ArticleInclude) SetTags(v []TagData)`

SetTags sets Tags field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


