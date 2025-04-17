# UpdateContentResponseDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the content update operation. | 
**Message** | **string** | A message providing additional information about the operation. | 
**SectionNumber** | **int32** |  | 
**Content** | Pointer to [**Content**](Content.md) |  | [optional] 

## Methods

### NewUpdateContentResponseDataAttributes

`func NewUpdateContentResponseDataAttributes(status string, message string, sectionNumber int32, ) *UpdateContentResponseDataAttributes`

NewUpdateContentResponseDataAttributes instantiates a new UpdateContentResponseDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateContentResponseDataAttributesWithDefaults

`func NewUpdateContentResponseDataAttributesWithDefaults() *UpdateContentResponseDataAttributes`

NewUpdateContentResponseDataAttributesWithDefaults instantiates a new UpdateContentResponseDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *UpdateContentResponseDataAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdateContentResponseDataAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdateContentResponseDataAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *UpdateContentResponseDataAttributes) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *UpdateContentResponseDataAttributes) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *UpdateContentResponseDataAttributes) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetSectionNumber

`func (o *UpdateContentResponseDataAttributes) GetSectionNumber() int32`

GetSectionNumber returns the SectionNumber field if non-nil, zero value otherwise.

### GetSectionNumberOk

`func (o *UpdateContentResponseDataAttributes) GetSectionNumberOk() (*int32, bool)`

GetSectionNumberOk returns a tuple with the SectionNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionNumber

`func (o *UpdateContentResponseDataAttributes) SetSectionNumber(v int32)`

SetSectionNumber sets SectionNumber field to given value.


### GetContent

`func (o *UpdateContentResponseDataAttributes) GetContent() Content`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *UpdateContentResponseDataAttributes) GetContentOk() (*Content, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *UpdateContentResponseDataAttributes) SetContent(v Content)`

SetContent sets Content field to given value.

### HasContent

`func (o *UpdateContentResponseDataAttributes) HasContent() bool`

HasContent returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


