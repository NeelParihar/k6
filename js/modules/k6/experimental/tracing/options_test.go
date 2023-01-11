package tracing

import "testing"

func TestOptionsValidate(t *testing.T) {
	t.Parallel()

	var testSampling int = 10.0

	type fields struct {
		Propagator string
		Sampling   *int
		Baggage    map[string]string
	}
	testCases := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "w3c propagator is valid",
			fields: fields{
				Propagator: "w3c",
			},
			wantErr: false,
		},
		{
			name: "b3 propagator is valid",
			fields: fields{
				Propagator: "b3",
			},
			wantErr: false,
		},
		{
			name: "jaeger propagator is valid",
			fields: fields{
				Propagator: "jaeger",
			},
			wantErr: false,
		},
		{
			name: "invalid propagator is invalid",
			fields: fields{
				Propagator: "invalid",
			},
			wantErr: true,
		},
		{
			name: "sampling is not yet supported",
			fields: fields{
				Sampling: &testSampling,
			},
			wantErr: true,
		},
		{
			name: "baggage is not yet supported",
			fields: fields{
				Baggage: map[string]string{"key": "value"},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			i := &options{
				Propagator: tc.fields.Propagator,
				Sampling:   tc.fields.Sampling,
				Baggage:    tc.fields.Baggage,
			}

			if err := i.validate(); (err != nil) != tc.wantErr {
				t.Errorf("instrumentationOptions.validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
