package job

import (
	"testing"

	mockConfig "github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/task"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		dep  func() Dependency
		want []task.Runner
	}{
		{
			name: "FeatureFlagOff",
			dep: func() Dependency {
				configMock := mockConfig.NewMockConfig(t)

				configMock.EXPECT().
					GetBool("feature.flag.todo.job").
					Return(false)

				return Dependency{
					Config: configMock,
				}
			},
			want: nil,
		},
		{
			name: "FeatureFlagOn",
			dep: func() Dependency {
				configMock := mockConfig.NewMockConfig(t)

				configMock.EXPECT().
					GetBool("feature.flag.todo.job").
					Return(true)

				return Dependency{
					Config: configMock,
				}
			},
			want: []task.Runner{
				&todoPublisher{
					cjson: nil,
					mc:    nil,
					tel:   nil,
					topic: "todo.creator.topic",
				},
				&todoSubscriber{
					cjson:        nil,
					mc:           nil,
					tel:          nil,
					createUC:     nil,
					topic:        "todo.creator.topic",
					subscription: "gostarter.todo.creator.subscription",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New(tt.dep())
			assert.Equal(t, tt.want, got)
		})
	}
}
