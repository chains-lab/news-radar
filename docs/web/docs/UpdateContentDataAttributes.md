# UpdateContentDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SectionID** | **int32** | The ID of the section to be updated. | 
**Content** | [**Section**](Section.md) |  | 

## Methods

### NewUpdateContentDataAttributes

`func NewUpdateContentDataAttributes(sectionID int32, content Section, ) *UpdateContentDataAttributes`

NewUpdateContentDataAttributes instantiates a new UpdateContentDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateContentDataAttributesWithDefaults

`func NewUpdateContentDataAttributesWithDefaults() *UpdateContentDataAttributes`

NewUpdateContentDataAttributesWithDefaults instantiates a new UpdateContentDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSectionID

`func (o *UpdateContentDataAttributes) GetSectionID() int32`

GetSectionID returns the SectionID field if non-nil, zero value otherwise.

### GetSectionIDOk

`func (o *UpdateContentDataAttributes) GetSectionIDOk() (*int32, bool)`

GetSectionIDOk returns a tuple with the SectionID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionID

`func (o *UpdateContentDataAttributes) SetSectionID(v int32)`

SetSectionID sets SectionID field to given value.


### GetContent

`func (o *UpdateContentDataAttributes) GetContent() Section`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *UpdateContentDataAttributes) GetContentOk() (*Section, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *UpdateContentDataAttributes) SetContent(v Section)`

SetContent sets Content field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


