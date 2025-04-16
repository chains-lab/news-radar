# ContentUpdateAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SectionID** | **string** | The ID of the section to be updated. | 
**Action** | **string** | The action to be performed on the content. | 
**Content** | [**Content**](Content.md) |  | 

## Methods

### NewContentUpdateAttributes

`func NewContentUpdateAttributes(sectionID string, action string, content Content, ) *ContentUpdateAttributes`

NewContentUpdateAttributes instantiates a new ContentUpdateAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContentUpdateAttributesWithDefaults

`func NewContentUpdateAttributesWithDefaults() *ContentUpdateAttributes`

NewContentUpdateAttributesWithDefaults instantiates a new ContentUpdateAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSectionID

`func (o *ContentUpdateAttributes) GetSectionID() string`

GetSectionID returns the SectionID field if non-nil, zero value otherwise.

### GetSectionIDOk

`func (o *ContentUpdateAttributes) GetSectionIDOk() (*string, bool)`

GetSectionIDOk returns a tuple with the SectionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionID

`func (o *ContentUpdateAttributes) SetSectionID(v string)`

SetSectionID sets SectionID field to given value.


### GetAction

`func (o *ContentUpdateAttributes) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *ContentUpdateAttributes) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *ContentUpdateAttributes) SetAction(v string)`

SetAction sets Action field to given value.


### GetContent

`func (o *ContentUpdateAttributes) GetContent() Content`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ContentUpdateAttributes) GetContentOk() (*Content, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ContentUpdateAttributes) SetContent(v Content)`

SetContent sets Content field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


