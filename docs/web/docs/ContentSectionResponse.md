# ContentSectionResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | The status of the content update operation. | 
**Code** | **int32** | A code representing the result of the operation. | 
**Message** | **string** | A message providing additional information about the operation. | 
**Type** | **string** |  | 
**Section** | [**Section**](Section.md) |  | 

## Methods

### NewContentSectionResponse

`func NewContentSectionResponse(status string, code int32, message string, type_ string, section Section, ) *ContentSectionResponse`

NewContentSectionResponse instantiates a new ContentSectionResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContentSectionResponseWithDefaults

`func NewContentSectionResponseWithDefaults() *ContentSectionResponse`

NewContentSectionResponseWithDefaults instantiates a new ContentSectionResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *ContentSectionResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ContentSectionResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ContentSectionResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetCode

`func (o *ContentSectionResponse) GetCode() int32`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ContentSectionResponse) GetCodeOk() (*int32, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ContentSectionResponse) SetCode(v int32)`

SetCode sets Code field to given value.


### GetMessage

`func (o *ContentSectionResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ContentSectionResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ContentSectionResponse) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetType

`func (o *ContentSectionResponse) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ContentSectionResponse) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ContentSectionResponse) SetType(v string)`

SetType sets Type field to given value.


### GetSection

`func (o *ContentSectionResponse) GetSection() Section`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *ContentSectionResponse) GetSectionOk() (*Section, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *ContentSectionResponse) SetSection(v Section)`

SetSection sets Section field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


