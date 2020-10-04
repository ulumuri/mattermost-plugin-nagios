package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ulumuri/go-nagios/nagios"
)

func Test_gettingReportUnsuccessfulMessage(t *testing.T) {
	type args struct {
		reportPart string
		message    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				reportPart: "a part",
				message:    "a message",
			},
			want: "Getting monitoring report unsuccessful (a part): a message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gettingReportUnsuccessfulMessage(tt.args.reportPart, tt.args.message); got != tt.want {
				t.Errorf("gettingReportUnsuccessfulMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reportPreamble(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{
			name: "basic",
			t:    now,
			want: "#### :bar_chart: System monitoring report (" + now.Format(time.UnixDate) + ")\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reportPreamble(tt.t); got != tt.want {
				t.Errorf("reportPreamble() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatHostCount(t *testing.T) {
	tests := []struct {
		name  string
		count nagios.HostCount
		want  string
	}{
		{
			name: "basic",
			count: nagios.HostCount{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: nagios.HostCountData{
					Count: nagios.HostStatusCount{
						Up:          1,
						Down:        2,
						Unreachable: 3,
						Pending:     4,
					},
				},
			},
			want: "##### HOST SUMMARY\n\n:up: Up: **1**  :small_red_triangle_" +
				"down: Down: **2**  :mailbox_with_no_mail: Unreachable: **3**" +
				"  :hourglass_flowing_sand: Pending: **4**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatHostCount(tt.count); got != tt.want {
				t.Errorf("formatHostCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateHomogenousHostListData(state string, n int) nagios.HostListData {
	ret := nagios.HostListData{
		HostList: make(map[string]json.RawMessage),
	}

	m := json.RawMessage(fmt.Sprintf(`"%s"`, state))

	for i := 0; i < n; i++ {
		ret.HostList[strconv.Itoa(i)] = m
	}

	return ret
}

// Test_formatHostList only tests the edge cases to make sure we don't hit the
// Mattermost's message limits.
func Test_formatHostList(t *testing.T) {
	tests := []struct {
		name string
		list nagios.HostList
		want string
	}{
		{
			name: "empty",
			list: nagios.HostList{},
			want: gettingReportUnsuccessfulMessage("host list", ""),
		},
		{
			name: "empty successful",
			list: nagios.HostList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "##### HOST LIST\n\nNo hosts to show.",
		},
		{
			name: "too many hosts (all UP)",
			list: nagios.HostList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: generateHomogenousHostListData(upState, maximumReportLength+1),
			},
			want: "##### HOST LIST\n\n**Too many hosts. Showing only abnormal" +
				" state hosts.**\n\nNo hosts to show.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatHostList(tt.list); got != tt.want {
				t.Errorf("formatHostList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatServiceCount(t *testing.T) {
	tests := []struct {
		name  string
		count nagios.ServiceCount
		want  string
	}{
		{
			name: "basic",
			count: nagios.ServiceCount{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: nagios.ServiceCountData{
					Count: nagios.ServiceStatusCount{
						Ok:       1,
						Warning:  2,
						Critical: 3,
						Unknown:  4,
						Pending:  5,
					},
				},
			},
			want: "##### SERVICE SUMMARY\n\n:white_check_mark: OK: **1**  :wa" +
				"rning: Warning: **2**  :bangbang: Critical: **3**  :question" +
				": Unknown: **4**  :hourglass_flowing_sand: Pending: **5**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatServiceCount(tt.count); got != tt.want {
				t.Errorf("formatServiceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_calculateServiceState(t *testing.T) {
//	tests := []struct {
//		name       string
//		rawMessage json.RawMessage
//		want       string
//	}{
//		{
//			name: "all OK",
//			rawMessage: []byte(`{"Bandwidth Spike":"ok","Facebook Usage":"ok"` +
//				`,"Port 80 Bandwidth":"ok","Windows Failed Logins":"ok","Yout` +
//				`ube Usage":"ok"}`),
//			want: okState,
//		},
//		{
//			name: "most warning",
//			rawMessage: []byte(`{"/ Disk Usage":"ok","Apache 404 Errors":"ok"` +
//				`,"Apache Web Server":"ok","Bandwidth Spike":"ok","CPU Stats"` +
//				`:"ok","Cron Scheduling Daemon":"ok","Failed SSH Logins":"ok"` +
//				`,"Linux Failed Logins":"ok","Load":"ok","Memory Usage":"ok",` +
//				`"MySQL Crashed Tables":"ok","MySQL Server":"ok","Open Files"` +
//				`:"ok","Ping":"ok","Port 22 Bandwidth":"ok","Port 80 Bandwidt` +
//				`h":"ok","SSH Server":"ok","Swap Usage":"ok","System Logging ` +
//				`Daemon":"ok","Total Processes":"warning","Users":"ok","Youtu` +
//				`be Usage":"ok","Yum Updates":"warning"}`),
//			want: warningState,
//		},
//		{
//			name: "most critical",
//			rawMessage: []byte(`{"Auroral Activity":"ok","Weather Carteret No` +
//				`rth Carolina":"critical","Weather King Washington":"critical` +
//				`","Weather Ramsey Minnesota":"critical","Weather San Bernard` +
//				`ino California":"critical","Weather Strafford New Hampshire"` +
//				`:"critical","Weather Tulsa Oklahoma":"critical"}`),
//			want: criticalState,
//		},
//		{
//			name: "some warnings, most critical",
//			rawMessage: []byte(`{"Bandwidth Spike":"ok","Ping":"ok","Port 1 B` +
//				`andwidth":"ok","Port 1 Status":"ok","Port 10 Bandwidth":"ok"` +
//				`,"Port 10 Status":"critical","Port 11 Bandwidth":"ok","Port ` +
//				`11 Status":"critical","Port 12 Bandwidth":"ok","Port 12 Stat` +
//				`us":"warning","Port 13 Bandwidth":"ok","Port 13 Status":"cri` +
//				`tical","Port 14 Bandwidth":"ok","Port 14 Status":"critical",` +
//				`"Port 15 Bandwidth":"ok","Port 15 Status":"critical","Port 1` +
//				`6 Bandwidth":"ok","Port 16 Status":"warning","Port 17 Bandwi` +
//				`dth":"ok","Port 17 Status":"warning","Port 18 Bandwidth":"ok` +
//				`","Port 18 Status":"critical","Port 19 Bandwidth":"ok","Port` +
//				` 19 Status":"critical","Port 2 Bandwidth":"ok","Port 2 Statu` +
//				`s":"critical","Port 20 Bandwidth":"ok","Port 20 Status":"ok"` +
//				`,"Port 21 Bandwidth":"ok","Port 21 Status":"ok","Port 22 Ban` +
//				`dwidth":"ok","Port 22 Status":"warning","Port 23 Bandwidth":` +
//				`"ok","Port 23 Status":"warning","Port 24 Bandwidth":"ok","Po` +
//				`rt 24 Status":"ok","Port 25 Bandwidth":"ok","Port 25 Status"` +
//				`:"ok","Port 3 Bandwidth":"ok","Port 3 Status":"ok","Port 4 B` +
//				`andwidth":"ok","Port 4 Status":"warning","Port 5 Bandwidth":` +
//				`"ok","Port 5 Status":"ok","Port 6 Bandwidth":"ok","Port 6 St` +
//				`atus":"critical","Port 7 Bandwidth":"ok","Port 7 Status":"cr` +
//				`itical","Port 8 Bandwidth":"ok","Port 8 Status":"critical","` +
//				`Port 9 Bandwidth":"ok","Port 9 Status":"ok","Youtube Usage":` +
//				`"warning"}`),
//			want: criticalState,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := calculateServiceState(tt.rawMessage); got != tt.want {
//				t.Errorf("calculateServiceState() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func generateHomogenousServiceListData(state string, n int) nagios.ServiceListData {
	ret := nagios.ServiceListData{
		ServiceList: make(map[string]json.RawMessage),
	}

	m := json.RawMessage(fmt.Sprintf(`{"parameter":"%s"}`, state))

	for i := 0; i < n; i++ {
		ret.ServiceList[strconv.Itoa(i)] = m
	}

	return ret
}

// Test_formatServiceList only tests the edge cases to make sure we don't hit
// the Mattermost's message limits.
func Test_formatServiceList(t *testing.T) {
	tests := []struct {
		name string
		list nagios.ServiceList
		want string
	}{
		{
			name: "empty",
			list: nagios.ServiceList{},
			want: gettingReportUnsuccessfulMessage("service list", ""),
		},
		{
			name: "empty successful",
			list: nagios.ServiceList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "##### SERVICE LIST\n\nNo services to show.",
		},
		{
			name: "too many services (all OK)",
			list: nagios.ServiceList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: generateHomogenousServiceListData(okState, maximumReportLength+1),
			},
			want: "##### SERVICE LIST\n\n**Too many services. Showing only ab" +
				"normal state services.**\n\nNo services to show.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatServiceList(tt.list); got != tt.want {
				t.Errorf("formatServiceList() = %v, want %v", got, tt.want)
			}
		})
	}
}