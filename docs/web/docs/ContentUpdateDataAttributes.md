# ContentUpdateDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SectionID** | **int32** | The ID of the section to be updated. | 
**Action** | **string** | The action to be performed on the content. | 
**Content** | [**Content**](Content.md) |  | 

## Methods

### NewContentUpdateDataAttributes

`func NewContentUpdateDataAttributes(sectionID int32, action string, content Content, ) *ContentUpdateDataAttributes`

NewContentUpdateDataAttributes instantiates a new ContentUpdateDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContentUpdateDataAttributesWithDefaults

`func NewContentUpdateDataAttributesWithDefaults() *ContentUpdateDataAttributes`

NewContentUpdateDataAttributesWithDefaults instantiates a new ContentUpdateDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSectionID

`func (o *ContentUpdateDataAttributes) GetSectionID() int32`

GetSectionID returns the SectionID field if non-nil, zero value otherwise.

### GetSectionIDOk

`func (o *ContentUpdateDataAttributes) GetSectionIDOk() (*int32, bool)`

GetSectionIDOk returns a tuple with the SectionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionID

`func (o *ContentUpdateDataAttributes) SetSectionID(v int32)`

SetSectionID sets SectionID field to given value.


### GetAction

`func (o *ContentUpdateDataAttributes) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *ContentUpdateDataAttributes) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *ContentUpdateDataAttributes) SetAction(v string)`

SetAction sets Action field to given value.


### GetContent

`func (o *ContentUpdateDataAttributes) GetContent() Content`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ContentUpdateDataAttributes) GetContentOk() (*Content, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ContentUpdateDataAttributes) SetContent(v Content)`

SetContent sets Content field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


