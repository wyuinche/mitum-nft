package nft

import (
	"fmt"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

func ErrStringPreProcess(t any) string {
	return fmt.Sprintf("failed to preprocess %T", t)
}

func ErrStringProcess(t any) string {
	return fmt.Sprintf("failed to process %T", t)
}

func ErrBaseOperationProcess(msg, k string, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf("%s, %s", msg, k), e))
}

func BaseErrStateNotFound(name string, k string, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf("%s not found, %s", name, k), e))
}

func BaseErrStateAlreadyExists(name, k string, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(fmt.Sprintf("%s already exists, %s", name, k), e))
}

func BaseErrInvalid(t any, e error) base.BaseOperationProcessReasonError {
	return base.NewBaseOperationProcessReasonError(utils.ErrStringWrap(utils.ErrStringInvalid(t), e))
}

func ErrStateNotFound(name string, k string, e error) error {
	return errors.Errorf(utils.ErrStringWrap(fmt.Sprintf("%s not found, %s", name, k), e))
}

func ErrStateAlreadyExists(name, k string, e error) error {
	return errors.Errorf(utils.ErrStringWrap(fmt.Sprintf("%s already exists, %s", name, k), e))
}

func ErrInvalid(t any, e error) error {
	return errors.Errorf(utils.ErrStringWrap(utils.ErrStringInvalid(t), e))
}
