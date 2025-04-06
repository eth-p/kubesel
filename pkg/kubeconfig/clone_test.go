package kubeconfig

import (
	"reflect"
	"testing"

	"github.com/eth-p/kubesel/internal/testutil"
)

var cloneTestcases = map[string]struct {
	Type reflect.Type
}{
	"Config": {
		Type: reflect.TypeFor[Config](),
	},
	"NamedCluster": {
		Type: reflect.TypeFor[NamedCluster](),
	},
	"Cluster": {
		Type: reflect.TypeFor[Cluster](),
	},
	"NamedContext": {
		Type: reflect.TypeFor[NamedContext](),
	},
	"Context": {
		Type: reflect.TypeFor[Context](),
	},
	"NamedAuthInfo": {
		Type: reflect.TypeFor[NamedAuthInfo](),
	},
	"AuthInfo": {
		Type: reflect.TypeFor[AuthInfo](),
	},
	"AuthProviderConfig": {
		Type: reflect.TypeFor[AuthProviderConfig](),
	},
	"ExecConfig": {
		Type: reflect.TypeFor[ExecConfig](),
	},
	"ExecEnvVar": {
		Type: reflect.TypeFor[ExecEnvVar](),
	},
	"NamedExtension": {
		Type: reflect.TypeFor[NamedExtension](),
	},
	"Preferences": {
		Type: reflect.TypeFor[Preferences](),
	},
	"Extension": {
		Type: reflect.TypeFor[Extension](),
	},
}

func TestClone(t *testing.T) {
	t.Parallel()
	for name, tc := range cloneTestcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var gen testutil.BasicTypeGenerator
			inst := testutil.GenerateCloneTestdata(&gen, reflect.PointerTo(tc.Type))

			// Get the Clone() method.
			cloneMethod, ok := inst.Type().MethodByName("Clone")
			if !ok {
				t.Fatal("does not have Clone() pointer-receiver method")
			}

			// Validate the Clone() method.
			if cloneMethod.Type.NumIn() != 1 ||
				cloneMethod.Type.NumOut() != 1 ||
				cloneMethod.Type.Out(0) != inst.Type() {
				t.Fatalf(
					"Clone() should be `func(*kubeconfig.%v) *kubeconfig.%v`, got `%v`",
					tc.Type.Name(),
					tc.Type.Name(),
					cloneMethod.Type,
				)
			}

			// Call the Clone() method.
			cloned := inst.Method(cloneMethod.Index).Call([]reflect.Value{})[0]

			// Compare the values.
			issue := testutil.VerifyClone(inst, cloned, tc.Type.Name())
			if issue != nil {
				t.Fatal(issue)
			}
		})
	}
}

func TestCloneInto(t *testing.T) {
	t.Parallel()
	for name, tc := range cloneTestcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var gen testutil.BasicTypeGenerator
			inst := testutil.GenerateCloneTestdata(&gen, tc.Type)

			// Get the CloneInto() method.
			cloneMethod, ok := inst.Type().MethodByName("CloneInto")
			if !ok {
				t.Fatal("does not have CloneInto() value-receiver method")
			}

			// Validate the Clone() method.
			if cloneMethod.Type.NumIn() != 2 ||
				cloneMethod.Type.NumOut() != 0 ||
				cloneMethod.Type.In(1) != reflect.PointerTo(inst.Type()) {
				t.Fatalf(
					"CloneInto() should be `func(kubeconfig.%v, *kubeconfig.%v)`, got `%v`",
					tc.Type.Name(),
					tc.Type.Name(),
					cloneMethod.Type,
				)
			}

			// Call the Clone() method.
			cloned := reflect.New(tc.Type)
			inst.Method(cloneMethod.Index).Call([]reflect.Value{cloned})

			// Compare the values.
			issue := testutil.VerifyClone(inst, cloned.Elem(), tc.Type.Name())
			if issue != nil {
				t.Fatal(issue)
			}
		})
	}
}
