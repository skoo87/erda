// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quality_chart

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func Test_radar(t *testing.T) {
	ctx := context.Background()
	context.WithValue(ctx, cptype.GlobalInnerKeyCtxSDK, &cptype.SDK{})

	radar := charts.NewRadar()
	radar.Indicator = []*opts.Indicator{
		{Name: "test-case", Max: 100, Min: 0, Color: ""},
		{Name: "config-sheet", Max: 100, Min: 0, Color: ""},
		{Name: "test-plan", Max: 100, Min: 0, Color: ""},
		{Name: "code-quality", Max: 100, Min: 0, Color: ""},
		{Name: "issue-bug", Max: 100, Min: 0, Color: ""},
	}
	radar.AddSeries(
		"quality",
		[]opts.RadarData{
			{
				Name:  "",
				Value: []uint64{100, 90, 80, 70, 60},
			},
		},
		charts.WithAreaStyleOpts(opts.AreaStyle{Color: "", Opacity: 0.2}),
		charts.WithLabelOpts(opts.Label{Show: true, Color: "", Position: "", Formatter: ""}),
	)
	radar.SetGlobalOptions(
		charts.WithTooltipOpts(opts.Tooltip{Show: false, Trigger: "item"}),
		charts.WithTitleOpts(opts.Title{Title: "quality score"}),
	)
	b, err := json.MarshalIndent(radar.JSON(), "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(b))
}

func Test_polishScore(t *testing.T) {
	type args struct {
		score float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "less than 0",
			args: args{
				score: -1,
			},
			want: float64(0),
		},
		{
			name: "bigger than 100",
			args: args{
				score: 101,
			},
			want: float64(100),
		},
		{
			name: "3 digit",
			args: args{
				score: 11.234,
			},
			want: float64(11.23),
		},
		{
			name: "1 digit",
			args: args{
				score: 11.2,
			},
			want: float64(11.2),
		},
		{
			name: "no digit",
			args: args{
				score: 11,
			},
			want: float64(11),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := polishScore(tt.args.score); got != tt.want {
				t.Errorf("polishScore() = %v, want %v", got, tt.want)
			}
		})
	}
}