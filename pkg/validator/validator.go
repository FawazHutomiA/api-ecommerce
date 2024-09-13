package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"example/pkg/response"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/id_ID"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"

	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func ValidateRequest(req *http.Request, referenceStruct interface{}) (res response.Response, err error) {

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return response.Error(response.StatusUnprocessableEntity, ERR_INPUT, []string{ERR_EMPTY}), err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return response.Error(response.StatusUnprocessableEntity, ERR_INPUT, []string{ERR_EMPTY}), errors.New(ERR_EMPTY)
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(referenceStruct); err != nil {
		var errValue interface{}
		text := ERR_DATATYPE
		switch err.(type) {
		case *json.UnmarshalTypeError:
			typeError := err.(*json.UnmarshalTypeError)
			text = ERR_INPUT
			errValue = []ApiError{{typeError.Field, fmt.Sprintf("expected type %s, got %s", typeError.Type, typeError.Value)}}
		default:
			text = ERR_WRONG
			errValue = []string{fmt.Sprintf("failed when parsing body: %v", err)}
		}
		return response.Error(response.StatusUnprocessableEntity, text, errValue), errors.New(ERR_DATATYPE)
	}

	val := reflect.ValueOf(referenceStruct)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return response.Error(response.StatusUnprocessableEntity, ERR_WRONG, []string{"referenceStruct should be a non-nil pointer"}), errors.New(ERR_WRONG)
	}

	ptrType := reflect.TypeOf(referenceStruct)
	if ptrType.Elem().Kind() != reflect.Slice {
		res := validateStruct(referenceStruct)
		if res != nil {
			return response.Error(response.StatusUnprocessableEntity, ERR_INPUT, res), errors.New(ERR_VALIDATE)
		}
		return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, res), nil
	}

	for i := 0; i < val.Elem().Len(); i++ {
		element := val.Elem().Index(i)
		if err := validateStruct(element.Interface()); err != nil {
			return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, err), nil
		}
	}

	return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, err), nil
}

func ValidateFormRequest(req *http.Request, referenceStruct interface{}) (res response.Response, err error) {

	err = req.ParseForm()
	if err != nil {
		return response.Error(response.StatusUnprocessableEntity, ERR_INPUT, []string{ERR_EMPTY}), err
	}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(referenceStruct, req.Form); err != nil {
		return response.Error(response.StatusUnprocessableEntity, ERR_WRONG, []string{fmt.Sprintf("failed when parsing body: %v", err)}), errors.New(ERR_DATATYPE)
	}

	val := reflect.ValueOf(referenceStruct)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return response.Error(response.StatusUnprocessableEntity, ERR_WRONG, []string{"referenceStruct should be a non-nil pointer"}), errors.New(ERR_WRONG)
	}

	ptrType := reflect.TypeOf(referenceStruct)
	if ptrType.Elem().Kind() != reflect.Slice {
		res := validateStruct(referenceStruct)
		if res != nil {
			return response.Error(response.StatusUnprocessableEntity, ERR_INPUT, res), errors.New(ERR_VALIDATE)
		}
		return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, res), nil
	}

	for i := 0; i < val.Elem().Len(); i++ {
		element := val.Elem().Index(i)
		if err := validateStruct(element.Interface()); err != nil {
			return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, err), nil
		}
	}

	return response.Error(response.StatusUnprocessableEntity, ACC_INPUT, err), nil
}

func validateStruct(referenceStruct interface{}) (res interface{}) {

	validate := validator.New()

	//translation
	en := en_US.New()
	id := id_ID.New()
	uni := ut.New(en, en, id)
	transId, _ := uni.GetTranslator("id")
	idTranslations.RegisterDefaultTranslations(validate, transId)

	useJsonFieldValidation(validate)
	validate.RegisterValidation("maxSizeFile", validateMaxSizeFile)
	validate.RegisterValidation("typeFile", validateTypeFile)
	validate.RegisterValidation("requiredFile", validateRequiredFile)

	// Validate the fields of the structure using validator/v10
	if err := validate.Struct(referenceStruct); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {

			errors := err.(validator.ValidationErrors)
			out := make([]interface{}, len(errors))
			for i, e := range errors {
				out[i] = ApiError{e.Field(), msgForTag(e.Tag(), e.Param(), e.Translate(transId))}
			}

			return out
		}
	}

	return res
}

func msgForTag(tag string, param string, msgError string) string {
	switch tag {
	case "required_if":
		return "Kolom ini wajib diisi"
	case "phoneNumber":
		return "Format nomor telepon tidak valid"
	case "maxSizeFile":
		return "Ukuran file harus kurang dari atau sama dengan " + param + " Mb"
	case "typeFile":
		return "Jenis file tidak valid. Jenis file harus salah satu dari [" + param + "]"
	case "requiredFile":
		return "Kolom ini wajib diisi"
	}
	return msgError
}

func useJsonFieldValidation(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func validateMaxSizeFile(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	fileSize := fileHeader.Size
	maxSize := fl.Param()

	if maxSize == "" {
		return true
	}
	maxSizeInt, err := strconv.Atoi(maxSize)
	if err != nil {
		return false
	}
	maxSizeInt64 := int64(maxSizeInt * 1024 * 1024)

	return fileSize <= maxSizeInt64
}

func validateTypeFile(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	fileType := strings.Split(fileHeader.Header.Get("Content-Type"), "/")
	validFileType := strings.Split(fl.Param(), " ")

	if len(fileType) == 1 {
		return true
	}

	for _, v := range validFileType {
		if fileType[1] == v {
			return true
		}
	}
	return false
}

func validateRequiredFile(fl validator.FieldLevel) bool {
	fileHeader := fl.Field().Interface().(multipart.FileHeader)
	fileSize := fileHeader.Size
	return fileSize != 0
}
