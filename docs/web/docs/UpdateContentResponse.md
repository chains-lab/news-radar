# UpdateContentResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the content update operation. | 
**Message** | **string** | A message providing additional information about the operation. | 
**Type** | **string** |  | 
**Section** | [**Section**](Section.md) |  | 

## Methods

### NewUpdateContentResponse

`func NewUpdateContentResponse(status string, message string, type_ string, section Section, ) *UpdateContentResponse`

NewUpdateContentResponse instantiates a new UpdateContentResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateContentResponseWithDefaults

`func NewUpdateContentResponseWithDefaults() *UpdateContentResponse`

NewUpdateContentResponseWithDefaults instantiates a new UpdateContentResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *UpdateContentResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UpdateContentResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UpdateContentResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetMessage

`func (o *UpdateContentResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *UpdateContentResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *UpdateContentResponse) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetType

`func (o *UpdateContentResponse) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateContentResponse) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateContentResponse) SetType(v string)`

SetType sets Type field to given value.


### GetSection

`func (o *UpdateContentResponse) GetSection() Section`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *UpdateContentResponse) GetSectionOk() (*Section, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *UpdateContentResponse) SetSection(v Section)`

SetSection sets Section field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


