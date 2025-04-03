# ArticleAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | Article title | 
**Status** | **string** | Article status | 
**Icon** | **string** | Article link | 
**Desc** | **string** | Article description | 
**Content** | [**Content**](Content.md) |  | 
**Likes** | **int32** | Article likes | 
**Reposts** | **int32** | Article reposts | 
**UpdatedAt** | Pointer to **time.Time** | Updated at | [optional] 
**CreatedAt** | **time.Time** | Created at | 

## Methods

### NewArticleAttributes

`func NewArticleAttributes(title string, status string, icon string, desc string, content Content, likes int32, reposts int32, createdAt time.Time, ) *ArticleAttributes`

NewArticleAttributes instantiates a new ArticleAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleAttributesWithDefaults

`func NewArticleAttributesWithDefaults() *ArticleAttributes`

NewArticleAttributesWithDefaults instantiates a new ArticleAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *ArticleAttributes) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ArticleAttributes) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ArticleAttributes) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetStatus

`func (o *ArticleAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ArticleAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ArticleAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetIcon

`func (o *ArticleAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *ArticleAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *ArticleAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.


### GetDesc

`func (o *ArticleAttributes) GetDesc() string`

GetDesc returns the Desc field if non-nil, zero value otherwise.

### GetDescOk

`func (o *ArticleAttributes) GetDescOk() (*string, bool)`

GetDescOk returns a tuple with the Desc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesc

`func (o *ArticleAttributes) SetDesc(v string)`

SetDesc sets Desc field to given value.


### GetContent

`func (o *ArticleAttributes) GetContent() Content`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ArticleAttributes) GetContentOk() (*Content, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ArticleAttributes) SetContent(v Content)`

SetContent sets Content field to given value.


### GetLikes

`func (o *ArticleAttributes) GetLikes() int32`

GetLikes returns the Likes field if non-nil, zero value otherwise.

### GetLikesOk

`func (o *ArticleAttributes) GetLikesOk() (*int32, bool)`

GetLikesOk returns a tuple with the Likes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLikes

`func (o *ArticleAttributes) SetLikes(v int32)`

SetLikes sets Likes field to given value.


### GetReposts

`func (o *ArticleAttributes) GetReposts() int32`

GetReposts returns the Reposts field if non-nil, zero value otherwise.

### GetRepostsOk

`func (o *ArticleAttributes) GetRepostsOk() (*int32, bool)`

GetRepostsOk returns a tuple with the Reposts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReposts

`func (o *ArticleAttributes) SetReposts(v int32)`

SetReposts sets Reposts field to given value.


### GetUpdatedAt

`func (o *ArticleAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ArticleAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ArticleAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ArticleAttributes) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetCreatedAt

`func (o *ArticleAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ArticleAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ArticleAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


