# Section

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | Unique identifier for the content in this content block. (Num in section) | 
**Type** | **string** | The type of content in this section. | 
**Text** | Pointer to [**[]SectionTextInner**](SectionTextInner.md) |  | [optional] 
**Media** | Pointer to [**[]SectionMediaInner**](SectionMediaInner.md) |  | [optional] 
**Audio** | Pointer to [**[]SectionAudioInner**](SectionAudioInner.md) |  | [optional] 

## Methods

### NewSection

`func NewSection(id int32, type_ string, ) *Section`

NewSection instantiates a new Section object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSectionWithDefaults

`func NewSectionWithDefaults() *Section`

NewSectionWithDefaults instantiates a new Section object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Section) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Section) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Section) SetId(v int32)`

SetId sets Id field to given value.


### GetType

`func (o *Section) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Section) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Section) SetType(v string)`

SetType sets Type field to given value.


### GetText

`func (o *Section) GetText() []SectionTextInner`

GetText returns the Text field if non-nil, zero value otherwise.

### GetTextOk

`func (o *Section) GetTextOk() (*[]SectionTextInner, bool)`

GetTextOk returns a tuple with the Text field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetText

`func (o *Section) SetText(v []SectionTextInner)`

SetText sets Text field to given value.

### HasText

`func (o *Section) HasText() bool`

HasText returns a boolean if a field has been set.

### GetMedia

`func (o *Section) GetMedia() []SectionMediaInner`

GetMedia returns the Media field if non-nil, zero value otherwise.

### GetMediaOk

`func (o *Section) GetMediaOk() (*[]SectionMediaInner, bool)`

GetMediaOk returns a tuple with the Media field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMedia

`func (o *Section) SetMedia(v []SectionMediaInner)`

SetMedia sets Media field to given value.

### HasMedia

`func (o *Section) HasMedia() bool`

HasMedia returns a boolean if a field has been set.

### GetAudio

`func (o *Section) GetAudio() []SectionAudioInner`

GetAudio returns the Audio field if non-nil, zero value otherwise.

### GetAudioOk

`func (o *Section) GetAudioOk() (*[]SectionAudioInner, bool)`

GetAudioOk returns a tuple with the Audio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudio

`func (o *Section) SetAudio(v []SectionAudioInner)`

SetAudio sets Audio field to given value.

### HasAudio

`func (o *Section) HasAudio() bool`

HasAudio returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


