// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/treeverse/lakefs/api/gen/models"
)

// ListObjectsReader is a Reader for the ListObjects structure.
type ListObjectsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListObjectsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListObjectsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListObjectsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewListObjectsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewListObjectsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListObjectsOK creates a ListObjectsOK with default headers values
func NewListObjectsOK() *ListObjectsOK {
	return &ListObjectsOK{}
}

/*ListObjectsOK handles this case with default header values.

entry list
*/
type ListObjectsOK struct {
	Payload *ListObjectsOKBody
}

func (o *ListObjectsOK) Error() string {
	return fmt.Sprintf("[GET /repositories/{repositoryId}/branches/{branchId}/objects/ls][%d] listObjectsOK  %+v", 200, o.Payload)
}

func (o *ListObjectsOK) GetPayload() *ListObjectsOKBody {
	return o.Payload
}

func (o *ListObjectsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ListObjectsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListObjectsUnauthorized creates a ListObjectsUnauthorized with default headers values
func NewListObjectsUnauthorized() *ListObjectsUnauthorized {
	return &ListObjectsUnauthorized{}
}

/*ListObjectsUnauthorized handles this case with default header values.

Unauthorized
*/
type ListObjectsUnauthorized struct {
	Payload *models.Error
}

func (o *ListObjectsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /repositories/{repositoryId}/branches/{branchId}/objects/ls][%d] listObjectsUnauthorized  %+v", 401, o.Payload)
}

func (o *ListObjectsUnauthorized) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListObjectsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListObjectsNotFound creates a ListObjectsNotFound with default headers values
func NewListObjectsNotFound() *ListObjectsNotFound {
	return &ListObjectsNotFound{}
}

/*ListObjectsNotFound handles this case with default header values.

tree or branch not found
*/
type ListObjectsNotFound struct {
	Payload *models.Error
}

func (o *ListObjectsNotFound) Error() string {
	return fmt.Sprintf("[GET /repositories/{repositoryId}/branches/{branchId}/objects/ls][%d] listObjectsNotFound  %+v", 404, o.Payload)
}

func (o *ListObjectsNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListObjectsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListObjectsDefault creates a ListObjectsDefault with default headers values
func NewListObjectsDefault(code int) *ListObjectsDefault {
	return &ListObjectsDefault{
		_statusCode: code,
	}
}

/*ListObjectsDefault handles this case with default header values.

generic error response
*/
type ListObjectsDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the list objects default response
func (o *ListObjectsDefault) Code() int {
	return o._statusCode
}

func (o *ListObjectsDefault) Error() string {
	return fmt.Sprintf("[GET /repositories/{repositoryId}/branches/{branchId}/objects/ls][%d] listObjects default  %+v", o._statusCode, o.Payload)
}

func (o *ListObjectsDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *ListObjectsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ListObjectsOKBody list objects o k body
swagger:model ListObjectsOKBody
*/
type ListObjectsOKBody struct {

	// pagination
	Pagination *models.Pagination `json:"pagination,omitempty"`

	// results
	Results []*models.ObjectStats `json:"results"`
}

// Validate validates this list objects o k body
func (o *ListObjectsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePagination(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListObjectsOKBody) validatePagination(formats strfmt.Registry) error {

	if swag.IsZero(o.Pagination) { // not required
		return nil
	}

	if o.Pagination != nil {
		if err := o.Pagination.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("listObjectsOK" + "." + "pagination")
			}
			return err
		}
	}

	return nil
}

func (o *ListObjectsOKBody) validateResults(formats strfmt.Registry) error {

	if swag.IsZero(o.Results) { // not required
		return nil
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listObjectsOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *ListObjectsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListObjectsOKBody) UnmarshalBinary(b []byte) error {
	var res ListObjectsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}