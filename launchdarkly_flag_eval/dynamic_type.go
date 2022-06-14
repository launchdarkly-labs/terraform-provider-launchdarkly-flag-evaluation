package launchdarkly_flag_eval

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type DynamicType struct{}

func (d DynamicType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.DynamicPseudoType
}

func (d DynamicType) ValueFromTerraform(_ context.Context, val tftypes.Value) (attr.Value, error) {
	vals := map[string]tftypes.Value{}
	err := val.As(&vals)
	return Dynamic{
		Value: vals,
	}, nil
}

func (d DynamicType) Equal(other attr.Type) bool {
	_, ok := other.(DynamicType)
	if !ok {
		return false
	}
	return true
}

func (d DynamicType) String() string {
	return "DynamicPseudoType"
}

func (d DynamicType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("can't step into dynamic pseudo type")
}

func (d DynamicType) Validate(_ context.Context, val tftypes.Value, path *tftypes.AttributePath) diag.Diagnostics {
	if !val.Type().Is(tftypes.Object{}) {
		return diag.NewAttributeErrorDiagnostic(path, "Invalid type", "Can only be an object")
	}
	return nil
}

type Dynamic struct {
	Value map[string]tftypes.Value
}

func (d Dynamic) Type(_ context.Context) attr.Type {
	return DynamicType{}
}

func (d Dynamic) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	return d.Value, nil
}

func (d Dynamic) Equal(other attr.Value) bool {
	o, ok := other.(Dynamic)
	if !ok {
		return false
	}
	return d.Value.Equal(o.Value)
}
