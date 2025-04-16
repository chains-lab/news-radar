# ContentUpdateResponseAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the content update operation. | 
**Message** | **string** | A message providing additional information about the operation. | 
**Content** | [**Content**](Content.md) |  | 

## Methods

### NewContentUpdateResponseAttributes

`func NewContentUpdateResponseAttributes(status string, message string, content Content, ) *ContentUpdateResponseAttributes`

NewContentUpdateResponseAttributes instantiates a new ContentUpdateResponseAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContentUpdateResponseAttributesWithDefaults

`func NewContentUpdateResponseAttributesWithDefaults() *ContentUpdateResponseAttributes`

NewContentUpdateResponseAttributesWithDefaults instantiates a new ContentUpdateResponseAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *ContentUpdateResponseAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ContentUpdateResponseAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ContentUpdateResponseAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *ContentUpdateResponseAttributes) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ContentUpdateResponseAttributes) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ContentUpdateResponseAttributes) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetContent

`func (o *ContentUpdateResponseAttributes) GetContent() Content`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *ContentUpdateResponseAttributes) GetContentOk() (*Content, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *ContentUpdateResponseAttributes) SetContent(v Content)`

SetContent sets Content field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


