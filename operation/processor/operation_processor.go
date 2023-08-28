package processor

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var did string
	var didtype currencytypes.DuplicationType
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccounts:
		fact, ok := t.Fact().(currency.CreateAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.KeyUpdater:
		fact, ok := t.Fact().(currency.KeyUpdaterFact)
		if !ok {
			return errors.Errorf("expected KeyUpdaterFact, not %T", t.Fact())
		}
		did = fact.Target().String()
		didtype = DuplicationTypeSender
	case currency.Transfers:
		fact, ok := t.Fact().(currency.TransfersFact)
		if !ok {
			return errors.Errorf("expected TransfersFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case extension.CreateContractAccounts:
		fact, ok := t.Fact().(extension.CreateContractAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
	case extension.Withdraws:
		fact, ok := t.Fact().(extension.WithdrawsFact)
		if !ok {
			return errors.Errorf("expected WithdrawsFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case currency.CurrencyRegister:
		fact, ok := t.Fact().(currency.CurrencyRegisterFact)
		if !ok {
			return errors.Errorf("expected CurrencyRegisterFact, not %T", t.Fact())
		}
		did = fact.Currency().Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.CurrencyPolicyUpdater:
		fact, ok := t.Fact().(currency.CurrencyPolicyUpdaterFact)
		if !ok {
			return errors.Errorf("expected CurrencyPolicyUpdaterFact, not %T", t.Fact())
		}
		did = fact.Currency().String()
		didtype = DuplicationTypeCurrency
	case currency.SuffrageInflation:
	case nft.CollectionRegister:
		fact, ok := t.Fact().(nft.CollectionRegisterFact)
		if !ok {
			return errors.Errorf("expected CollectionRegisterFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.CollectionPolicyUpdater:
		fact, ok := t.Fact().(nft.CollectionPolicyUpdaterFact)
		if !ok {
			return errors.Errorf("expected CollectionPolicyUpdaterFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.Mint:
		fact, ok := t.Fact().(nft.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.NFTTransfer:
		fact, ok := t.Fact().(nft.NFTTransferFact)
		if !ok {
			return errors.Errorf("expected NFTTransferFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.Delegate:
		fact, ok := t.Fact().(nft.DelegateFact)
		if !ok {
			return errors.Errorf("expected DelegateFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.Approve:
		fact, ok := t.Fact().(nft.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	case nft.NFTSign:
		fact, ok := t.Fact().(nft.NFTSignFact)
		if !ok {
			return errors.Errorf("expected NFTSignFact, not %T", t.Fact())
		}
		did = fact.Sender().String()
		didtype = DuplicationTypeSender
	default:
		return nil
	}

	if len(did) > 0 {
		if _, found := opr.Duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicate currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.Duplicated[did] = didtype
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op base.Operation) (base.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccounts,
		currency.KeyUpdater,
		currency.Transfers,
		extension.CreateContractAccounts,
		extension.Withdraws,
		currency.CurrencyRegister,
		currency.CurrencyPolicyUpdater,
		currency.SuffrageInflation,
		nft.CollectionRegister,
		nft.CollectionPolicyUpdater,
		nft.Mint,
		nft.NFTTransfer,
		nft.Delegate,
		nft.Approve,
		nft.NFTSign:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
