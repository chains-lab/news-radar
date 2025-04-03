# AuthorAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Status** | Pointer to **string** |  | [optional] 
**Desc** | **string** |  | 
**Avatar** | **string** |  | 
**Email** | Pointer to **string** |  | [optional] 
**Telegram** | Pointer to **string** |  | [optional] 
**Twitter** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewAuthorAttributes

`func NewAuthorAttributes(name string, desc string, avatar string, ) *AuthorAttributes`

NewAuthorAttributes instantiates a new AuthorAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthorAttributesWithDefaults

`func NewAuthorAttributesWithDefaults() *AuthorAttributes`

NewAuthorAttributesWithDefaults instantiates a new AuthorAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AuthorAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AuthorAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AuthorAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *AuthorAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *AuthorAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *AuthorAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *AuthorAttributes) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetDesc

`func (o *AuthorAttributes) GetDesc() string`

GetDesc returns the Desc field if non-nil, zero value otherwise.

### GetDescOk

`func (o *AuthorAttributes) GetDescOk() (*string, bool)`

GetDescOk returns a tuple with the Desc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDesc

`func (o *AuthorAttributes) SetDesc(v string)`

SetDesc sets Desc field to given value.


### GetAvatar

`func (o *AuthorAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *AuthorAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *AuthorAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.


### GetEmail

`func (o *AuthorAttributes) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AuthorAttributes) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AuthorAttributes) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *AuthorAttributes) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetTelegram

`func (o *AuthorAttributes) GetTelegram() string`

GetTelegram returns the Telegram field if non-nil, zero value otherwise.

### GetTelegramOk

`func (o *AuthorAttributes) GetTelegramOk() (*string, bool)`

GetTelegramOk returns a tuple with the Telegram field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTelegram

`func (o *AuthorAttributes) SetTelegram(v string)`

SetTelegram sets Telegram field to given value.

### HasTelegram

`func (o *AuthorAttributes) HasTelegram() bool`

HasTelegram returns a boolean if a field has been set.

### GetTwitter

`func (o *AuthorAttributes) GetTwitter() string`

GetTwitter returns the Twitter field if non-nil, zero value otherwise.

### GetTwitterOk

`func (o *AuthorAttributes) GetTwitterOk() (*string, bool)`

GetTwitterOk returns a tuple with the Twitter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTwitter

`func (o *AuthorAttributes) SetTwitter(v string)`

SetTwitter sets Twitter field to given value.

### HasTwitter

`func (o *AuthorAttributes) HasTwitter() bool`

HasTwitter returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *AuthorAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *AuthorAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *AuthorAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *AuthorAttributes) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetCreatedAt

`func (o *AuthorAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *AuthorAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *AuthorAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *AuthorAttributes) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


