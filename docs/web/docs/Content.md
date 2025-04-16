# Content

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Type** | **string** |  | 
**Text** | Pointer to [**[]ContentTextInner**](ContentTextInner.md) |  | [optional] 
**Media** | Pointer to [**ContentMedia**](ContentMedia.md) |  | [optional] 
**Audio** | Pointer to [**ContentAudio**](ContentAudio.md) |  | [optional] 

## Methods

### NewContent

`func NewContent(id string, type_ string, ) *Content`

NewContent instantiates a new Content object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContentWithDefaults

`func NewContentWithDefaults() *Content`

NewContentWithDefaults instantiates a new Content object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Content) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Content) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Content) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *Content) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Content) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Content) SetType(v string)`

SetType sets Type field to given value.


### GetText

`func (o *Content) GetText() []ContentTextInner`

GetText returns the Text field if non-nil, zero value otherwise.

### GetTextOk

`func (o *Content) GetTextOk() (*[]ContentTextInner, bool)`

GetTextOk returns a tuple with the Text field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetText

`func (o *Content) SetText(v []ContentTextInner)`

SetText sets Text field to given value.

### HasText

`func (o *Content) HasText() bool`

HasText returns a boolean if a field has been set.

### GetMedia

`func (o *Content) GetMedia() ContentMedia`

GetMedia returns the Media field if non-nil, zero value otherwise.

### GetMediaOk

`func (o *Content) GetMediaOk() (*ContentMedia, bool)`

GetMediaOk returns a tuple with the Media field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMedia

`func (o *Content) SetMedia(v ContentMedia)`

SetMedia sets Media field to given value.

### HasMedia

`func (o *Content) HasMedia() bool`

HasMedia returns a boolean if a field has been set.

### GetAudio

`func (o *Content) GetAudio() ContentAudio`

GetAudio returns the Audio field if non-nil, zero value otherwise.

### GetAudioOk

`func (o *Content) GetAudioOk() (*ContentAudio, bool)`

GetAudioOk returns a tuple with the Audio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudio

`func (o *Content) SetAudio(v ContentAudio)`

SetAudio sets Audio field to given value.

### HasAudio

`func (o *Content) HasAudio() bool`

HasAudio returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


