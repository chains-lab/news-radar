# UpdateContentSectionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the content update operation. | 
**Message** | **string** | A message providing additional information about the operation. | 
**Type** | **string** |  | 
**Section** | [**Section**](Section.md) |  | 

## Methods

### NewUpdateContentSectionResponse

`func NewUpdateContentSectionResponse(status string, message string, type_ string, section Section, ) *UpdateContentSectionResponse`

NewUpdateContentSectionResponse instantiates a new UpdateContentSectionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateContentSectionResponseWithDefaults

`func NewUpdateContentSectionResponseWithDefaults() *UpdateContentSectionResponse`

NewUpdateContentSectionResponseWithDefaults instantiates a new UpdateContentSectionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *UpdateContentSectionResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdateContentSectionResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdateContentSectionResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *UpdateContentSectionResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *UpdateContentSectionResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *UpdateContentSectionResponse) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetType

`func (o *UpdateContentSectionResponse) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateContentSectionResponse) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateContentSectionResponse) SetType(v string)`

SetType sets Type field to given value.


### GetSection

`func (o *UpdateContentSectionResponse) GetSection() Section`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *UpdateContentSectionResponse) GetSectionOk() (*Section, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *UpdateContentSectionResponse) SetSection(v Section)`

SetSection sets Section field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


