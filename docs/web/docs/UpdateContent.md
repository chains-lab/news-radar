# UpdateContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Content** | [**Section**](Section.md) |  | 

## Methods

### NewUpdateContent

`func NewUpdateContent(type_ string, content Section, ) *UpdateContent`

NewUpdateContent instantiates a new UpdateContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateContentWithDefaults

`func NewUpdateContentWithDefaults() *UpdateContent`

NewUpdateContentWithDefaults instantiates a new UpdateContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *UpdateContent) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateContent) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateContent) SetType(v string)`

SetType sets Type field to given value.


### GetContent

`func (o *UpdateContent) GetContent() Section`

GetContent returns the Content field if non-nil, zero value otherwise.

### GetContentOk

`func (o *UpdateContent) GetContentOk() (*Section, bool)`

GetContentOk returns a tuple with the Content field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContent

`func (o *UpdateContent) SetContent(v Section)`

SetContent sets Content field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


