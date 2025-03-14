# AuthorData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Author ID | 
**Type** | **string** |  | 
**Attributes** | [**AuthorAttributes**](AuthorAttributes.md) |  | 

## Methods

### NewAuthorData

`func NewAuthorData(id string, type_ string, attributes AuthorAttributes, ) *AuthorData`

NewAuthorData instantiates a new AuthorData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthorDataWithDefaults

`func NewAuthorDataWithDefaults() *AuthorData`

NewAuthorDataWithDefaults instantiates a new AuthorData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AuthorData) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AuthorData) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AuthorData) SetId(v string)`

SetId sets Id field to given value.


### GetType

`func (o *AuthorData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *AuthorData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *AuthorData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *AuthorData) GetAttributes() AuthorAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *AuthorData) GetAttributesOk() (*AuthorAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *AuthorData) SetAttributes(v AuthorAttributes)`

SetAttributes sets Attributes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


