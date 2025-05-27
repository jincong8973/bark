package deepseek

import "testing"

func Test_removeThinkProcess(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "移除单个think标签",
			args: args{
				content: "这是正常文本<think>这是思考过程</think>这是后续文本",
			},
			want: "这是正常文本这是后续文本",
		},
		{
			name: "移除多个think标签",
			args: args{
				content: "开始<think>思考1</think>中间<think>思考2</think>结束",
			},
			want: "开始中间结束",
		},
		{
			name: "没有think标签",
			args: args{
				content: "这是普通文本，没有think标签",
			},
			want: "这是普通文本，没有think标签",
		},
		{
			name: "think标签包含多行内容",
			args: args{
				content: "开始<think>第一行\n第二行\n第三行</think>结束",
			},
			want: "开始结束",
		},
		{
			name: "空字符串",
			args: args{
				content: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeThinkProcess(tt.args.content); got != tt.want {
				t.Errorf("removeThinkProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
