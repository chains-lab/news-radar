# ArticleDataRelationships

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authors** | [**[]Relationships**](Relationships.md) |  | 
**Tags** | [**[]Relationships**](Relationships.md) |  | 

## Methods

### NewArticleDataRelationships

`func NewArticleDataRelationships(authors []Relationships, tags []Relationships, ) *ArticleDataRelationships`

NewArticleDataRelationships instantiates a new ArticleDataRelationships object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleDataRelationshipsWithDefaults

`func NewArticleDataRelationshipsWithDefaults() *ArticleDataRelationships`

NewArticleDataRelationshipsWithDefaults instantiates a new ArticleDataRelationships object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthors

`func (o *ArticleDataRelationships) GetAuthors() []Relationships`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ArticleDataRelationships) GetAuthorsOk() (*[]Relationships, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ArticleDataRelationships) SetAuthors(v []Relationships)`

SetAuthors sets Authors field to given value.


### GetTags

`func (o *ArticleDataRelationships) GetTags() []Relationships`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ArticleDataRelationships) GetTagsOk() (*[]Relationships, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ArticleDataRelationships) SetTags(v []Relationships)`

SetTags sets Tags field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


