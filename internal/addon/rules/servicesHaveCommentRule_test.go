package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser/meta"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestServicesHaveCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name                         string
		inputProto                   *parser.Proto
		inputShouldFollowGolangStyle bool
		wantFailures                 []report.Failure
	}{
		{
			name: "no failures for proto without service",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{},
			},
		},
		{
			name: "no failures for proto including valid services with comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "ServiceName",
						Comments: []*parser.Comment{
							{
								Raw: "// a service name.",
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto including valid services with Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "ServiceName",
						Comments: []*parser.Comment{
							{
								Raw: "// ServiceName is a service name.",
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
		},
		{
			name: "failures for proto with invalid services",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "ServiceName",
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   150,
								Line:     7,
								Column:   15,
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"SERVICES_HAVE_COMMENT",
					`Service "ServiceName" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid services without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "ServiceName",
						Comments: []*parser.Comment{
							{
								Raw: "// a service name.",
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   150,
								Line:     7,
								Column:   15,
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"SERVICES_HAVE_COMMENT",
					`Service "ServiceName" should have a comment of the form "// ServiceName ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewServicesHaveCommentRule(test.inputShouldFollowGolangStyle)

			got, err := rule.Apply(test.inputProto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}
