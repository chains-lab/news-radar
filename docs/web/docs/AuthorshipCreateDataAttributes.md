# AuthorshipCreateDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorID** | Pointer to **[]string** |  | [optional] 
**ArticleID** | **string** | The ID of the article. | 

## Methods

### NewAuthorshipCreateDataAttributes

`func NewAuthorshipCreateDataAttributes(articleID string, ) *AuthorshipCreateDataAttributes`

NewAuthorshipCreateDataAttributes instantiates a new AuthorshipCreateDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthorshipCreateDataAttributesWithDefaults

`func NewAuthorshipCreateDataAttributesWithDefaults() *AuthorshipCreateDataAttributes`

NewAuthorshipCreateDataAttributesWithDefaults instantiates a new AuthorshipCreateDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthorID

`func (o *AuthorshipCreateDataAttributes) GetAuthorID() []string`

GetAuthorID returns the AuthorID field if non-nil, zero value otherwise.

### GetAuthorIDOk

`func (o *AuthorshipCreateDataAttributes) GetAuthorIDOk() (*[]string, bool)`

GetAuthorIDOk returns a tuple with the AuthorID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorID

`func (o *AuthorshipCreateDataAttributes) SetAuthorID(v []string)`

SetAuthorID sets AuthorID field to given value.

### HasAuthorID

`func (o *AuthorshipCreateDataAttributes) HasAuthorID() bool`

HasAuthorID returns a boolean if a field has been set.

### GetArticleID

`func (o *AuthorshipCreateDataAttributes) GetArticleID() string`

GetArticleID returns the ArticleID field if non-nil, zero value otherwise.

### GetArticleIDOk

`func (o *AuthorshipCreateDataAttributes) GetArticleIDOk() (*string, bool)`

GetArticleIDOk returns a tuple with the ArticleID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArticleID

`func (o *AuthorshipCreateDataAttributes) SetArticleID(v string)`

SetArticleID sets ArticleID field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


