# ArticleUpdateDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** | Article status | [optional] 
**Title** | Pointer to **string** | Article title | [optional] 
**Icon** | Pointer to **string** | Article link | [optional] 
**Desc** | Pointer to **string** | Article description | [optional] 
**Content** | Pointer to **map[string]interface{}** | Article content | [optional] 
**Authors** | Pointer to **[]string** | Authors ID uuid | [optional] 
**Tags** | Pointer to **[]string** |  | [optional] 

## Methods

### NewArticleUpdateDataAttributes

`func NewArticleUpdateDataAttributes() *ArticleUpdateDataAttributes`

NewArticleUpdateDataAttributes instantiates a new ArticleUpdateDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleUpdateDataAttributesWithDefaults

`func NewArticleUpdateDataAttributesWithDefaults() *ArticleUpdateDataAttributes`

NewArticleUpdateDataAttributesWithDefaults instantiates a new ArticleUpdateDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *ArticleUpdateDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ArticleUpdateDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ArticleUpdateDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ArticleUpdateDataAttributes) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTitle

`func (o *ArticleUpdateDataAttributes) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ArticleUpdateDataAttributes) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ArticleUpdateDataAttributes) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *ArticleUpdateDataAttributes) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetIcon

`func (o *ArticleUpdateDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *ArticleUpdateDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *ArticleUpdateDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.

### HasIcon

`func (o *ArticleUpdateDataAttributes) HasIcon() bool`

HasIcon returns a boolean if a field has been set.

### GetDesc

`func (o *ArticleUpdateDataAttributes) GetDesc() string`

GetDesc returns the Desc field if non-nil, zero value otherwise.

### GetDescOk

`func (o *ArticleUpdateDataAttributes) GetDescOk() (*string, bool)`

GetDescOk returns a tuple with the Desc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesc

`func (o *ArticleUpdateDataAttributes) SetDesc(v string)`

SetDesc sets Desc field to given value.

### HasDesc

`func (o *ArticleUpdateDataAttributes) HasDesc() bool`

HasDesc returns a boolean if a field has been set.

### GetContent

`func (o *ArticleUpdateDataAttributes) GetContent() map[string]interface{}`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ArticleUpdateDataAttributes) GetContentOk() (*map[string]interface{}, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ArticleUpdateDataAttributes) SetContent(v map[string]interface{})`

SetContent sets Content field to given value.

### HasContent

`func (o *ArticleUpdateDataAttributes) HasContent() bool`

HasContent returns a boolean if a field has been set.

### GetAuthors

`func (o *ArticleUpdateDataAttributes) GetAuthors() []string`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ArticleUpdateDataAttributes) GetAuthorsOk() (*[]string, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ArticleUpdateDataAttributes) SetAuthors(v []string)`

SetAuthors sets Authors field to given value.

### HasAuthors

`func (o *ArticleUpdateDataAttributes) HasAuthors() bool`

HasAuthors returns a boolean if a field has been set.

### GetTags

`func (o *ArticleUpdateDataAttributes) GetTags() []string`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ArticleUpdateDataAttributes) GetTagsOk() (*[]string, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ArticleUpdateDataAttributes) SetTags(v []string)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ArticleUpdateDataAttributes) HasTags() bool`

HasTags returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


