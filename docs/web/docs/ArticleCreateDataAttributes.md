# ArticleCreateDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | Article title | 
**Icon** | **string** | Article link | 
**Desc** | **string** | Article description | 
**Authors** | **[]string** | Authors ID uuid | 
**Content** | **string** | Article content | 
**Tags** | [**[]TagDataYaml**](TagDataYaml.md) |  | 
**Status** | **string** | Article status | 

## Methods

### NewArticleCreateDataAttributes

`func NewArticleCreateDataAttributes(title string, icon string, desc string, authors []string, content string, tags []TagDataYaml, status string, ) *ArticleCreateDataAttributes`

NewArticleCreateDataAttributes instantiates a new ArticleCreateDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewArticleCreateDataAttributesWithDefaults

`func NewArticleCreateDataAttributesWithDefaults() *ArticleCreateDataAttributes`

NewArticleCreateDataAttributesWithDefaults instantiates a new ArticleCreateDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *ArticleCreateDataAttributes) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ArticleCreateDataAttributes) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ArticleCreateDataAttributes) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetIcon

`func (o *ArticleCreateDataAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *ArticleCreateDataAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *ArticleCreateDataAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.


### GetDesc

`func (o *ArticleCreateDataAttributes) GetDesc() string`

GetDesc returns the Desc field if non-nil, zero value otherwise.

### GetDescOk

`func (o *ArticleCreateDataAttributes) GetDescOk() (*string, bool)`

GetDescOk returns a tuple with the Desc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesc

`func (o *ArticleCreateDataAttributes) SetDesc(v string)`

SetDesc sets Desc field to given value.


### GetAuthors

`func (o *ArticleCreateDataAttributes) GetAuthors() []string`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ArticleCreateDataAttributes) GetAuthorsOk() (*[]string, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ArticleCreateDataAttributes) SetAuthors(v []string)`

SetAuthors sets Authors field to given value.


### GetContent

`func (o *ArticleCreateDataAttributes) GetContent() string`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ArticleCreateDataAttributes) GetContentOk() (*string, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ArticleCreateDataAttributes) SetContent(v string)`

SetContent sets Content field to given value.


### GetTags

`func (o *ArticleCreateDataAttributes) GetTags() []TagDataYaml`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ArticleCreateDataAttributes) GetTagsOk() (*[]TagDataYaml, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ArticleCreateDataAttributes) SetTags(v []TagDataYaml)`

SetTags sets Tags field to given value.


### GetStatus

`func (o *ArticleCreateDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ArticleCreateDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ArticleCreateDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


