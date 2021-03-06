// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/grpc/v1/location.proto

package v1

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Category) Validate() error {
	return nil
}
func (this *Location) Validate() error {
	return nil
}
func (this *CreateCategoryRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	return nil
}
func (this *CreateCategoryResponse) Validate() error {
	if this.Category != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Category); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Category", err)
		}
	}
	return nil
}
func (this *GetCategoriesRequest) Validate() error {
	return nil
}
func (this *GetCategoriesResponse) Validate() error {
	for _, item := range this.Categories {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Categories", err)
			}
		}
	}
	return nil
}

var _regex_GetCategoryRequest_Id = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetCategoryRequest) Validate() error {
	if !_regex_GetCategoryRequest_Id.MatchString(this.Id) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.Id))
	}
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	return nil
}
func (this *GetCategoryResponse) Validate() error {
	if this.Category != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Category); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Category", err)
		}
	}
	return nil
}
func (this *UpdateCategoryRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	return nil
}
func (this *UpdateCategoryResponse) Validate() error {
	if this.Category != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Category); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Category", err)
		}
	}
	return nil
}

var _regex_DeleteCategoryRequest_Id = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *DeleteCategoryRequest) Validate() error {
	if !_regex_DeleteCategoryRequest_Id.MatchString(this.Id) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.Id))
	}
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	return nil
}
func (this *DeleteCategoryResponse) Validate() error {
	return nil
}

var _regex_CreateLocationRequest_CategoryId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *CreateLocationRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if this.Address == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Address", fmt.Errorf(`value '%v' must not be an empty string`, this.Address))
	}
	if !_regex_CreateLocationRequest_CategoryId.MatchString(this.CategoryId) {
		return github_com_mwitkow_go_proto_validators.FieldError("CategoryId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.CategoryId))
	}
	if this.CategoryId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("CategoryId", fmt.Errorf(`value '%v' must not be an empty string`, this.CategoryId))
	}
	return nil
}
func (this *CreateLocationResponse) Validate() error {
	if this.Location != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Location); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Location", err)
		}
	}
	return nil
}
func (this *GetLocationsRequest) Validate() error {
	return nil
}
func (this *GetLocationsResponse) Validate() error {
	for _, item := range this.Categories {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Categories", err)
			}
		}
	}
	return nil
}

var _regex_GetLocationRequest_Id = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetLocationRequest) Validate() error {
	if !_regex_GetLocationRequest_Id.MatchString(this.Id) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.Id))
	}
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	return nil
}
func (this *GetLocationResponse) Validate() error {
	if this.Location != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Location); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Location", err)
		}
	}
	return nil
}

var _regex_UpdateLocationRequest_CategoryId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *UpdateLocationRequest) Validate() error {
	if !_regex_UpdateLocationRequest_CategoryId.MatchString(this.CategoryId) {
		return github_com_mwitkow_go_proto_validators.FieldError("CategoryId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.CategoryId))
	}
	return nil
}
func (this *UpdateLocationResponse) Validate() error {
	if this.Location != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Location); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Location", err)
		}
	}
	return nil
}

var _regex_DeleteLocationRequest_Id = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *DeleteLocationRequest) Validate() error {
	if !_regex_DeleteLocationRequest_Id.MatchString(this.Id) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.Id))
	}
	if this.Id == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`value '%v' must not be an empty string`, this.Id))
	}
	return nil
}
func (this *DeleteLocationResponse) Validate() error {
	return nil
}
