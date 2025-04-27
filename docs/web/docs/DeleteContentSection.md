# DeleteContentSection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**SectionId** | **int32** | Unique identifier for the content section to be deleted. | 

## Methods

### NewDeleteContentSection

`func NewDeleteContentSection(type_ string, sectionId int32, ) *DeleteContentSection`

NewDeleteContentSection instantiates a new DeleteContentSection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteContentSectionWithDefaults

`func NewDeleteContentSectionWithDefaults() *DeleteContentSection`

NewDeleteContentSectionWithDefaults instantiates a new DeleteContentSection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *DeleteContentSection) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DeleteContentSection) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DeleteContentSection) SetType(v string)`

SetType sets Type field to given value.


### GetSectionId

`func (o *DeleteContentSection) GetSectionId() int32`

GetSectionId returns the SectionId field if non-nil, zero value otherwise.

### GetSectionIdOk

`func (o *DeleteContentSection) GetSectionIdOk() (*int32, bool)`

GetSectionIdOk returns a tuple with the SectionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionId

`func (o *DeleteContentSection) SetSectionId(v int32)`

SetSectionId sets SectionId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


